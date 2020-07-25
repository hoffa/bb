// Package cache provides a very basic thread-safe cache with random replacement.
package cache

import "sync"

type Cache struct {
	data map[string][]byte
	size int
	sync.RWMutex
}

func New(size int) *Cache {
	return &Cache{
		data: make(map[string][]byte),
		size: size,
	}
}

func (c *Cache) evict() {
	if len(c.data) >= c.size {
		for k := range c.data {
			delete(c.data, k)
			break
		}
	}
}

func (c *Cache) Put(k string, b []byte) {
	c.Lock()
	defer c.Unlock()
	c.evict()
	c.data[k] = b
}

func (c *Cache) Get(k string) []byte {
	c.RLock()
	defer c.RUnlock()
	return c.data[k]
}
