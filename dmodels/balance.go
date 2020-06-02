package dmodels

import "time"

type DayBalance struct {
	StartOfPeriod         time.Time `db:"start_of_period"`
	GeneralBalance        uint64    `db:"general_balance"`
	EscrowBalanceActive   uint64    `db:"escrow_balance_active"`
	EscrowDebondingActive uint64    `db:"escrow_debonding_active"`
}
