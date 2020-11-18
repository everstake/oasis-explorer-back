package smodels

import (
	"github.com/oasisprotocol/oasis-core/go/common/quantity"
	"github.com/oasisprotocol/oasis-core/go/consensus/api/transaction"
	"github.com/oasisprotocol/oasis-core/go/staking/api"
	"time"
)

type GenesisDocument struct {
	GenesisTime time.Time    `json:"genesis_time"`
	EpochTime   GenesisEpoch `json:"epochtime"`
	ChainID     string       `json:"chain_id"`
	Registry    Registry     `json:"registry"`
	Staking     Staking      `json:"staking"`
}

type Registry struct {
	Entities []transaction.SignedTransaction `json:"entities"`
	Nodes    []transaction.SignedTransaction `json:"nodes"`
}
type Staking struct {
	Ledger               map[api.Address]api.Account                         `json:"ledger"`
	Delegations          map[api.Address]map[string]GenesisDelegation        `json:"delegations"`
	DebondingDelegations map[string]map[string][]GenesisDebondingDelegations `json:"debonding_delegations"`
}

type GenesisEpoch struct {
	Base uint64 `json:"base"`
}

type GenesisDelegation struct {
	Shares quantity.Quantity `json:"shares"`
}

type GenesisDebondingDelegations struct {
	Shares    quantity.Quantity `json:"shares"`
	DebondEnd uint64            `json:"debond_end"`
}
