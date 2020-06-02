package smodels

type TransactionsParams struct {
	Limit         uint64
	Offset        uint64
	BlockID       []string `schema:"block_id"`
	BlockLevel    []int64  `schema:"block_level"`
	OperationID   []string `schema:"operation_id"`
	OperationKind []string `schema:"operation_kind"`
}

type Transaction struct {
	Amount    uint64 `json:"amount,omitempty"`
	Fee       uint64 `json:"fee,omitempty"`
	From      string `json:"from,omitempty"`
	GasPrice  uint64 `json:"gas_price,omitempty"`
	GasUsed   uint64 `json:"gas_used,omitempty"`
	Hash      string `json:"hash,omitempty"`
	Level     uint64 `json:"level,omitempty"`
	Nonce     uint64 `json:"nonce,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	To        string `json:"to,omitempty"`
	Type      string `json:"type,omitempty"`
}
