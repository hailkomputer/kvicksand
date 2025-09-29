package cache_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/hailkomputer/kvicksand/internal/cache"
)

func TestBasicOperations(t *testing.T) {
	cache := cache.NewCache(30 * time.Minute)
	defer cache.Close()

	tests := []struct {
		name      string
		operation string
		key       string
		value     string
		want      string
		wantOk    bool
	}{
		{
			name:      "set_and_get_simple_value",
			operation: "set",
			key:       "key1",
			value:     "value1",
			want:      "value1",
			wantOk:    true,
		},
		{
			name:      "get_existing_key",
			operation: "get",
			key:       "key1",
			want:      "value1",
			wantOk:    true,
		},
		{
			name:      "get_non_existent_key",
			operation: "get",
			key:       "nonexistent",
			want:      "",
			wantOk:    false,
		},
		{
			name:      "overwrite_existing_key",
			operation: "set",
			key:       "key1",
			value:     "newvalue1",
			want:      "newvalue1",
			wantOk:    true,
		},
		{
			name:      "set_empty_value",
			operation: "set",
			key:       "emptykey",
			value:     "",
			want:      "",
			wantOk:    true,
		},
		{
			name:      "get_empty_value",
			operation: "get",
			key:       "emptykey",
			want:      "",
			wantOk:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.operation {
			case "set":
				cache.Set(tt.key, tt.value)
				// Verify the set worked
				got, ok := cache.Get(tt.key)
				if ok != tt.wantOk {
					t.Errorf("Get() after Set() ok = %v, wantOk %v", ok, tt.wantOk)
				}
				if got != tt.want {
					t.Errorf("Get() after Set() = %v, want %v", got, tt.want)
				}
			case "get":
				got, ok := cache.Get(tt.key)
				if ok != tt.wantOk {
					t.Errorf("Get() ok = %v, wantOk %v", ok, tt.wantOk)
				}
				if got != tt.want {
					t.Errorf("Get() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	tests := []struct {
		name            string
		numWriters      int
		numReaders      int
		numKeys         int
		writesPerWriter int
		readsPerReader  int
		expiration      time.Duration
	}{
		{
			name:            "balanced_load",
			numWriters:      10,
			numReaders:      10,
			numKeys:         100,
			writesPerWriter: 100,
			readsPerReader:  100,
			expiration:      1 * time.Second,
		},
		{
			name:            "read_heavy_80_20",
			numWriters:      2,
			numReaders:      8,
			numKeys:         100,
			writesPerWriter: 100,
			readsPerReader:  400,
			expiration:      1 * time.Second,
		},
		{
			name:            "write_heavy",
			numWriters:      8,
			numReaders:      2,
			numKeys:         100,
			writesPerWriter: 100,
			readsPerReader:  100,
			expiration:      1 * time.Second,
		},
		{
			name:            "high_contention_few_keys",
			numWriters:      50,
			numReaders:      50,
			numKeys:         5,
			writesPerWriter: 100,
			readsPerReader:  100,
			expiration:      1 * time.Second,
		},
		{
			name:            "low_contention_many_keys",
			numWriters:      10,
			numReaders:      10,
			numKeys:         1000,
			writesPerWriter: 100,
			readsPerReader:  100,
			expiration:      1 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := cache.NewCache(tt.expiration)
			defer cache.Close()

			var wg sync.WaitGroup

			// Launch writers
			for w := 0; w < tt.numWriters; w++ {
				wg.Add(1)
				go func(writerID int) {
					defer wg.Done()
					for i := 0; i < tt.writesPerWriter; i++ {
						key := fmt.Sprintf("key%d", i%tt.numKeys)
						value := fmt.Sprintf("writer%d_value%d", writerID, i)
						cache.Set(key, value)
					}
				}(w)
			}

			// Launch readers
			for r := 0; r < tt.numReaders; r++ {
				wg.Add(1)
				go func(readerID int) {
					defer wg.Done()
					for i := 0; i < tt.readsPerReader; i++ {
						key := fmt.Sprintf("key%d", i%tt.numKeys)
						cache.Get(key)
					}
				}(r)
			}

			// Wait for all operations to complete
			wg.Wait()

			// Verify cache is still functional
			cache.Set("test_key", "test_value")
			val, ok := cache.Get("test_key")
			if !ok || val != "test_value" {
				t.Error("Cache not functional after concurrent access")
			}
		})
	}
}

// BenchmarkCacheOperations benchmarks different cache operations
func BenchmarkCacheOperations(b *testing.B) {
	benchmarks := []struct {
		name          string
		cacheSize     int
		numShards     int
		operationType string
		keyPattern    string
	}{
		{
			name:          "Get_Small_Cache_Sequential",
			cacheSize:     100,
			numShards:     1,
			operationType: "get",
			keyPattern:    "sequential",
		},
		{
			name:          "Get_Small_Cache_Random",
			cacheSize:     100,
			numShards:     1,
			operationType: "get",
			keyPattern:    "random",
		},
		{
			name:          "Get_Large_Cache_Sequential",
			cacheSize:     10000,
			numShards:     1,
			operationType: "get",
			keyPattern:    "sequential",
		},
		{
			name:          "Get_Sharded_256_Sequential",
			cacheSize:     10000,
			numShards:     256,
			operationType: "get",
			keyPattern:    "sequential",
		},
		{
			name:          "Set_Small_Cache",
			cacheSize:     100,
			numShards:     1,
			operationType: "set",
			keyPattern:    "sequential",
		},
		{
			name:          "Set_Sharded_256",
			cacheSize:     10000,
			numShards:     256,
			operationType: "set",
			keyPattern:    "sequential",
		},
		{
			name:          "Mixed_80_20_Read_Heavy",
			cacheSize:     10000,
			numShards:     256,
			operationType: "mixed",
			keyPattern:    "random",
		},
	}

	for _, bm := range benchmarks {
		b.Run(bm.name, func(b *testing.B) {
			// Note: This assumes your cache implementation can accept numShards parameter
			// Adjust based on your actual implementation
			cache := cache.NewCache(30 * time.Minute)
			defer cache.Close()

			// Pre-populate cache
			for i := 0; i < bm.cacheSize; i++ {
				cache.Set(fmt.Sprintf("key%d", i), fmt.Sprintf("value%d", i))
			}

			b.ResetTimer()

			switch bm.operationType {
			case "get":
				b.RunParallel(func(pb *testing.PB) {
					i := 0
					for pb.Next() {
						var key string
						if bm.keyPattern == "sequential" {
							key = fmt.Sprintf("key%d", i%bm.cacheSize)
						} else {
							key = fmt.Sprintf("key%d", fastrand()%uint32(bm.cacheSize))
						}
						cache.Get(key)
						i++
					}
				})

			case "set":
				b.RunParallel(func(pb *testing.PB) {
					i := 0
					for pb.Next() {
						key := fmt.Sprintf("newkey%d", i)
						cache.Set(key, "value")
						i++
					}
				})

			case "mixed":
				b.RunParallel(func(pb *testing.PB) {
					i := 0
					for pb.Next() {
						// 80% reads, 20% writes
						if fastrand()%5 == 0 {
							key := fmt.Sprintf("writekey%d", i)
							cache.Set(key, "value")
						} else {
							key := fmt.Sprintf("key%d", fastrand()%uint32(bm.cacheSize))
							cache.Get(key)
						}
						i++
					}
				})
			}
		})
	}
}

// fastrand is a fast pseudo-random number generator for benchmarks
// This avoids the overhead of math/rand in benchmarks
func fastrand() uint32 {
	var x uint32 = 0x123456
	x ^= x << 13
	x ^= x >> 17
	x ^= x << 5
	return x
}
