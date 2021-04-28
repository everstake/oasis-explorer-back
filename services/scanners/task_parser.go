package scanners

import (
	"context"
	"encoding/hex"
	"fmt"
	"oasisTracker/dmodels"
	"oasisTracker/dmodels/oasis"
	"reflect"
	"runtime"

	beaconAPI "github.com/oasisprotocol/oasis-core/go/beacon/api"
	"github.com/oasisprotocol/oasis-core/go/common/cbor"
	"github.com/oasisprotocol/oasis-core/go/common/crypto/address"
	consensusAPI "github.com/oasisprotocol/oasis-core/go/consensus/api"
	"github.com/oasisprotocol/oasis-core/go/consensus/api/transaction"
	registryAPI "github.com/oasisprotocol/oasis-core/go/registry/api"
	stakingAPI "github.com/oasisprotocol/oasis-core/go/staking/api"
	"github.com/tendermint/tendermint/crypto"
	"google.golang.org/grpc"
)

//Struct for direct work with connection from worker
type ParserTask struct {
	ctx             context.Context
	consensusAPI    consensusAPI.ClientBackend
	beaconAPI       beaconAPI.Backend
	stakingAPI      stakingAPI.Backend
	registryAPI     registryAPI.Backend
	parserContainer *ParseContainer
}

func NewParserTask(ctx context.Context, conn *grpc.ClientConn, parserContainer *ParseContainer) (*ParserTask, error) {
	consensusAPI := consensusAPI.NewConsensusClient(conn)
	stakingAPI := stakingAPI.NewStakingClient(conn)
	registryAPI := registryAPI.NewRegistryClient(conn)
	beaconAPI := beaconAPI.NewBeaconClient(conn)

	return &ParserTask{ctx: ctx, consensusAPI: consensusAPI, stakingAPI: stakingAPI, registryAPI: registryAPI, beaconAPI: beaconAPI, parserContainer: parserContainer}, nil
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
	epoch, err := p.beaconAPI.GetEpoch(p.ctx, block.Header.Height)
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
	txsWithResults, err := p.consensusAPI.GetTransactionsWithResults(p.ctx, block.Header.Height)
	if err != nil {
		return err
	}

	dTxs := make([]dmodels.Transaction, len(txsWithResults.Transactions))
	var nodeRegisterTxs []dmodels.NodeRegistryTransaction
	var entityRegisterTxs []dmodels.EntityRegistryTransaction

	accountBalanceUpdates := make([]dmodels.AccountBalance, 0, len(txsWithResults.Transactions))
	tx := transaction.SignedTransaction{}

	for i := range txsWithResults.Transactions {
		err = cbor.Unmarshal(txsWithResults.Transactions[i], &tx)
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
			nodeRegisterTx.Hash = tx.Hash().String()
			nodeRegisterTxs = append(nodeRegisterTxs, nodeRegisterTx)
		}

		//Parse EntityRegister Tx
		entityRegisterTx, err := p.parseEntityRegistryTransaction(txType, block, raw)
		if err != nil {
			return err
		}

		if entityRegisterTx.ID != "" {
			entityRegisterTx.Hash = tx.Hash().String()
			entityRegisterTxs = append(entityRegisterTxs, entityRegisterTx)
		}

		//Save updated balances
		balanceUpdates, err := p.parseAccountBalances(block, tx, raw)
		if err != nil {
			return err
		}
		accountBalanceUpdates = append(accountBalanceUpdates, balanceUpdates...)

		receiver := raw.Body.To.String()
		if txType.Type() == dmodels.TransactionTypeAddEscrow || txType.Type() == dmodels.TransactionTypeReclaimEscrow {
			receiver = raw.Body.Account.String()
		}

		dTxs[i] = dmodels.Transaction{
			BlockLevel: uint64(block.Header.Height),
			BlockHash:  hex.EncodeToString(block.Hash),
			Hash:       tx.Hash().String(),
			Time:       block.Header.Time,
			//Same field
			//Todo probably merge to single field
			Amount:              raw.Body.Amount.ToBigInt().Uint64(),
			EscrowAmount:        raw.Body.Amount.ToBigInt().Uint64(),
			EscrowReclaimAmount: raw.Body.Shares.ToBigInt().Uint64(),

			//Save tx status and error if presented
			Status: txsWithResults.Results[i].IsSuccess(),
			Error:  txsWithResults.Results[i].Error.Message,

			Type:     txType.Type(),
			Sender:   stakingAPI.NewAddress(tx.Signature.PublicKey).String(),
			Receiver: receiver,
			Nonce:    raw.Nonce,
			Fee:      raw.Fee.Amount.ToBigInt().Uint64(),
			GasLimit: uint64(raw.Fee.Gas),
			GasPrice: raw.Fee.GasPrice().ToBigInt().Uint64(),
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

	entities, err := p.registryAPI.GetEntities(p.ctx, block.Header.Height)
	if err != nil {
		return err
	}

	if len(entities) == 0 {
		return nil
	}

	epoch, err := p.beaconAPI.GetEpoch(p.ctx, block.Header.Height)
	if err != nil {
		return err
	}

	updates := make([]dmodels.AccountBalance, 0, len(entities))
	rewards := make([]dmodels.Reward, 0, len(entities))
	var rewardsAmount uint64

	var entityAddress stakingAPI.Address
	for i := range entities {

		entityAddress = stakingAPI.NewAddress(entities[i].ID)

		balance, err := p.getAccountBalance(block.Header.Height, entityAddress)
		if err != nil {
			return err
		}

		balance.Time = block.Header.Time
		updates = append(updates, balance)

		prevBalance, err := p.getAccountBalance(block.Header.Height-1, entityAddress)
		if err != nil {
			return err
		}

		//Todo check txs
		rewardsAmount = (balance.EscrowBalanceActive - prevBalance.EscrowBalanceActive) + (balance.GeneralBalance - prevBalance.GeneralBalance)

		if rewardsAmount > 0 {
			rewards = append(rewards, dmodels.Reward{
				EntityAddress: entityAddress.String(),
				BlockLevel:    block.Header.Height,
				Epoch:         uint64(epoch),

				Amount:    rewardsAmount,
				CreatedAt: block.Header.Time,
			})
		}
	}

	debondingUpdates, err := p.processDebondingDelegations(block)
	if err != nil {
		return err
	}

	updates = append(updates, debondingUpdates...)

	p.parserContainer.balances.Add(updates)
	p.parserContainer.rewards.Add(rewards)

	return nil
}

func (p *ParserTask) processDebondingDelegations(block oasis.Block) (updates []dmodels.AccountBalance, err error) {
	//Save undelegations
	addresses, err := p.stakingAPI.Addresses(context.Background(), block.Header.Height)
	if err != nil {
		return updates, err
	}

	for _, address := range addresses {

		debondingDelegations, err := p.stakingAPI.DebondingDelegationsFor(p.ctx, &stakingAPI.OwnerQuery{
			Height: block.Header.Height,
			Owner:  address,
		})
		if err != nil {
			return updates, err
		}

		previousDebondingDelegations, err := p.stakingAPI.DebondingDelegationsFor(p.ctx, &stakingAPI.OwnerQuery{
			Height: block.Header.Height - 1,
			Owner:  address,
		})
		if err != nil {
			return updates, err
		}

		//Equal so skip
		if compareDebondingDelegations(debondingDelegations, previousDebondingDelegations) {
			continue
		}

		accountBalance, err := p.getAccountBalance(block.Header.Height, address)
		if err != nil {
			return updates, err
		}

		accountBalance.Time = block.Header.Time

		updates = append(updates, accountBalance)
	}

	return updates, nil
}

func compareDebondingDelegations(debondingDelegations, previousDebondingDelegations map[stakingAPI.Address][]*stakingAPI.DebondingDelegation) bool {

	if len(debondingDelegations) != len(previousDebondingDelegations) {
		return false
	}

	for address, debondings := range debondingDelegations {
		prevDebonding, ok := previousDebondingDelegations[address]
		if !ok {
			return false
		}

		if len(prevDebonding) != len(debondings) {
			return false
		}

		if prevDebonding[0].Shares.Cmp(&debondings[0].Shares) != 0 || prevDebonding[0].DebondEndTime != debondings[0].DebondEndTime {
			return false
		}
	}

	return true
}

func (p *ParserTask) parseAccountBalances(block oasis.Block, tx transaction.SignedTransaction, rawTX oasis.UntrustedRawValue) (updates []dmodels.AccountBalance, err error) {
	addresses := []stakingAPI.Address{stakingAPI.NewAddress(tx.Signature.PublicKey), rawTX.Body.Account, rawTX.Body.To}

	updates = make([]dmodels.AccountBalance, 0, len(addresses))

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

	delegations, err := p.stakingAPI.DelegationsFor(p.ctx, &stakingAPI.OwnerQuery{
		Height: height,
		Owner:  address,
	})

	var delegationsBalance uint64
	for _, balance := range delegations {
		delegationsBalance += balance.Shares.ToBigInt().Uint64()
	}

	debondingDelegations, err := p.stakingAPI.DebondingDelegationsFor(p.ctx, &stakingAPI.OwnerQuery{
		Height: height,
		Owner:  address,
	})

	var debondingDelegationsBalance uint64
	for _, debondings := range debondingDelegations {
		for _, value := range debondings {
			debondingDelegationsBalance += value.Shares.ToBigInt().Uint64()
		}
	}

	return dmodels.AccountBalance{
		Account:                     address.String(),
		Height:                      height,
		Nonce:                       accInfo.General.Nonce,
		GeneralBalance:              accInfo.General.Balance.ToBigInt().Uint64(),
		EscrowBalanceActive:         accInfo.Escrow.Active.Balance.ToBigInt().Uint64(),
		EscrowBalanceShare:          accInfo.Escrow.Active.TotalShares.ToBigInt().Uint64(),
		EscrowDebondingActive:       accInfo.Escrow.Debonding.Balance.ToBigInt().Uint64(),
		EscrowDebondingShare:        accInfo.Escrow.Debonding.TotalShares.ToBigInt().Uint64(),
		DelegationsBalance:          delegationsBalance,
		DebondingDelegationsBalance: debondingDelegationsBalance,
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
		BlockLevel: uint64(block.Header.Height),
		Time:       block.Header.Time,
		ID:         regEntity.ID.String(),
		Address:    stakingAPI.NewAddress(regEntity.ID).String(),
		Nodes:      nodes,
	}, nil
}
