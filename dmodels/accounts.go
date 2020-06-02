package dmodels

import "time"

const (
	AccountBalanceTable          = "oasis.account_balance"
	DayTotalBalanceView          = "oasis.day_total_balance_view"
	TopEscrowBalanceAccountsView = "oasis.top_escrow_balance_accounts_view"
)

type AccountTime struct {
	CreatedAt  time.Time `db:"created_at"`
	LastActive time.Time `db:"last_active"`
}

type AccountBalance struct {
	Account               string    `db:"acb_account"`
	Time                  time.Time `db:"blk_time"`
	Height                int64     `db:"blk_lvl"`
	Nonce                 uint64    `db:"acb_nonce"`
	GeneralBalance        string    `db:"acb_general_balance"`
	EscrowBalanceActive   string    `db:"acb_escrow_balance_active"`
	EscrowBalanceShare    string    `db:"acb_escrow_balance_share"`
	EscrowDebondingActive string    `db:"acb_escrow_debonding_active"`
	EscrowDebondingShare  string    `db:"acb_escrow_debonding_share"`
}
