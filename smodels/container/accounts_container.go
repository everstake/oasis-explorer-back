package container

import (
	"oasisTracker/dmodels"
	"sync"
)

type AccountsContainer struct {
	balances []dmodels.AccountBalance
	mu       *sync.Mutex
}

func NewAccountsContainer() *AccountsContainer {
	return &AccountsContainer{
		mu: &sync.Mutex{},
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
