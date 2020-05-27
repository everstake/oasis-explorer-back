package render

import (
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
)

func Blocks(bs []dmodels.RowBlock) []smodels.Block {
	blocks := make([]smodels.Block, len(bs))
	for i := range bs {
		blocks[i] = Block(bs[i])
	}
	return blocks
}

func Block(b dmodels.RowBlock) smodels.Block {

	return smodels.Block{
		Epoch:              b.Epoch,
		Fees:               b.Fee,
		GasUsed:            b.GasUsed,
		Hash:               b.Hash,
		Level:              b.Height,
		NumberOfTxs:        b.TxsCount,
		NumberOfSignatures: b.SigCount,
		Proposer:           b.ProposerAddress,
		Timestamp:          b.CreatedAt.Unix(),
	}
}
