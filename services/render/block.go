package render

import (
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
)

func Blocks(bs []dmodels.Block) []smodels.Block {
	blocks := make([]smodels.Block, len(bs))
	for i := range bs {
		blocks[i] = Block(bs[i])
	}
	return blocks
}

func Block(b dmodels.Block) smodels.Block {

	return smodels.Block{
		Epoch:              b.Epoch,
		Fees:               b.Fees,
		GasUsed:            b.GasUsed,
		Hash:               b.Hash,
		Level:              b.Height,
		NumberOfTxs:        b.NumberOfTxs,
		NumberOfSignatures: b.NumberOfSignatures,
		Proposer:           b.ProposerAddress,
		Timestamp:          b.CreatedAt.Unix(),
	}
}
