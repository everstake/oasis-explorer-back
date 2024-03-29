package smodels

import (
	"time"

	"github.com/oasisprotocol/oasis-core/go/common/quantity"
	"github.com/oasisprotocol/oasis-core/go/consensus/api/transaction"
	"github.com/oasisprotocol/oasis-core/go/staking/api"
)

type GenesisDocument struct {
	GenesisTime   time.Time `json:"genesis_time"`
	GenesisHeight uint64    `json:"height"`
	ChainID       string    `json:"chain_id"`
	Registry      Registry  `json:"registry"`
	Staking       Staking   `json:"staking"`
	Beacon        Beacon    `json:"beacon"`
}

type Registry struct {
	Entities []transaction.SignedTransaction `json:"entities"`
	Nodes    []transaction.SignedTransaction `json:"nodes"`
}
type Staking struct {
	Ledger               map[api.Address]*api.Account                                  `json:"ledger"`
	Delegations          map[api.Address]map[api.Address]GenesisDelegation             `json:"delegations"`
	DebondingDelegations map[api.Address]map[api.Address][]GenesisDebondingDelegations `json:"debonding_delegations"`
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

type Beacon struct {
	Base   uint64      `json:"base"`   //base epoch
	Params interface{} `json:"params"` // Todo add params definition
}
