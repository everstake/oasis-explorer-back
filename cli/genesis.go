package cli

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/oasisprotocol/oasis-core/go/common/cbor"
	"github.com/oasisprotocol/oasis-core/go/staking/api"
	"github.com/tendermint/tendermint/crypto"
	"oasisTracker/common/log"
	"oasisTracker/dao"
	"oasisTracker/dmodels"
	"oasisTracker/dmodels/oasis"
	"oasisTracker/smodels"
	"os"
)

const SetupGenesisJson = "setup-genesis"

type (
	ICli interface {
		Setup(args []string) error
		SetupGenesisJson(args []string) error
	}

	Cli struct {
		DAO dao.ParserDAO
	}
)

func NewCli(d dao.DAO) ICli {

	pDAO, _ := d.GetParserDAO()

	return &Cli{
		DAO: pDAO,
	}
}

const genesisHeight = 0

func (cli *Cli) Setup(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("invalid arguments length %d", len(args))
	}

	switch args[0] {
	case SetupGenesisJson:
		return cli.SetupGenesisJson(args[1:])
	default:
		return fmt.Errorf("unsupported setup mode %s", args[0])
	}

}

func (cli *Cli) SetupGenesisJson(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("invalid arguments length %d, expected 1 arguments for setup-genesis", len(args))
	}

	//Use root folder as default
	file, err := os.Open(fmt.Sprint("./", args[0]))
	if err != nil {
		return err
	}

	gen := smodels.GenesisDocument{}

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&gen)
	if err != nil {
		return err
	}

	genesisBlockHash := sha256.Sum256([]byte(gen.ChainID))

	balances := make([]dmodels.AccountBalance, len(gen.Staking.Ledger))

	i := 0
	//Genesis balances
	for accountAddress, balance := range gen.Staking.Ledger {

		balances[i] = dmodels.AccountBalance{
			Account:               accountAddress.String(),
			Time:                  gen.GenesisTime,
			Height:                genesisHeight,
			Nonce:                 balance.General.Nonce,
			GeneralBalance:        balance.General.Balance.ToBigInt().Uint64(),
			EscrowBalanceActive:   balance.Escrow.Active.Balance.ToBigInt().Uint64(),
			EscrowBalanceShare:    balance.Escrow.Active.TotalShares.ToBigInt().Uint64(),
			EscrowDebondingActive: balance.Escrow.Debonding.Balance.ToBigInt().Uint64(),
			EscrowDebondingShare:  balance.Escrow.Debonding.TotalShares.ToBigInt().Uint64(),
		}

		i++
	}

	err = cli.DAO.CreateAccountBalances(balances)
	if err != nil {
		return err
	}

	txs := make([]dmodels.Transaction, 0, len(gen.Staking.Delegations)) //+len(gen.Staking.DebondingDelegations)

	//Genesis delegations
	for delegateAddress, receiver := range gen.Staking.Delegations {

		for receiverAddress, share := range receiver {
			txHash := sha256.Sum256([]byte(fmt.Sprint(gen.ChainID, "delegation", delegateAddress, receiver, share.Shares.String())))

			txs = append(txs, dmodels.Transaction{
				BlockLevel:          genesisHeight,
				BlockHash:           hex.EncodeToString(genesisBlockHash[:]),
				Hash:                hex.EncodeToString(txHash[:]),
				Time:                gen.GenesisTime,
				Amount:              0,
				EscrowAmount:        share.Shares.ToBigInt().Uint64(),
				EscrowReclaimAmount: 0,
				Type:                "addescrow",
				Sender:              delegateAddress.String(),
				Receiver:            receiverAddress,
				Nonce:               0,
				Fee:                 0,
				GasLimit:            0,
				GasPrice:            0,
			})
		}
	}

	//In this genesis not used
	//Genesis escrowreclaim
	//for debonder, staker := range gen.Staking.DebondingDelegations {
	//
	//	for staker, shareArr := range staker {
	//
	//		for i := range shareArr {
	//
	//			txHash := sha256.Sum256([]byte(fmt.Sprint(gen.ChainID, "reclaim", debonder, staker, shareArr[i].Shares.String())))
	//
	//			txs = append(txs, dmodels.Transaction{
	//				BlockLevel:          genesisHeight,
	//				BlockHash:           hex.EncodeToString(genesisBlockHash[:]),
	//				Hash:                hex.EncodeToString(txHash[:]),
	//				Time:                gen.GenesisTime,
	//				Amount:              0,
	//				EscrowAmount:        0,
	//				EscrowReclaimAmount: shareArr[i].Shares.ToBigInt().Uint64(),
	//				EscrowAccount:       staker,
	//				Type:                "reclaimescrow",
	//				Sender:              debonder,
	//				Receiver:            (api.Address)(oasis.SystemAddress).String(),
	//				Nonce:               0,
	//				Fee:                 0,
	//				GasLimit:            0,
	//				GasPrice:            0,
	//			})
	//
	//		}
	//	}
	//}

	err = cli.DAO.CreateTransfers(txs)
	if err != nil {
		return err
	}

	nodes := make([]dmodels.NodeRegistryTransaction, len(gen.Registry.Nodes))
	//Genesis nodes
	for i := range gen.Registry.Nodes {
		node := oasis.RegisterNode{}
		err = cbor.Unmarshal(gen.Registry.Nodes[i].UntrustedRawValue, &node)
		if err != nil {
			return err
		}

		txHash := sha256.Sum256([]byte(fmt.Sprint(gen.ChainID, "registernode", node.ID.String())))

		consensusIDBytes, err := node.Consensus.ID.MarshalBinary()
		if err != nil {
			return err
		}

		var physicalAddress string
		if len(node.Consensus.Addresses) > 0 {
			physicalAddress = node.Consensus.Addresses[0].Address.String()
		}

		nodes[i] = dmodels.NodeRegistryTransaction{
			BlockLevel:       genesisHeight,
			Hash:             hex.EncodeToString(txHash[:]),
			Time:             gen.GenesisTime,
			ID:               node.ID.String(),
			Address:          api.NewAddress(node.ID).String(),
			EntityID:         node.EntityID.String(),
			EntityAddress:    api.NewAddress(node.EntityID).String(),
			Expiration:       node.Expiration,
			P2PID:            node.P2P.ID.String(),
			ConsensusID:      node.Consensus.ID.String(),
			ConsensusAddress: crypto.AddressHash(consensusIDBytes).String(),
			PhysicalAddress:  physicalAddress,
			Roles:            uint32(node.Roles),
		}
	}

	err = cli.DAO.CreateRegisterNodeTransactions(nodes)
	if err != nil {
		return err
	}

	entities := make([]dmodels.EntityRegistryTransaction, len(gen.Registry.Entities))
	//Genesis nodes
	for i := range gen.Registry.Entities {
		entity := oasis.RegisterEntity{}
		err = cbor.Unmarshal(gen.Registry.Entities[i].UntrustedRawValue, &entity)
		if err != nil {
			return err
		}

		txHash := sha256.Sum256([]byte(fmt.Sprint(gen.ChainID, "registerentity", entity.ID.String())))

		nodes := make([]string, len(entity.Nodes))
		for i := range entity.Nodes {
			nodes[i] = entity.Nodes[i].String()
		}

		entities[i] = dmodels.EntityRegistryTransaction{
			BlockLevel:             genesisHeight,
			Hash:                   hex.EncodeToString(txHash[:]),
			Time:                   gen.GenesisTime,
			ID:                     entity.ID.String(),
			Address:                api.NewAddress(entity.ID).String(),
			Nodes:                  nodes,
			AllowEntitySignedNodes: entity.AllowEntitySignedNodes,
		}
	}

	err = cli.DAO.CreateRegisterEntityTransactions(entities)
	if err != nil {
		return err
	}

	log.Info("Genesis json were successfully parsed")

	return nil
}
