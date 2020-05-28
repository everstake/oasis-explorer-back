package smodels

import (
	"oasisTracker/dmodels"
	"sync"
)

type TxsContainer struct {
	txs         []dmodels.Transaction
	registryTxs []dmodels.RegistryTransaction
	mu          *sync.Mutex
}

func NewTxsContainer() *TxsContainer {
	return &TxsContainer{
		mu:  &sync.Mutex{},
		txs: []dmodels.Transaction{},
	}
}

func (c *TxsContainer) Add(txs []dmodels.Transaction, registerTxs []dmodels.RegistryTransaction) {
	c.mu.Lock()
	c.txs = append(c.txs, txs...)
	c.registryTxs = append(c.registryTxs, registerTxs...)
	c.mu.Unlock()
}

func (c *TxsContainer) Txs() []dmodels.Transaction {
	return c.txs
}

func (c *TxsContainer) RegistryTxs() []dmodels.RegistryTransaction {
	return c.registryTxs
}

func (c *TxsContainer) IsEmpty() bool {
	return len(c.txs) == 0 && len(c.registryTxs) == 0
}

func (c *TxsContainer) Flush() {
	c.txs = []dmodels.Transaction{}
	c.registryTxs = []dmodels.RegistryTransaction{}
}
