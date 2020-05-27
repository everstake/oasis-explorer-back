package smodels

import (
	"github.com/wedancedalot/decimal"
	"time"
)

type Account struct {
	Address          string          `json:"address"`
	LiquidBalance    decimal.Decimal `json:"liquid_balance"`
	EscrowBalance    decimal.Decimal `json:"escrow_balance"`
	DebondingBalance decimal.Decimal `json:"debonding_balance"`
	TotalBalance     decimal.Decimal `json:"total_balance"`
	CreatedAt        time.Time       `json:"created_at"`
	LastActive       time.Time       `json:"last_active"`
	Nonce            uint64          `json:"nonce"`
	Type             string          `json:"type"`
	NodeAddress      string          `json:"node_address,omitempty"`
}
