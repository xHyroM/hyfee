package utils

import (
	"time"
)

type CacheItem struct {
	value      interface{}
	expiration time.Time
}

type Cache struct {
	items map[string]CacheItem
}

func NewCache() *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
	}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	expiration := time.Now().Add(ttl)
	c.items[key] = CacheItem{
		value:      value,
		expiration: expiration,
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	item, found := c.items[key]
	if !found {
		return nil, false
	}

	if time.Now().After(item.expiration) {
		delete(c.items, key)
		return nil, false
	}

	return item.value, true
}