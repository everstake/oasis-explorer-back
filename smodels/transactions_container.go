package smodels

import (
	"oasisTracker/dmodels"
	"sync"
)

type TxsContainer struct {
	txs []dmodels.Transaction
	mu  *sync.Mutex
}

func NewTxsContainer() *TxsContainer {
	return &TxsContainer{
		mu:  &sync.Mutex{},
		txs: []dmodels.Transaction{},
	}
}

func (c *TxsContainer) Add(txs []dmodels.Transaction) {
	c.mu.Lock()
	c.txs = append(c.txs, txs...)
	c.mu.Unlock()
}

func (c *TxsContainer) Txs() []dmodels.Transaction {
	return c.txs
}

func (c *TxsContainer) IsEmpty() bool {
	return len(c.txs) == 0
}

func (c *TxsContainer) Flush() {
	c.txs = []dmodels.Transaction{}
}
