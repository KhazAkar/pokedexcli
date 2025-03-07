package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mtx     *sync.Mutex
}

func NewCache(reapInterval time.Duration) *Cache {
	c := &Cache{
		entries: make(map[string]cacheEntry),
		mtx:     &sync.Mutex{},
	}
	go c.reapLoop(reapInterval)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mtx.Lock()
	defer c.mtx.Unlock()

	entry, exists := c.entries[key]
	return entry.val, exists
}

func (c *Cache) reapLoop(reapInterval time.Duration) {
	ticker := time.NewTicker(reapInterval)
	for range ticker.C {
		c.mtx.Lock()
		for key, entry := range c.entries {
			if time.Since(entry.createdAt) > reapInterval {
				delete(c.entries, key)
			}
		}
		c.mtx.Unlock()
	}
}
