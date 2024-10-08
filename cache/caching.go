package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	createdAt time.Time
	value     []byte
}

type Cache struct {
	entries map[string]CacheEntry
	mutex   *sync.Mutex
}

// NewCache creates a new cache with the given interval.
func NewCache(interval time.Duration) *Cache {
	newCache := &Cache{
		entries: make(map[string]CacheEntry),
		mutex:   &sync.Mutex{},
	}

	go newCache.reapLoop(interval)
	return newCache
}

// AddToCache adds a value to the cache with the given key.
func (cache *Cache) AddToCache(key string, val []byte) {
	cache.mutex.Lock()

	cache.entries[key] = CacheEntry{
		createdAt: time.Now(),
		value:     val,
	}

	defer cache.mutex.Unlock()
}

// GetFromCache gets a value from the cache with the given key.
func (cache *Cache) GetFromCache(key string) ([]byte, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if cache.entries[key].value != nil {
		return cache.entries[key].value, true
	}

	return nil, false
}

// reapLoop reaps the cache every interval.
func (cache *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		<-ticker.C
		cache.mutex.Lock()

		for key, item := range cache.entries {
			if time.Since(item.createdAt) > interval {
				delete(cache.entries, key)
			}
		}

		cache.mutex.Unlock()
	}
}
