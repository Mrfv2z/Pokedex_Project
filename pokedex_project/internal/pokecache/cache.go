package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	ourMap   map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	ourCache := &Cache{
		ourMap:   make(map[string]cacheEntry),
		interval: interval,
	}
	go ourCache.reapLoop()
	return ourCache
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.ourMap[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.ourMap[key]
	if !ok {
		return []byte{}, false
	} else {
		return entry.val, true
	}
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for range ticker.C {
		c.mu.Lock()
		for i := range c.ourMap {
			if time.Since(c.ourMap[i].createdAt) > c.interval {
				delete(c.ourMap, i)
			}
		}
		c.mu.Unlock()
	}
}
