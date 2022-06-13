package cache

import (
	"sync"
	"time"
)

// Value is a struct that holds the value and expiry time for a key
type Value struct {
	Text   string
	Expiry int64
}

// isExpired checks whether the Value is expired or not
func (v Value) isExpired() bool {
	return time.Now().Unix() > v.Expiry
}

// Cache is a struct that holds the key value pairs and expiration duration
type Cache struct {
	keyValuePair map[string]Value
	expiration   time.Duration
	mutex        sync.RWMutex
}

// NewCache initializes a new cache with a given expiration duration for items in the cache
// and also starts to run a cleanup process
func NewCache(expiration time.Duration) *Cache {
	keyValuePair := make(map[string]Value)
	cache := &Cache{keyValuePair: keyValuePair, expiration: expiration}

	go func() {
		cache.cleanup()
	}()

	return cache
}

// Set writes the value for the given key using write mutex locks.
func (c *Cache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.keyValuePair[key] = Value{
		Text:   value,
		Expiry: time.Now().Add(c.expiration).Unix(),
	}
}

// Get tries to fetch the value for the given key using read mutexes.
// It will return second argument as false, in cases where value is not found or expired
func (c *Cache) Get(key string) (string, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, ok := c.keyValuePair[key]
	if !ok {
		return "", false
	}
	if value.isExpired() {
		return "", false
	}
	return value.Text, true
}

// cleanup runs every half a second until timeout context forces it on shutdown.
// It will delete all expired key value pairs while locking with a write mutex
func (c *Cache) cleanup() {
	ticker := time.NewTicker(30 * time.Millisecond)
	defer ticker.Stop()

	for range ticker.C {
		c.mutex.Lock()
		for key, value := range c.keyValuePair {
			if value.isExpired() {
				delete(c.keyValuePair, key)
			}
		}
		c.mutex.Unlock()
	}
}
