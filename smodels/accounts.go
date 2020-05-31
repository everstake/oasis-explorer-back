package smodels

import (
	"github.com/wedancedalot/decimal"
	"oasisTracker/dmodels"
	"sync"
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

type AccountsContainer struct {
	balances []dmodels.AccountBalance
	mu       *sync.Mutex
}

func NewAccountsContainer() *AccountsContainer {
	return &AccountsContainer{
		mu:       &sync.Mutex{},
		balances: []dmodels.AccountBalance{},
	}
}

func (c *AccountsContainer) Add(balances []dmodels.AccountBalance) {
	if len(balances) == 0 {
		return
	}

	c.mu.Lock()
	c.balances = append(c.balances, balances...)
	c.mu.Unlock()
}

func (c *AccountsContainer) Balances() []dmodels.AccountBalance {
	return c.balances
}

func (c *AccountsContainer) IsEmpty() bool {
	return len(c.balances) == 0
}

func (c *AccountsContainer) Flush() {
	c.balances = []dmodels.AccountBalance{}
}
