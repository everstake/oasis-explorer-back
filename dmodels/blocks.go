package dmodels

import (
	"time"
)

const (
	BlocksTable          = "blocks"
	BlocksRowView        = "block_row_view"
	BlocksSigCountView   = "blocks_sig_count"
	BlockSignaturesTable = "block_signatures"
)

type RowBlock struct {
	Block
	GasUsed  uint64 `db:"gas_used"`
	Fee      uint64 `db:"fee"`
	TxsCount uint64 `db:"txs_count"`
	SigCount uint64 `db:"sig_count"`
}

type Block struct {
	Height          uint64    `db:"blk_lvl"`
	Hash            string    `db:"blk_hash"`
	CreatedAt       time.Time `db:"blk_created_at"`
	Epoch           uint64    `db:"blk_epoch"`
	ProposerAddress string    `db:"blk_proposer_address"`
	ValidatorHash   string    `db:"blk_validator_hash"`
}

type BlockSignature struct {
	BlockHeight      int64
	Timestamp        time.Time
	BlockIDFlag      uint64
	ValidatorAddress string
	Signature        string
}
