// Package cache implements a thread-safe in-memory cache with expiration and periodic cleanup.
package cache

import (
	"hash/fnv"
	"sync"
	"time"
)

const numShards = 256

type Value struct {
	Text   string
	Expiry int64
}

func (v Value) isExpired() bool {
	return time.Now().Unix() > v.Expiry
}

type shard struct {
	items map[string]Value
	mu    sync.RWMutex
}

type Cache struct {
	shards      [numShards]*shard
	expiration  time.Duration
	stopCleanUp chan struct{}
	cleanUpDone sync.WaitGroup
}

func NewCache(expiration time.Duration) *Cache {
	cache := &Cache{
		expiration:  expiration,
		stopCleanUp: make(chan struct{}),
	}

	for i := 0; i < numShards; i++ {
		cache.shards[i] = &shard{
			items: make(map[string]Value),
		}
	}

	cache.startCleanUp()

	return cache
}

func (c *Cache) getShard(key string) *shard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return c.shards[h.Sum32()&(numShards-1)]
}

func (c *Cache) Set(key, value string) {
	shard := c.getShard(key)

	shard.mu.Lock()
	shard.items[key] = Value{
		Text:   value,
		Expiry: time.Now().Add(c.expiration).Unix(),
	}
	shard.mu.Unlock()
}

func (c *Cache) Get(key string) (string, bool) {
	shard := c.getShard(key)

	shard.mu.RLock()
	value, exists := shard.items[key]
	shard.mu.RUnlock()

	if !exists {
		return "", false
	}

	if value.isExpired() {
		shard.mu.Lock()
		if value, exists := shard.items[key]; exists && value.isExpired() {
			delete(shard.items, key)
		}
		shard.mu.Unlock()

		return "", false
	}

	return value.Text, true
}

func (c *Cache) startCleanUp() {
	c.cleanUpDone.Add(1)

	go func() {
		defer c.cleanUpDone.Done()

		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		shardIndex := 0
		shardsPerCleanUp := 16

		for {
			select {
			case <-ticker.C:
				for i := 0; i < shardsPerCleanUp; i++ {
					shard := c.shards[shardIndex]
					c.cleanUpShard(shard)
					shardIndex = (shardIndex + 1) % numShards
				}
			case <-c.stopCleanUp:
				return
			}
		}
	}()
}

func (c *Cache) cleanUpShard(shard *shard) {
	shard.mu.RLock()
	expiredKeys := make([]string, 0)
	for key, value := range shard.items {
		if value.isExpired() {
			expiredKeys = append(expiredKeys, key)
		}
	}
	shard.mu.RUnlock()

	if len(expiredKeys) > 0 {
		shard.mu.Lock()
		for _, key := range expiredKeys {
			if value, exists := shard.items[key]; exists &&
				value.isExpired() {
				delete(shard.items, key)
			}
		}
		shard.mu.Unlock()
	}
}

func (c *Cache) Close() {
	close(c.stopCleanUp)
	c.cleanUpDone.Wait()
}
