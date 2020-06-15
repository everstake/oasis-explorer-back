package smodels

import (
	"fmt"
	"github.com/oasislabs/oasis-core/go/staking/api"
	"oasisTracker/dmodels"
	"sort"
	"sync"
	"time"
)

func NewAccountListParams() AccountListParams {
	return AccountListParams{
		CommonParams: CommonParams{
			Limit:  50,
			Offset: 0,
		},
		SortColumn: sortCreatedAt,
		SortSide:   sortDesc,
	}
}

const (
	sortCreatedAt        = "created_at"
	sortBalance          = "general_balance"
	sortEscrowBalance    = "escrow_balance"
	sortEscrowShare      = "escrow_share"
	sortOperationsAmount = "operations_amount"
	sortAsc              = "asc"
	sortDesc             = "desc"
	AccountTypeAccount   = "account"
	AccountTypeNode      = "node"
	AccountTypeEntity    = "entity"
)

//Sorted
var sortColumns = []string{sortCreatedAt, sortEscrowBalance, sortEscrowShare, sortBalance, sortOperationsAmount}
var sortSides = []string{sortAsc, sortDesc}

func (b AccountListParams) Validate() error {

	if err := b.CommonParams.Validate(); err != nil {
		return err
	}

	if b.SortColumn == "" {
		return fmt.Errorf("sort column not present")
	}

	//Not found
	index := sort.SearchStrings(sortColumns, b.SortColumn)
	if index == len(sortColumns) || sortColumns[index] != b.SortColumn {
		return fmt.Errorf("sort column unknown")
	}

	if b.SortSide == "" {
		return fmt.Errorf("sort side not present")
	}

	//Not found
	index = sort.SearchStrings(sortSides, b.SortSide)
	if index == len(sortSides) || sortSides[index] != b.SortSide {
		return fmt.Errorf("sort side unknown")
	}

	return nil
}

type AccountListParams struct {
	CommonParams
	SortColumn string `schema:"sort_column"`
	SortSide   string `schema:"sort_side"`
}

type AccountList struct {
	Account            string `json:"account_id"`
	CreatedAt          int64  `json:"created_at"`
	OperationsAmount   uint64 `json:"operations_amount"`
	GeneralBalance     uint64 `json:"general_balance"`
	EscrowBalance      uint64 `json:"escrow_balance"`
	EscrowBalanceShare uint64 `json:"escrow_balance_share"`
	Delegate           string `json:"delegate"`
	Type               string `json:"type"`
}

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

	EntityAddress string         `json:"entity_address,omitempty"`
	Validator     *ValidatorInfo `json:"validator"`
}

var TestNetGenesis = api.CommissionScheduleRules{
	RateChangeInterval: 1,
	RateBoundLead:      14,
	MaxRateSteps:       21,
	MaxBoundSteps:      21,
}

type ValidatorInfo struct {
	api.CommissionScheduleRules
	Status           string `json:"status"`
	NodeAddress      string `json:"node_address,omitempty"`
	ConsensusAddress string `json:"consensus_address,omitempty"`
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
