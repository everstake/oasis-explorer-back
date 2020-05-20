package oasis

import (
	"github.com/oasislabs/oasis-core/go/common/crypto/signature"
	"github.com/oasislabs/oasis-core/go/common/entity"
	"github.com/oasislabs/oasis-core/go/common/node"
	"github.com/oasislabs/oasis-core/go/common/quantity"
	tx "github.com/oasislabs/oasis-core/go/consensus/api/transaction"
	registry "github.com/oasislabs/oasis-core/go/registry/api"
	staking "github.com/oasislabs/oasis-core/go/staking/api"
)

type TxRaw struct {
	UntrustedRawValue []byte              `json:"untrusted_raw_value"`
	Signature         signature.Signature `cbor:"signature"`
}

type UntrustedRawValue struct {
	Fee    tx.Fee `cbor:"fee"`
	Nonce  uint64 `cbor:"nonce"`
	Method string `cbor:"method"`
	Body   TxBody `cbor:"body"`
}

type TxBody struct {
	staking.Transfer
	staking.Burn
	staking.AmendCommissionSchedule
	//staking.Escrow staking.ReclaimEscrow
	EscrowTx

	// RegisterEntity RegisterRuntime RegisterNode
	RegisterTx

	//UnfreezeNode
	registry.UnfreezeNode
}

type EscrowTx struct {
	Account signature.PublicKey `json:"escrow_account"`
	//Escrow
	Tokens quantity.Quantity `json:"escrow_tokens"`
	//ReclaimEscrow
	Shares quantity.Quantity `json:"reclaim_shares"`
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
