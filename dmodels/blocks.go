package dmodels

import (
	"time"
)

const (
	BlocksNewTable       = "blocks_new"
	BlocksOldTable       = "block_row_view"
	BlocksRowView        = "block_row_view"
	BlocksSigCountView   = "block_signatures_count_view"
	BlockSignaturesTable = "block_signatures"

	BlocksPostgresTable    = "blocks"
	BlocksDayPostgresTable = "day_blocks"
)

type Block struct {
	Height             uint64    `db:"blk_lvl"`
	Hash               string    `db:"blk_hash"`
	CreatedAt          time.Time `db:"blk_created_at"`
	Epoch              uint64    `db:"blk_epoch"`
	ProposerAddress    string    `db:"blk_proposer_address"`
	ValidatorHash      string    `db:"blk_validator_hash"`
	NumberOfTxs        uint64    `db:"blk_number_of_txs"`
	NumberOfSignatures uint64    `db:"blk_number_of_signatures"`
	Fees               uint64    `db:"blk_fees"`
	GasUsed            uint64    `db:"blk_gas_used"`
}

type BlockSignature struct {
	BlockHeight      int64
	Timestamp        time.Time
	BlockIDFlag      uint64
	ValidatorAddress string
	Signature        string
}

type BlockInfo struct {
	ID          uint64 `gorm:"column:id;PRIMARY_KEY"`
	TotalBlocks uint64 `gorm:"column:total_count"`
	LastLvl     uint64 `gorm:"column:last_lvl"`
	Epoch       uint64 `gorm:"column:epoch"`
}

type BlockDayInfo struct {
	ID             uint64    `gorm:"column:id;PRIMARY_KEY"`
	TotalDayBlocks uint64    `gorm:"column:day_total_count"`
	Day            time.Time `gorm:"column:day"`
}
