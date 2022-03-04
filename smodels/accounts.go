package smodels

import (
	"fmt"
	"sort"
	"time"

	"github.com/oasisprotocol/oasis-core/go/staking/api"
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
	sortCreatedAt                   = "created_at"
	sortDelegationsBalance          = "delegations_balance"
	sortDebondingDelegationsBalance = "debonding_delegations_balance"
	sortBalance                     = "general_balance"
	sortEscrowBalance               = "escrow_balance"
	sortEscrowBalanceShare          = "escrow_balance_share"
	sortOperationsAmount            = "operations_amount"
	sortOperationsNumber            = "operations_number"
	sortAsc                         = "asc"
	sortDesc                        = "desc"
	AccountTypeAccount              = "account"
	AccountTypeNode                 = "node"
	AccountTypeEntity               = "entity"
	AccountTypeValidator            = "validator"
)

//Sorted
var sortColumns = []string{sortCreatedAt, sortDebondingDelegationsBalance, sortDelegationsBalance, sortEscrowBalance, sortEscrowBalanceShare /*sortEscrowShare,*/, sortBalance, sortOperationsAmount, sortOperationsNumber}
var sortSides = []string{sortAsc, sortDesc}

func (b *AccountListParams) Validate() error {

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
	OperationsNumber   uint64 `json:"operations_number"`
	GeneralBalance     uint64 `json:"general_balance"`
	EscrowBalance      uint64 `json:"escrow_balance"`
	EscrowBalanceShare uint64 `json:"escrow_balance_share"`

	DelegationsBalance          uint64 `json:"delegations_balance"`
	DebondingDelegationsBalance uint64 `json:"debonding_delegations_balance"`
	SelfDelegationBalance       uint64 `json:"self_delegation_balance"`

	Delegate string `json:"delegate"`
	Type     string `json:"type"`
}

type Account struct {
	Address          string `json:"address"`
	LiquidBalance    uint64 `json:"liquid_balance"`
	EscrowBalance    uint64 `json:"escrow_balance"`
	DebondingBalance uint64 `json:"escrow_debonding_balance"`

	DelegationsBalance          uint64 `json:"delegations_balance"`
	DebondingDelegationsBalance uint64 `json:"debonding_delegations_balance"`
	SelfDelegationBalance       uint64 `json:"self_delegation_balance"`

	TotalBalance uint64    `json:"total_balance"`
	CreatedAt    time.Time `json:"created_at"`
	LastActive   time.Time `json:"last_active"`
	Nonce        *uint64   `json:"nonce"`
	Type         string    `json:"type"`

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
