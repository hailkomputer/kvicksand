package cache

import (
	"sync"
	"time"
)

type Value struct {
	Text   string
	Expiry int64
}

// isExpired checks whether the Value is expired or not
func (v Value) isExpired() bool {
	return time.Now().Unix() > v.Expiry
}

type Cache struct {
	keyValuePair map[string]Value
	mutex        sync.RWMutex
	wg           sync.WaitGroup
	sc           chan bool
}

// NewCache initializes a new cache and also starts to run a cleanup process
func NewCache() *Cache {
	keyValuePair := make(map[string]Value)
	cache := &Cache{keyValuePair: keyValuePair, sc: make(chan bool)}

	cache.wg.Add(1)
	go func() {
		defer cache.wg.Done()
		cache.cleanup()
	}()

	return cache
}

// Set writes the value for the given key using write mutex locks.
// Expiration duration is hard coded as 30 minutes
func (c *Cache) Set(key, value string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.keyValuePair[key] = Value{
		Text:   value,
		Expiry: time.Now().Add(time.Minute * 30).Unix(),
	}
}

// Get tries to fetch the value for the given key using read mutexes.
// It will always return second argument as false, in cases where value is
// not found or expired
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

// cleanup runs every 5 seconds until timeout context forces it on shutdown.
// It will delete all expired key value pairs while using a write mutex
func (c *Cache) cleanup() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.sc:
			return
		case <-ticker.C:
			c.mutex.Lock()
			for key, value := range c.keyValuePair {
				if value.isExpired() {
					delete(c.keyValuePair, key)
				}
			}
			c.mutex.Unlock()
		}
	}
}
