package oasis

import (
	"github.com/oasisprotocol/oasis-core/go/common/crypto/address"
	"github.com/oasisprotocol/oasis-core/go/common/crypto/signature"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/bytes"
	tmtypes "github.com/tendermint/tendermint/types"
	"time"
)

const EpochBlocksNumber = 600

//AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=
var SystemPublicKey = signature.PublicKey{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var SystemAddress = address.Address{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

type Block struct {
	Hash       bytes.HexBytes  `cbor:"-"`
	Header     tmtypes.Header  `cbor:"header"`
	LastCommit BlockLastCommit `cbor:"last_commit"`
}

func (b Block) IsEpochBlock() bool {
	return (b.Header.Height % EpochBlocksNumber) == 0
}

type BlockLastCommit struct {
	Round      uint64      `cbor:"round"`
	Height     uint64      `cbor:"height"`
	BlockID    BlockID     `cbor:"block_id"`
	Signatures []Signature `cbor:"signatures"`
}

type Signature struct {
	Timestamp        time.Time      `cbor:"timestamp"`
	BlockIDFlag      uint64         `cbor:"block_id_flag"`
	ValidatorAddress crypto.Address `cbor:"validator_address"`
	Signature        []byte         `cbor:"signature"`
}

type BlockID struct {
	Hash  bytes.HexBytes `cbor:"hash"`
	Parts Parts          `cbor:"parts"`
}

type BlockHeader struct {
	ChainID            string         `cbor:"chain_id"`
	EvidenceHash       bytes.HexBytes `cbor:"evidence_hash"`
	ConsensusHash      bytes.HexBytes `cbor:"consensus_hash"`
	LastCommitHash     bytes.HexBytes `cbor:"last_commit_hash"`
	NextValidatorsHash bytes.HexBytes `cbor:"next_validators_hash"`
	Height             int64          `cbor:"height"`
	AppHash            bytes.HexBytes `cbor:"app_hash"`
	Time               time.Time      `cbor:"time"`
	ValidatorsHash     bytes.HexBytes `cbor:"validators_hash"`
	ProposerAddress    crypto.Address `cbor:"proposer_address"`
	DataHash           bytes.HexBytes `cbor:"data_hash"`
	LastResultsHash    bytes.HexBytes `cbor:"last_results_hash"`

	NumTxs   int64 `cbor:"num_txs"`
	TotalTxs int64 `cbor:"total_txs"`

	LastBlockID LastBlockID `cbor:"last_block_id"`
	Version     Version     `cbor:"version"`
}

type LastBlockID struct {
	Hash  bytes.HexBytes `cbor:"hash"`
	Parts Parts          `cbor:"parts"`
}

type Parts struct {
	Hash  bytes.HexBytes `cbor:"hash"`
	Total int            `cbor:"total"`
}

type Version struct {
	App   uint64 `cbor:"app"`
	Block uint64 `cbor:"block"`
}
