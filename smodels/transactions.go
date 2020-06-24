package smodels

//Constructor to setup default values
func NewTransactionsParams() TransactionsParams {
	return TransactionsParams{
		CommonParams: CommonParams{
			Limit: 20,
		},
	}
}

type TransactionsParams struct {
	CommonParams
	//Time range
	From uint64
	To   uint64

	BlockID       []string `schema:"block_id"`
	BlockLevel    []int64  `schema:"block_level"`
	OperationID   []string `schema:"operation_id"`
	OperationKind []string `schema:"operation_kind"`
	Sender        string
	Receiver      string
	//Get transaction by account
	AccountID string `schema:"account_id"`
}

type Transaction struct {
	Amount              uint64 `json:"amount"`
	EscrowAmount        uint64 `json:"escrow_amount,omitempty"`
	ReclaimEscrowAmount uint64 `json:"reclaim_escrow_amount,omitempty"`
	Fee                 uint64 `json:"fee,omitempty"`
	From                string `json:"from,omitempty"`
	GasPrice            uint64 `json:"gas_price,omitempty"`
	GasUsed             uint64 `json:"gas_used,omitempty"`
	Hash                string `json:"hash,omitempty"`
	Level               uint64 `json:"level"`
	Nonce               uint64 `json:"nonce"`
	Timestamp           int64  `json:"timestamp,omitempty"`
	To                  string `json:"to,omitempty"`
	Type                string `json:"type,omitempty"`
}
