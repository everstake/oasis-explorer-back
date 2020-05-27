package dmodels

import (
	"fmt"
	"strings"
	"time"
)

const (
	TransactionsTable = "oasis.transactions"
	Precision         = 9
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
	TransactionTypeTransfer      TransactionType = "Transfer"
	TransactionTypeBurn          TransactionType = "Burn"
	TransactionTypeAddEscrow     TransactionType = "AddEscrow"
	TransactionTypeReclaimEscrow TransactionType = "ReclaimEscrow"
)

type Transaction struct {
	BlockLevel          uint64          `db:"blk_lvl"`
	BlockHash           string          `db:"blk_hash"`
	Hash                string          `db:"tx_hash"`
	Time                time.Time       `db:"tx_time"`
	Amount              string          `db:"tx_amount"`
	EscrowAmount        string          `db:"tx_escrow_amount"`
	EscrowReclaimAmount string          `db:"tx_escrow_reclaim_amount"`
	EscrowAccount       string          `db:"tx_escrow_account"`
	Type                TransactionType `db:"tx_type"`
	Sender              string          `db:"tx_sender"`
	Receiver            string          `db:"tx_receiver"`
	Nonce               uint64          `db:"tx_nonce"`
	Fee                 uint64          `db:"tx_fee"`
	GasLimit            uint64          `db:"tx_gas_limit"` //Probably GasUsed
	GasPrice            uint64          `db:"tx_gas_price"`
}
