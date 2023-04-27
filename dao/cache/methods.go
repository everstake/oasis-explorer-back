package cache

import "time"

func (c *Cache) Save(key string, item interface{}, expiration time.Duration) (err error) {
	c.cache.Set(key, item, expiration)
	return nil
}

func (c *Cache) Get(key string) (interface{}, bool, error) {
	item, ok := c.cache.Get(key)
	if !ok {
		return item, false, nil
	}

	return item, true, nil
}

func (c *Cache) Remove(key string) error {
	c.cache.Delete(key)
	return nil
}
