package oasis

import (
	beacon "github.com/oasisprotocol/oasis-core/go/beacon/api"
	"github.com/oasisprotocol/oasis-core/go/common"
	"github.com/oasisprotocol/oasis-core/go/common/cbor"
	"github.com/oasisprotocol/oasis-core/go/common/crypto/pvss"
	"github.com/oasisprotocol/oasis-core/go/common/crypto/signature"
	"github.com/oasisprotocol/oasis-core/go/common/entity"
	"github.com/oasisprotocol/oasis-core/go/common/node"
	"github.com/oasisprotocol/oasis-core/go/common/quantity"
	tx "github.com/oasisprotocol/oasis-core/go/consensus/api/transaction"
	registry "github.com/oasisprotocol/oasis-core/go/registry/api"
	"github.com/oasisprotocol/oasis-core/go/roothash/api/commitment"
	scheduler "github.com/oasisprotocol/oasis-core/go/scheduler/api"
	staking "github.com/oasisprotocol/oasis-core/go/staking/api"
)

type UntrustedRawValue struct {
	Fee    tx.Fee `cbor:"fee"`
	Nonce  uint64 `cbor:"nonce"`
	Method string `cbor:"method"`
	Body   TxBody `cbor:"body"`
}

//TODO refactor txBody
type TxBody struct {
	//staking.Transfer
	To staking.Address `json:"to"`

	staking.Allow

	//staking.Withdraw
	From staking.Address `json:"from"`

	//staking.Transfer staking.Burn staking.Escrow staking.Withdraw
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

	//ExecutorCommit registry.Runtime
	ID common.Namespace `json:"id"`

	//ExecutorCommit
	Commits []commitment.ExecutorCommitment `json:"commits"`

	//ExecutorCommitExpired,PVSSCommit,PVSSReveal
	Round uint64 `json:"round"`

	//PVSSCommit,PVSSReveal
	Epoch beacon.EpochTime `json:"epoch"`
	Index int              `json:"index"`

	//PVSSCommit
	Commit *pvss.Commit `json:"commit,omitempty"`

	// PVSSReveal
	Reveal *pvss.Reveal `json:"reveal,omitempty"`

	//PVSSEvent
	Height       int64                 `json:"height,omitempty"`
	State        beacon.RoundState     `json:"state,omitempty"`
	Participants []signature.PublicKey `json:"participants,omitempty"`

	//registry.Runtime
	cbor.Versioned
	EntityID        signature.PublicKey                                                           `json:"entity_id"`
	Genesis         registry.RuntimeGenesis                                                       `json:"genesis"`
	Kind            registry.RuntimeKind                                                          `json:"kind"`
	TEEHardware     node.TEEHardware                                                              `json:"tee_hardware"`
	Version         registry.VersionInfo                                                          `json:"versions"`
	KeyManager      *common.Namespace                                                             `json:"key_manager,omitempty"`
	Executor        registry.ExecutorParameters                                                   `json:"executor,omitempty"`
	TxnScheduler    registry.TxnSchedulerParameters                                               `json:"txn_scheduler,omitempty"`
	Storage         registry.StorageParameters                                                    `json:"storage,omitempty"`
	AdmissionPolicy registry.RuntimeAdmissionPolicy                                               `json:"admission_policy"`
	Constraints     map[scheduler.CommitteeKind]map[scheduler.Role]registry.SchedulingConstraints `json:"constraints,omitempty"`
	Staking         registry.RuntimeStakingParameters                                             `json:"staking,omitempty"`
	GovernanceModel registry.RuntimeGovernanceModel                                               `json:"governance_model"`
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
