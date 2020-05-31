package scanners

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/fxamacker/cbor/v2"
	"github.com/oasislabs/oasis-core/go/common/crypto/signature"
	consensusAPI "github.com/oasislabs/oasis-core/go/consensus/api"
	stakingAPI "github.com/oasislabs/oasis-core/go/staking/api"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/types"
	"github.com/wedancedalot/decimal"
	"google.golang.org/grpc"
	"oasisTracker/dmodels"
	"oasisTracker/dmodels/oasis"
	"reflect"
	"runtime"
)

//Struct for direct work with connection from worker
type ParserTask struct {
	ctx             context.Context
	consensusAPI    consensusAPI.ClientBackend
	stakingAPI      stakingAPI.Backend
	parserContainer *ParseContainer
}

func NewParserTask(ctx context.Context, conn *grpc.ClientConn, parserContainer *ParseContainer) (*ParserTask, error) {
	consensusAPI := consensusAPI.NewConsensusClient(conn)
	stakingAPI := stakingAPI.NewStakingClient(conn)

	return &ParserTask{ctx: ctx, consensusAPI: consensusAPI, stakingAPI: stakingAPI, parserContainer: parserContainer}, nil
}

func (p *ParserTask) ParseBase(blockID uint64) error {
	blockData, err := p.consensusAPI.GetBlock(p.ctx, int64(blockID))
	if err != nil {
		return fmt.Errorf("api.Block.Get: %s", err.Error())
	}

	err = p.ParseOasisBase(blockData)
	if err != nil {
		return fmt.Errorf("p.parseOasisBase: %s", err.Error())
	}

	return nil
}

func (p *ParserTask) ParseOasisBase(blockData *consensusAPI.Block) (err error) {

	b := oasis.Block{}
	err = cbor.Unmarshal(blockData.Meta, &b)
	if err != nil {
		return err
	}

	b.Hash = blockData.Hash

	pipes := []func(data oasis.Block) error{
		p.parseBlock,
		p.parseBlockSignatures,
		p.parseBlockTransactions,
		p.epochBalanceSnapshot,
	}

	for _, pipe := range pipes {
		err = pipe(b)
		if err != nil {
			funcName := runtime.FuncForPC(reflect.ValueOf(pipe).Pointer()).Name()
			return fmt.Errorf("%s (block:%d): %s", funcName, blockData.Height, err.Error())
		}
	}
	return nil
}

func (p *ParserTask) parseBlock(block oasis.Block) error {
	epoch, err := p.consensusAPI.GetEpoch(p.ctx, block.Header.Height)
	if err != nil {
		return err
	}

	p.parserContainer.blocks.Add([]dmodels.Block{{
		Height:          uint64(block.Header.Height),
		Hash:            hex.EncodeToString(block.Hash),
		CreatedAt:       block.Header.Time,
		Epoch:           uint64(epoch),
		ProposerAddress: block.Header.ProposerAddress.String(),
		ValidatorHash:   block.Header.ValidatorsHash.String(),
	}})

	return nil
}

func (p *ParserTask) parseBlockSignatures(block oasis.Block) error {

	blockSignatures := make([]dmodels.BlockSignature, 0, len(block.LastCommit.Signatures))

	for key := range block.LastCommit.Signatures {
		timestamp := block.LastCommit.Signatures[key].Timestamp
		if timestamp.IsZero() {
			timestamp = block.Header.Time
		}

		//Use block timestamp if signature time is zero
		blockSignatures = append(blockSignatures, dmodels.BlockSignature{
			BlockHeight:      block.Header.Height,
			Timestamp:        timestamp,
			BlockIDFlag:      block.LastCommit.Signatures[key].BlockIDFlag,
			ValidatorAddress: block.LastCommit.Signatures[key].ValidatorAddress.String(),
			Signature:        hex.EncodeToString(block.LastCommit.Signatures[key].Signature),
		})

	}

	p.parserContainer.blockSignatures.Add(blockSignatures)

	return nil
}

func (p *ParserTask) parseBlockTransactions(block oasis.Block) (err error) {
	txs, err := p.consensusAPI.GetTransactions(p.ctx, block.Header.Height)
	if err != nil {
		return err
	}

	dTxs := make([]dmodels.Transaction, len(txs))
	var registerTxs []dmodels.RegistryTransaction
	accountBalanceUpdates := make([]dmodels.AccountBalance, 0, len(txs))
	for i := range txs {
		tx := oasis.TxRaw{}

		err = cbor.Unmarshal(txs[i], &tx)
		if err != nil {
			return err
		}

		raw := oasis.UntrustedRawValue{}
		err = cbor.Unmarshal(tx.UntrustedRawValue, &raw)
		if err != nil {
			return err
		}

		txType, err := dmodels.NewTransactionType(raw.Method)
		if err != nil {
			return err
		}

		if txType.Type() == dmodels.TransactionTypeRegisterNode {
			regNode := oasis.RegisterNode{}
			err = cbor.Unmarshal(raw.Body.RegisterTx.Blob, &regNode)
			if err != nil {
				return err
			}

			var physicalAddress string
			if len(regNode.Consensus.Addresses) > 0 {
				physicalAddress = regNode.Consensus.Addresses[0].Address.String()
			}

			entityBytes, err := regNode.EntityID.MarshalBinary()
			if err != nil {
				return err
			}

			consensusIDBytes, err := regNode.Consensus.ID.MarshalBinary()
			if err != nil {
				return err
			}

			registerTxs = append(registerTxs, dmodels.RegistryTransaction{
				BlockLevel:       uint64(block.Header.Height),
				Hash:             hex.EncodeToString(types.Tx(txs[i]).Hash()),
				Time:             block.Header.Time,
				ID:               regNode.ID.String(),
				EntityID:         regNode.EntityID.String(),
				EntityAddress:    crypto.AddressHash(entityBytes).String(),
				Expiration:       regNode.Expiration,
				P2PID:            regNode.P2P.ID.String(),
				ConsensusID:      regNode.Consensus.ID.String(),
				ConsensusAddress: crypto.AddressHash(consensusIDBytes).String(),
				PhysicalAddress:  physicalAddress,
				Roles:            uint32(regNode.Roles),
			})
		}

		if !block.IsEpochBlock() {
			balanceUpdates, err := p.parseAccountBalances(block, tx, raw)
			if err != nil {
				return err
			}
			accountBalanceUpdates = append(accountBalanceUpdates, balanceUpdates...)
		}

		dTxs[i] = dmodels.Transaction{
			BlockLevel:          uint64(block.Header.Height),
			BlockHash:           block.Hash.String(),
			Hash:                hex.EncodeToString(types.Tx(txs[i]).Hash()),
			Time:                block.Header.Time,
			Amount:              decimal.NewFromBigInt(raw.Body.Transfer.Tokens.ToBigInt(), -int32(dmodels.Precision)).String(),
			EscrowAmount:        decimal.NewFromBigInt(raw.Body.EscrowTx.Tokens.ToBigInt(), -int32(dmodels.Precision)).String(),
			EscrowReclaimAmount: decimal.NewFromBigInt(raw.Body.EscrowTx.Shares.ToBigInt(), -int32(dmodels.Precision)).String(),
			EscrowAccount:       raw.Body.EscrowTx.Account.String(),
			Type:                txType.Type(),
			Sender:              tx.Signature.PublicKey.String(),
			Receiver:            raw.Body.To.String(),
			Nonce:               raw.Nonce,
			Fee:                 raw.Fee.Amount.ToBigInt().Uint64(),
			GasLimit:            uint64(raw.Fee.Gas),
			GasPrice:            raw.Fee.GasPrice().ToBigInt().Uint64(),
		}
	}

	p.parserContainer.txs.Add(dTxs, registerTxs)
	p.parserContainer.balances.Add(accountBalanceUpdates)

	return nil
}

func (p *ParserTask) epochBalanceSnapshot(block oasis.Block) error {
	//Make snapshot only for epoch blocks
	if !block.IsEpochBlock() {
		return nil
	}

	accounts, err := p.stakingAPI.Accounts(p.ctx, block.Header.Height)
	if err != nil {
		return err
	}

	updates := make([]dmodels.AccountBalance, 0, len(accounts))

	for i := range accounts {

		balance, err := p.getAccountBalance(block.Header.Height, accounts[i])
		if err != nil {
			return err
		}

		balance.Time = block.Header.Time
		updates = append(updates, balance)
	}

	p.parserContainer.balances.Add(updates)

	return nil
}

func (p *ParserTask) parseAccountBalances(block oasis.Block, tx oasis.TxRaw, rawTX oasis.UntrustedRawValue) ([]dmodels.AccountBalance, error) {
	accounts := []signature.PublicKey{tx.Signature.PublicKey, rawTX.Body.EscrowTx.Account, rawTX.Body.To}
	updates := make([]dmodels.AccountBalance, 0, len(accounts))

	for i := range accounts {
		//Skip system address
		if accounts[i].Equal(oasis.SystemAddress) {
			continue
		}

		balance, err := p.getAccountBalance(block.Header.Height, accounts[i])
		if err != nil {
			return updates, err
		}

		balance.Time = block.Header.Time
		updates = append(updates, balance)
	}

	return updates, nil
}

func (p *ParserTask) getAccountBalance(height int64, pubKey signature.PublicKey) (balance dmodels.AccountBalance, err error) {

	accInfo, err := p.stakingAPI.AccountInfo(p.ctx, &stakingAPI.OwnerQuery{
		Height: height,
		Owner:  pubKey,
	})
	if err != nil {
		return balance, err
	}

	return dmodels.AccountBalance{
		Account:               pubKey.String(),
		Height:                height,
		GeneralBalance:        accInfo.General.Balance.String(),
		Nonce:                 accInfo.General.Nonce,
		EscrowBalanceActive:   accInfo.Escrow.Active.Balance.String(),
		EscrowBalanceShare:    accInfo.Escrow.Active.TotalShares.String(),
		EscrowDebondingActive: accInfo.Escrow.Debonding.Balance.String(),
		EscrowDebondingShare:  accInfo.Escrow.Debonding.TotalShares.String(),
	}, nil
}
