package dmodels

import (
	"time"
)

const (
	TransactionsTable = "transactions"
)

type TransactionType string

const (
	TransactionTypeTransfer      = "Transfer"
	TransactionTypeBurn          = "Burn"
	TransactionTypeAddEscrow     = "AddEscrow"
	TransactionTypeReclaimEscrow = "ReclaimEscrow"
)

type Transaction struct {
	BlockLevel    uint64
	Hash          string
	Time          time.Time
	Amount        uint64
	EscrowAmount  uint64
	EscrowAccount string
	Type          TransactionType
	Sender        string
	Receiver      string
	Nonce         uint64
	Fee           uint64
	GasLimit      uint64 //Probably GasUsed
	GasPrice      uint64
}
