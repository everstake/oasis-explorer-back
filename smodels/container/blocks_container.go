package container

import (
	"oasisTracker/dmodels"
	"sync"
)

type BlocksContainer struct {
	blocks      []dmodels.Block
	mu          *sync.Mutex
	expectedCap uint64
}

func NewBlocksContainer(cap ...uint64) *BlocksContainer {
	var blocks []dmodels.Block
	var expectedCap uint64

	if len(cap) == 1 {
		blocks = make([]dmodels.Block, 0, cap[0])
		expectedCap = cap[0]
	}

	return &BlocksContainer{
		mu:          &sync.Mutex{},
		blocks:      blocks,
		expectedCap: expectedCap,
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
	c.blocks = make([]dmodels.Block, 0, c.expectedCap)
}
