package service

import (
	"sync"
	"time"
)

// Cacher defines the interface for video metadata caching (in-memory or Redis)
type Cacher interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{}, ttl time.Duration)
	Delete(key string)
}

type cacheItem struct {
	value     interface{}
	expiresAt time.Time
}

// MemoryCache is an in-memory implementation of Cacher with TTL support
type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]*cacheItem
}

// NewMemoryCache creates a new MemoryCache and starts a background cleanup goroutine
func NewMemoryCache() *MemoryCache {
	c := &MemoryCache{
		items: make(map[string]*cacheItem),
	}
	go c.cleanup()
	return c
}

// Get retrieves a value from the cache; returns nil and false if not found or expired
func (c *MemoryCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	item, ok := c.items[key]
	c.mu.RUnlock()
	if !ok {
		return nil, false
	}
	if time.Now().After(item.expiresAt) {
		c.Delete(key)
		return nil, false
	}
	return item.value, true
}

// Set stores a value with the given TTL
func (c *MemoryCache) Set(key string, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	c.items[key] = &cacheItem{value: value, expiresAt: time.Now().Add(ttl)}
	c.mu.Unlock()
}

// Delete removes a key from the cache
func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	delete(c.items, key)
	c.mu.Unlock()
}

// cleanup periodically evicts expired entries (runs every 5 minutes)
func (c *MemoryCache) cleanup() {
	ticker := time.NewTicker(5 * time.Minute)
	for range ticker.C {
		now := time.Now()
		c.mu.Lock()
		for k, v := range c.items {
			if now.After(v.expiresAt) {
				delete(c.items, k)
			}
		}
		c.mu.Unlock()
	}
}
