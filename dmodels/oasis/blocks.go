package oasis

import (
	"github.com/tendermint/tendermint/crypto"
	cmn "github.com/tendermint/tendermint/libs/common"
	"time"
)

type Block struct {
	Hash       []byte          `cbor:"-"`
	Header     BlockHeader     `cbor:"header"`
	LastCommit BlockLastCommit `cbor:"last_commit"`
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
	Hash  cmn.HexBytes `cbor:"hash"`
	Parts Parts        `cbor:"parts"`
}

type BlockHeader struct {
	ChainID            string         `cbor:"chain_id"`
	EvidenceHash       cmn.HexBytes   `cbor:"evidence_hash"`
	ConsensusHash      cmn.HexBytes   `cbor:"consensus_hash"`
	LastCommitHash     cmn.HexBytes   `cbor:"last_commit_hash"`
	NextValidatorsHash cmn.HexBytes   `cbor:"next_validators_hash"`
	Height             int64          `cbor:"height"`
	AppHash            cmn.HexBytes   `cbor:"app_hash"`
	Time               time.Time      `cbor:"time"`
	ValidatorsHash     cmn.HexBytes   `cbor:"validators_hash"`
	ProposerAddress    crypto.Address `cbor:"proposer_address"`
	DataHash           cmn.HexBytes   `cbor:"data_hash"`
	LastResultsHash    cmn.HexBytes   `cbor:"last_results_hash"`

	NumTxs   int64 `cbor:"num_txs"`
	TotalTxs int64 `cbor:"total_txs"`

	LastBlockID LastBlockID `cbor:"last_block_id"`
	Version     Version     `cbor:"version"`
}

type LastBlockID struct {
	Hash  cmn.HexBytes `cbor:"hash"`
	Parts Parts        `cbor:"parts"`
}

type Parts struct {
	Hash  cmn.HexBytes `cbor:"hash"`
	Total int          `cbor:"total"`
}

type Version struct {
	App   uint64 `cbor:"app"`
	Block uint64 `cbor:"block"`
}
