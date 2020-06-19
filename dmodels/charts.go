package dmodels

import (
	"time"
)

type ChartData struct {
	BeginOfPeriod     time.Time `db:"start_of_period"`
	TransactionVolume string    `db:"transaction_volume"`
	EscrowRatio       float64   `db:"escrow_ratio"`
}

type BalanceChartData struct {
	BeginOfPeriod    time.Time `db:"start_of_period"`
	AccountID        string    `db:"acb_account"`
	GeneralBalance   uint64    `json:"acb_general_balance"`
	EscrowBalance    uint64    `json:"escrow_balance_active"`
	DebondingBalance uint64    `json:"escrow_debonding_active"`

	//Not implemented yet
	SelfStakeBalance uint64 `json:"self_stake_balance"`
}
