package pokecache

import (
	"sync"
	"time"
	"fmt"
)

type Cache struct {
	cacheEntries	map[string]cacheEntry
	mu				sync.Mutex
}

type cacheEntry struct {
	createdAt		time.Time
	val				[]byte
}

func NewCache(interval time.Duration) *Cache {
	cache := &Cache{
		cacheEntries: make(map[string]cacheEntry),
	}
	go cache.reapLoop(interval)
	return cache
}

func (c *Cache) Add(key string, val []byte) {
	new_entry := cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
	c.mu.Lock()
	c.cacheEntries[key] = new_entry
	c.mu.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	entry, ok := c.cacheEntries[key]
	c.mu.Unlock()

	if !ok {
		//key not found
		return nil, false
	}

	// otherwise return the cacheEntries at the given key
	fmt.Println("Accessing cached data...")
	return entry.val, true

}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)

	for {
		<-ticker.C 

		c.mu.Lock()

		now := time.Now()
		for key, entry := range c.cacheEntries {
			if now.Sub(entry.createdAt) > interval {
				delete(c.cacheEntries, key)
			}
		}

		c.mu.Unlock()
	}
}