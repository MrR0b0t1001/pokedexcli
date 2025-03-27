package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	cacheMap   map[string]cacheEntry
	expiration time.Duration
	mu         sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(expiration time.Duration) *Cache {
	cache := &Cache{
		cacheMap:   map[string]cacheEntry{},
		expiration: expiration,
	}

	go cache.reapLoop()

	return cache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheMap[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	cEntry, ok := c.cacheMap[key]
	if !ok {
		return []byte{}, false
	}
	return cEntry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.expiration)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			for key, value := range c.cacheMap {
				if time.Since(value.createdAt) > c.expiration {
					delete(c.cacheMap, key)
				}
			}
			c.mu.Unlock()
		}
	}
}
