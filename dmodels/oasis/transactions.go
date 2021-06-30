package oasis

import (
	"github.com/oasisprotocol/oasis-core/go/common/crypto/signature"
	"github.com/oasisprotocol/oasis-core/go/common/entity"
	"github.com/oasisprotocol/oasis-core/go/common/node"
	registry "github.com/oasisprotocol/oasis-core/go/registry/api"
)

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
