package smodels

//Constructor to setup default values
func NewBlockParams() BlockParams {
	return BlockParams{
		Limit: 20,
	}
}

func (b *BlockParams) Validate() error {
	return nil
}

type BlockParams struct {
	Limit      uint64
	Offset     uint64
	BlockID    []string `schema:"block_id"`
	BlockLevel []int64  `schema:"block_level"`
	Sender     string
	Receiver   string
	//Time range
	From uint64
	To   uint64
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
