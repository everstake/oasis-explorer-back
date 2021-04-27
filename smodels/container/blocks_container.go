package container

import (
	"oasisTracker/dmodels"
	"sync"
)

type BlocksContainer struct {
	blocks []dmodels.Block
	mu     *sync.Mutex
}

func NewBlocksContainer() *BlocksContainer {
	return &BlocksContainer{
		mu: &sync.Mutex{},
	}
}

func (c *BlocksContainer) Add(blocks []dmodels.Block) {
	if len(blocks) == 0 {
		return
	}
	c.mu.Lock()
	c.blocks = append(c.blocks, blocks...)
	c.mu.Unlock()
}

func (c *BlocksContainer) Blocks() []dmodels.Block {
	return c.blocks
}

func (c *BlocksContainer) IsEmpty() bool {
	return len(c.blocks) == 0
}

func (c *BlocksContainer) Flush() {
	c.blocks = []dmodels.Block{}
}
