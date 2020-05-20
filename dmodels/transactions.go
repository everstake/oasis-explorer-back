package dmodels

import (
	"github.com/wedancedalot/decimal"
	"time"
)

const (
	TransactionsTable = "transactions"
	Precision         = 9
)

type TransactionType string

const (
	TransactionTypeTransfer      = "Transfer"
	TransactionTypeBurn          = "Burn"
	TransactionTypeAddEscrow     = "AddEscrow"
	TransactionTypeReclaimEscrow = "ReclaimEscrow"
)

type Transaction struct {
	BlockLevel          uint64
	Hash                string
	Time                time.Time
	Amount              decimal.Decimal
	EscrowAmount        decimal.Decimal
	EscrowReclaimAmount decimal.Decimal
	EscrowAccount       string
	Type                TransactionType
	Sender              string
	Receiver            string
	Nonce               uint64
	Fee                 uint64
	GasLimit            uint64 //Probably GasUsed
	GasPrice            uint64
}
