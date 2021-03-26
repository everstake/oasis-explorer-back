package dmodels

import (
	"fmt"
	"strings"
	"time"
)

const (
	TransactionsTable   = "transactions"
	RegisterNodeTable   = "register_node_transactions"
	RegisterEntityTable = "register_entity_transactions"
	Precision           = 9
)

type TransactionMethod struct {
	api    string
	method string
}

type TransactionType string

func NewTransactionType(s string) (tt TransactionMethod, err error) {
	a := strings.Split(s, ".")

	if len(a) != 2 {
		return tt, fmt.Errorf("Wrong input TransactionType")
	}

	return TransactionMethod{
		api:    a[0],
		method: a[1],
	}, nil
}

func (tt TransactionMethod) Type() TransactionType {
	return TransactionType(strings.ToLower(tt.method))
}

const (
	TransactionTypeTransfer       TransactionType = "transfer"
	TransactionTypeBurn           TransactionType = "burn"
	TransactionTypeAddEscrow      TransactionType = "addescrow"
	TransactionTypeReclaimEscrow  TransactionType = "reclaimescrow"
	TransactionTypeRegisterNode   TransactionType = "registernode"
	TransactionTypeRegisterEntity TransactionType = "registerentity"
)

type Transaction struct {
	BlockLevel          uint64          `db:"blk_lvl"`
	BlockHash           string          `db:"blk_hash"`
	Hash                string          `db:"tx_hash"`
	Time                time.Time       `db:"tx_time"`
	Amount              uint64          `db:"tx_amount"`
	EscrowAmount        uint64          `db:"tx_escrow_amount"`
	EscrowReclaimAmount uint64          `db:"tx_escrow_reclaim_amount"`
	Type                TransactionType `db:"tx_type"`
	Status              bool            `db:"tx_status"`
	Error               string          `db:"tx_error"`
	Sender              string          `db:"tx_sender"`
	Receiver            string          `db:"tx_receiver"`
	Nonce               uint64          `db:"tx_nonce"`
	Fee                 uint64          `db:"tx_fee"`
	GasLimit            uint64          `db:"tx_gas_limit"` //Probably GasUsed
	GasPrice            uint64          `db:"tx_gas_price"`
}

type NodeRegistryTransaction struct {
	BlockLevel       uint64    `db:"blk_lvl"`
	Hash             string    `db:"tx_hash"`
	Time             time.Time `db:"tx_time"`
	ID               string
	Address          string
	EntityID         string
	EntityAddress    string
	Expiration       uint64
	P2PID            string
	ConsensusID      string
	ConsensusAddress string
	PhysicalAddress  string
	Roles            uint32
}

type EntityRegistryTransaction struct {
	BlockLevel uint64    `db:"blk_lvl"`
	Hash       string    `db:"tx_hash"`
	Time       time.Time `db:"tx_time"`
	ID         string
	Address    string
	Nodes      []string
}
