package render

import (
	"oasisTracker/dmodels"
	"oasisTracker/smodels"
)

func Transactions(txs []dmodels.Transaction) []smodels.Transaction {
	transactions := make([]smodels.Transaction, len(txs))
	for i := range txs {
		transactions[i] = Transaction(txs[i])
	}
	return transactions
}

func Transaction(tx dmodels.Transaction) smodels.Transaction {

	return smodels.Transaction{
		Amount:              tx.Amount,
		EscrowAmount:        tx.EscrowAmount,
		ReclaimEscrowAmount: tx.EscrowReclaimAmount,
		Fee:                 tx.Fee,
		From:                tx.Sender,
		GasPrice:            tx.GasLimit,
		GasUsed:             tx.GasPrice,
		Hash:                tx.Hash,
		Level:               tx.BlockLevel,
		Nonce:               tx.Nonce,
		Timestamp:           tx.Time.Unix(),
		To:                  tx.Receiver,
		Type:                string(tx.Type),
		Status:              tx.Status,
		Error:               tx.Error,
	}
}
