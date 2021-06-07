package oasis

import (
	beacon "github.com/oasisprotocol/oasis-core/go/beacon/api"
	"github.com/oasisprotocol/oasis-core/go/common/crypto/pvss"
	"github.com/oasisprotocol/oasis-core/go/common/crypto/signature"
	"github.com/oasisprotocol/oasis-core/go/common/entity"
	"github.com/oasisprotocol/oasis-core/go/common/node"
	"github.com/oasisprotocol/oasis-core/go/common/quantity"
	tx "github.com/oasisprotocol/oasis-core/go/consensus/api/transaction"
	registry "github.com/oasisprotocol/oasis-core/go/registry/api"
	roothash "github.com/oasisprotocol/oasis-core/go/roothash/api"
	staking "github.com/oasisprotocol/oasis-core/go/staking/api"
)

type UntrustedRawValue struct {
	Fee    tx.Fee `cbor:"fee"`
	Nonce  uint64 `cbor:"nonce"`
	Method string `cbor:"method"`
	Body   TxBody `cbor:"body"`
}

type TxBody struct {
	//staking.Transfer
	To staking.Address `json:"to"`

	//staking.Transfer staking.Burn staking.Escrow
	Amount quantity.Quantity `json:"amount"`

	staking.AmendCommissionSchedule

	//staking.Escrow staking.ReclaimEscrow
	Account staking.Address `json:"account"`

	//staking.ReclaimEscrow
	Shares quantity.Quantity `json:"shares"`

	// RegisterEntity RegisterRuntime RegisterNode
	RegisterTx

	//UnfreezeNode
	registry.UnfreezeNode

	//ExecutorCommit
	roothash.ExecutorCommit

	//ExecutorCommitExpired,PVSSCommit,PVSSReveal
	Round uint64 `json:"round"`

	//PVSSCommit,PVSSReveal
	Epoch beacon.EpochTime `json:"epoch"`

	//PVSSCommit
	Commit *pvss.Commit `json:"commit,omitempty"`

	// PVSSReveal
	Reveal *pvss.Reveal `json:"reveal,omitempty"`
}

type RegisterTx struct {
	// Blob is the signed blob.
	Blob []byte `json:"untrusted_raw_value"`

	// RegisterEntity RegisterRuntime
	// Signature is the signature over blob.
	Signature signature.Signature `json:"signature"`

	//RegisterNode
	// Signatures are the signatures over the blob.
	Signatures []signature.Signature `json:"signatures"`
}

type RegisterNode struct {
	node.Node
}

type RegisterRuntime struct {
	registry.Runtime
}

type RegisterEntity struct {
	entity.Entity
}
