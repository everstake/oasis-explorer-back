package container

import (
	"oasisTracker/dmodels"
	"sync"
)

type RewardsContainer struct {
	rewards []dmodels.Reward
	mu      *sync.Mutex
}

func NewRewardsContainer() *RewardsContainer {
	return &RewardsContainer{
		mu: &sync.Mutex{},
	}
}

func (c *RewardsContainer) Add(rewards []dmodels.Reward) {
	if len(rewards) == 0 {
		return
	}
	c.mu.Lock()
	c.rewards = append(c.rewards, rewards...)
	c.mu.Unlock()
}

func (c *RewardsContainer) Rewards() []dmodels.Reward {
	return c.rewards
}

func (c *RewardsContainer) IsEmpty() bool {
	return len(c.rewards) == 0
}

func (c *RewardsContainer) Flush() {
	c.rewards = []dmodels.Reward{}
}
