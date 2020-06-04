package smodels

import (
	"github.com/oasislabs/oasis-core/go/staking/api"
	"oasisTracker/dmodels"
	"sync"
	"time"
)

//const
type Account struct {
	Address          string    `json:"address"`
	LiquidBalance    uint64    `json:"liquid_balance"`
	EscrowBalance    uint64    `json:"escrow_balance"`
	DebondingBalance uint64    `json:"debonding_balance"`
	TotalBalance     uint64    `json:"total_balance"`
	CreatedAt        time.Time `json:"created_at"`
	LastActive       time.Time `json:"last_active"`
	Nonce            *uint64   `json:"nonce"`
	Type             string    `json:"type"`

	EntityAddress string     `json:"entity_address,omitempty"`
	Validator     *Validator `json:"validator"`
}

var TestNetGenesis = api.CommissionScheduleRules{
	RateChangeInterval: 1,
	RateBoundLead:      14,
	MaxRateSteps:       21,
	MaxBoundSteps:      21,
}

type Validator struct {
	api.CommissionScheduleRules
	Status           string `json:"status"`
	NodeAddress      string `json:"node_address,omitempty"`
	ConsensusAddress string `json:"consensus_address,"`
	DepositorsCount  uint64 `json:"depositors_count,omitempty"`
	BlocksCount      uint64 `json:"blocks_count"`
	SignaturesCount  uint64 `json:"signatures_count"`
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
