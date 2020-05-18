package dmodels

import (
	"time"
)

const (
	BlocksTable          = "blocks"
	BlockSignaturesTable = "block_signatures"
)

type Block struct {
	Height          uint64
	Hash            string
	CreatedAt       time.Time
	Epoch           uint64
	ProposerAddress string
	ValidatorHash   string
}

type BlockSignature struct {
	BlockHeight      int64
	Timestamp        time.Time
	BlockIDFlag      uint64
	ValidatorAddress string
	Signature        string
}
