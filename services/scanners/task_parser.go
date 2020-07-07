package scanners

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/oasisprotocol/oasis-core/go/common/cbor"
	"github.com/oasisprotocol/oasis-core/go/common/crypto/address"
	consensusAPI "github.com/oasisprotocol/oasis-core/go/consensus/api"
	"github.com/oasisprotocol/oasis-core/go/consensus/api/transaction"
	stakingAPI "github.com/oasisprotocol/oasis-core/go/staking/api"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/types"
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

	err = p.parseOasisBase(blockData, baseFlag)
	if err != nil {
		return fmt.Errorf("p.parseOasisBase: %s", err.Error())
	}

	return nil
}

func (p *ParserTask) BalanceSnapshot(blockID uint64) error {

	blockData, err := p.consensusAPI.GetBlock(p.ctx, int64(blockID))
	if err != nil {
		return fmt.Errorf("api.Block.Get: %s", err.Error())
	}

	err = p.parseOasisBase(blockData, balanceSnapshotFlag)
	if err != nil {
		return fmt.Errorf("p.parseOasisBase: %s", err.Error())
	}

	return nil
}

func (p *ParserTask) parseOasisBase(blockData *consensusAPI.Block, parseFlag ParseFlag) (err error) {

	b := oasis.Block{}
	//Nil pointer err
	err = cbor.Unmarshal(blockData.Meta, &b)
	if err != nil {
		return err
	}

	b.Hash = blockData.Hash

	var pipes []func(data oasis.Block) error

	if (parseFlag & baseFlag) != 0 {
		pipes = append(pipes, []func(data oasis.Block) error{
			p.parseBlock,
			p.parseBlockSignatures,
			p.parseBlockTransactions,
		}...)
	}

	if (parseFlag & balanceSnapshotFlag) != 0 {
		pipes = append(pipes, []func(data oasis.Block) error{
			p.epochBalanceSnapshot,
		}...)
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
	var nodeRegisterTxs []dmodels.NodeRegistryTransaction
	var entityRegisterTxs []dmodels.EntityRegistryTransaction

	accountBalanceUpdates := make([]dmodels.AccountBalance, 0, len(txs))
	for i := range txs {
		tx := transaction.SignedTransaction{}

		err = cbor.Unmarshal(txs[i], &tx)
		if err != nil {
			return err
		}

		raw := oasis.UntrustedRawValue{}
		err = cbor.Unmarshal(tx.Blob, &raw)
		if err != nil {
			return err
		}

		txType, err := dmodels.NewTransactionType(raw.Method)
		if err != nil {
			return err
		}

		//Parse NodeRegister Tx
		nodeRegisterTx, err := p.parseNodeRegistryTransaction(txType, block, raw)
		if err != nil {
			return err
		}

		if nodeRegisterTx.ID != "" {
			//Make hash from origin []byte
			nodeRegisterTx.Hash = hex.EncodeToString(types.Tx(txs[i]).Hash())
			nodeRegisterTxs = append(nodeRegisterTxs, nodeRegisterTx)
		}

		//Parse EntityRegister Tx
		entityRegisterTx, err := p.parseEntityRegistryTransaction(txType, block, raw)
		if err != nil {
			return err
		}

		if entityRegisterTx.ID != "" {
			//Make hash from origin []byte
			entityRegisterTx.Hash = hex.EncodeToString(types.Tx(txs[i]).Hash())
			entityRegisterTxs = append(entityRegisterTxs, entityRegisterTx)
		}

		//Epoch balance snapshots on another job
		if !block.IsEpochBlock() {
			balanceUpdates, err := p.parseAccountBalances(block, tx, raw)
			if err != nil {
				return err
			}
			accountBalanceUpdates = append(accountBalanceUpdates, balanceUpdates...)
		}

		receiver := raw.Body.To.String()
		if txType.Type() == dmodels.TransactionTypeAddEscrow || txType.Type() == dmodels.TransactionTypeReclaimEscrow {
			receiver = raw.Body.EscrowTx.Account.String()
		}

		dTxs[i] = dmodels.Transaction{
			BlockLevel:          uint64(block.Header.Height),
			BlockHash:           hex.EncodeToString(block.Hash),
			Hash:                tx.Hash().String(),
			Time:                block.Header.Time,
			Amount:              raw.Body.Transfer.Tokens.ToBigInt().Uint64(),
			EscrowAmount:        raw.Body.EscrowTx.Tokens.ToBigInt().Uint64(),
			EscrowReclaimAmount: raw.Body.EscrowTx.Shares.ToBigInt().Uint64(),
			Type:                txType.Type(),
			Sender:              stakingAPI.NewAddress(tx.Signature.PublicKey).String(),
			Receiver:            receiver,
			Nonce:               raw.Nonce,
			Fee:                 raw.Fee.Amount.ToBigInt().Uint64(),
			GasLimit:            uint64(raw.Fee.Gas),
			GasPrice:            raw.Fee.GasPrice().ToBigInt().Uint64(),
		}
	}

	p.parserContainer.txs.Add(dTxs, nodeRegisterTxs, entityRegisterTxs)
	p.parserContainer.balances.Add(accountBalanceUpdates)

	return nil
}

func (p *ParserTask) epochBalanceSnapshot(block oasis.Block) error {
	//Make snapshot only for epoch blocks
	if !block.IsEpochBlock() {
		return nil
	}

	addresses, err := p.stakingAPI.Addresses(p.ctx, block.Header.Height)
	if err != nil {
		return err
	}

	updates := make([]dmodels.AccountBalance, 0, len(addresses))

	for i := range addresses {

		balance, err := p.getAccountBalance(block.Header.Height, addresses[i])
		if err != nil {
			return err
		}

		balance.Time = block.Header.Time
		updates = append(updates, balance)
	}

	p.parserContainer.balances.Add(updates)

	return nil
}

func (p *ParserTask) parseAccountBalances(block oasis.Block, tx transaction.SignedTransaction, rawTX oasis.UntrustedRawValue) ([]dmodels.AccountBalance, error) {
	addresses := []stakingAPI.Address{stakingAPI.NewAddress(tx.Signature.PublicKey), rawTX.Body.EscrowTx.Account, rawTX.Body.To}

	updates := make([]dmodels.AccountBalance, 0, len(addresses))

	for i := range addresses {
		//Skip system address
		if (address.Address)(addresses[i]).Equal(oasis.SystemAddress) {
			continue
		}

		balance, err := p.getAccountBalance(block.Header.Height, addresses[i])
		if err != nil {
			return updates, err
		}

		balance.Time = block.Header.Time
		updates = append(updates, balance)
	}

	return updates, nil
}

func (p *ParserTask) getAccountBalance(height int64, address stakingAPI.Address) (balance dmodels.AccountBalance, err error) {

	accInfo, err := p.stakingAPI.Account(p.ctx, &stakingAPI.OwnerQuery{
		Height: height,
		Owner:  address,
	})
	if err != nil {
		return balance, err
	}

	return dmodels.AccountBalance{
		Account:               address.String(),
		Height:                height,
		Nonce:                 accInfo.General.Nonce,
		GeneralBalance:        accInfo.General.Balance.ToBigInt().Uint64(),
		EscrowBalanceActive:   accInfo.Escrow.Active.Balance.ToBigInt().Uint64(),
		EscrowBalanceShare:    accInfo.Escrow.Active.TotalShares.ToBigInt().Uint64(),
		EscrowDebondingActive: accInfo.Escrow.Debonding.Balance.ToBigInt().Uint64(),
		EscrowDebondingShare:  accInfo.Escrow.Debonding.TotalShares.ToBigInt().Uint64(),
	}, nil
}

func (p *ParserTask) parseNodeRegistryTransaction(txType dmodels.TransactionMethod, block oasis.Block, raw oasis.UntrustedRawValue) (registerTx dmodels.NodeRegistryTransaction, err error) {
	if txType.Type() != dmodels.TransactionTypeRegisterNode {
		return
	}

	regNode := oasis.RegisterNode{}
	err = cbor.Unmarshal(raw.Body.RegisterTx.Blob, &regNode)
	if err != nil {
		return registerTx, err
	}

	var physicalAddress string
	if len(regNode.Consensus.Addresses) > 0 {
		physicalAddress = regNode.Consensus.Addresses[0].Address.String()
	}

	consensusIDBytes, err := regNode.Consensus.ID.MarshalBinary()
	if err != nil {
		return registerTx, err
	}

	return dmodels.NodeRegistryTransaction{
		BlockLevel:       uint64(block.Header.Height),
		Time:             block.Header.Time,
		ID:               regNode.ID.String(),
		Address:          stakingAPI.NewAddress(regNode.ID).String(),
		EntityID:         regNode.EntityID.String(),
		EntityAddress:    stakingAPI.NewAddress(regNode.EntityID).String(),
		Expiration:       regNode.Expiration,
		P2PID:            regNode.P2P.ID.String(),
		ConsensusID:      regNode.Consensus.ID.String(),
		ConsensusAddress: crypto.AddressHash(consensusIDBytes).String(),
		PhysicalAddress:  physicalAddress,
		Roles:            uint32(regNode.Roles),
	}, nil
}

func (p *ParserTask) parseEntityRegistryTransaction(txType dmodels.TransactionMethod, block oasis.Block, raw oasis.UntrustedRawValue) (registerTx dmodels.EntityRegistryTransaction, err error) {
	if txType.Type() != dmodels.TransactionTypeRegisterEntity {
		return
	}

	regEntity := oasis.RegisterEntity{}
	err = cbor.Unmarshal(raw.Body.RegisterTx.Blob, &regEntity)
	if err != nil {
		return registerTx, err
	}

	nodes := make([]string, len(regEntity.Nodes))
	for i := range regEntity.Nodes {
		nodes[i] = regEntity.Nodes[i].String()
	}

	return dmodels.EntityRegistryTransaction{
		BlockLevel:             uint64(block.Header.Height),
		Time:                   block.Header.Time,
		ID:                     regEntity.ID.String(),
		Address:                stakingAPI.NewAddress(regEntity.ID).String(),
		Nodes:                  nodes,
		AllowEntitySignedNodes: regEntity.AllowEntitySignedNodes,
	}, nil
}
