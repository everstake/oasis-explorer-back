package container

import (
	"oasisTracker/dmodels"
	"sync"
)

type TxsContainer struct {
	txs               []dmodels.Transaction
	nodeRegistryTxs   []dmodels.NodeRegistryTransaction
	entityRegistryTxs []dmodels.EntityRegistryTransaction
	mu                *sync.Mutex
}

func NewTxsContainer() *TxsContainer {
	return &TxsContainer{
		mu: &sync.Mutex{},
	}
}

func (c *TxsContainer) Add(txs []dmodels.Transaction, nodeRegistryTxs []dmodels.NodeRegistryTransaction, entityRegistryTxs []dmodels.EntityRegistryTransaction) {
	c.mu.Lock()
	c.txs = append(c.txs, txs...)
	c.nodeRegistryTxs = append(c.nodeRegistryTxs, nodeRegistryTxs...)
	c.entityRegistryTxs = append(c.entityRegistryTxs, entityRegistryTxs...)
	c.mu.Unlock()
}

func (c *TxsContainer) Txs() []dmodels.Transaction {
	return c.txs
}

func (c *TxsContainer) NodeRegistryTxs() []dmodels.NodeRegistryTransaction {
	return c.nodeRegistryTxs
}

func (c *TxsContainer) EntityRegistryTxs() []dmodels.EntityRegistryTransaction {
	return c.entityRegistryTxs
}

func (c *TxsContainer) IsEmpty() bool {
	return len(c.txs) == 0 && len(c.nodeRegistryTxs) == 0 && len(c.entityRegistryTxs) == 0
}

func (c *TxsContainer) Flush() {
	c.txs = []dmodels.Transaction{}
	c.nodeRegistryTxs = []dmodels.NodeRegistryTransaction{}
	c.entityRegistryTxs = []dmodels.EntityRegistryTransaction{}
}
