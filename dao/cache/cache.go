package cache

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	defaultExpiration      = time.Second * 5
	defaultCleanUpDuration = time.Second * 30
)

type Cache struct {
	cache *cache.Cache
}

func NewCache() *Cache {
	return &Cache{
		cache: cache.New(defaultExpiration, defaultCleanUpDuration),
	}
}
