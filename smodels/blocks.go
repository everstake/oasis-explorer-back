package smodels

type BlockParams struct {
	Limit      uint64
	Offset     uint64
	BlockID    []string `schema:"block_id"`
	BlockLevel []int64  `schema:"block_level"`
}

type Block struct {
	Epoch     uint64 `json:"epoch,omitempty"`
	Hash      string `json:"hash"`
	Level     uint64 `json:"level"`
	Proposer  string `json:"proposer,omitempty"`
	Timestamp int64  `json:"timestamp"`

	NumberOfTxs        uint64 `json:"number_of_txs,omitempty"`
	NumberOfSignatures uint64 `json:"number_of_signatures,omitempty"`
	Fees               uint64 `json:"fees,omitempty"`
	GasUsed            uint64 `json:"gas_used,omitempty"`
}
