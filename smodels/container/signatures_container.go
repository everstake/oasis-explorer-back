package container

import (
	"oasisTracker/dmodels"
	"sync"
)

type BlockSignatureContainer struct {
	signs []dmodels.BlockSignature
	mu    *sync.Mutex
}

func NewBlockSignatureContainer() *BlockSignatureContainer {
	return &BlockSignatureContainer{
		mu: &sync.Mutex{},
	}
}

func (c *BlockSignatureContainer) Add(signs []dmodels.BlockSignature) {
	c.mu.Lock()
	c.signs = append(c.signs, signs...)
	c.mu.Unlock()
}

func (c *BlockSignatureContainer) Signatures() []dmodels.BlockSignature {
	return c.signs
}

func (c BlockSignatureContainer) IsEmpty() bool {
	return len(c.signs) == 0
}

func (c *BlockSignatureContainer) Flush() {
	c.signs = []dmodels.BlockSignature{}
}
