package dmodels

import "time"

const (
	AccountBalanceTable = "oasis.account_balance"
	DayTotalBalanceView = "oasis.day_total_balance_view"
)

type AccountTime struct {
	CreatedAt  time.Time `db:"created_at"`
	LastActive time.Time `db:"last_active"`
}

type AccountBalance struct {
	Account               string
	Time                  time.Time
	Height                int64
	GeneralBalance        string
	Nonce                 uint64
	EscrowBalanceActive   string
	EscrowBalanceShare    string
	EscrowDebondingActive string
	EscrowDebondingShare  string
}
