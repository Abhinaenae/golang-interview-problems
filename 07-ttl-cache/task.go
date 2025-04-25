package main

import (
	"context"
	"sync"
	"time"
)

type entry struct {
	value string
	ttl   time.Time
}
type TtlCache struct {
	cache  map[string]entry
	mu     *sync.RWMutex
	cancel context.CancelFunc
}

func NewTtlCache() *TtlCache {
	ctx, cancel := context.WithCancel(context.Background())
	c := &TtlCache{
		cache:  make(map[string]entry),
		mu:     &sync.RWMutex{},
		cancel: cancel,
	}

	// Start the cleanup goroutine
	go c.cleanupExpiredKeys(ctx)

	return c
}

func (c *TtlCache) Set(key string, value string, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiration := time.Time{}
	if ttl > 0 {
		expiration = time.Now().Add(ttl)
	}
	newEntry := entry{
		value: value,
		ttl:   expiration,
	}
	c.cache[key] = newEntry
}

func (c *TtlCache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, exists := c.cache[key]
	if !exists {
		return "", false
	}

	//Checks if cache is expired
	if !entry.ttl.IsZero() && time.Now().After(entry.ttl) {
		return "", false
	}

	return entry.value, true
}

func (c *TtlCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
}

func (c *TtlCache) Stop() {
	c.cancel()
}

func (c *TtlCache) cleanupExpiredKeys(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now()
			for key, entry := range c.cache {
				if !entry.ttl.IsZero() && now.After(entry.ttl) {
					delete(c.cache, key)
				}
			}
			c.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}
