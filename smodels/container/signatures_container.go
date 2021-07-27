package container

import (
	"oasisTracker/dmodels"
	"sync"
)

type BlockSignatureContainer struct {
	signs       []dmodels.BlockSignature
	mu          *sync.Mutex
	expectedCap uint64
}

func NewBlockSignatureContainer(cap ...uint64) *BlockSignatureContainer {
	var signs []dmodels.BlockSignature
	var expectedCap uint64
	if len(cap) == 1 {
		signs = make([]dmodels.BlockSignature, 0, cap[0])
		expectedCap = cap[0]
	}

	return &BlockSignatureContainer{
		mu:          &sync.Mutex{},
		signs:       signs,
		expectedCap: expectedCap,
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
	c.signs = make([]dmodels.BlockSignature, 0, c.expectedCap)
}
