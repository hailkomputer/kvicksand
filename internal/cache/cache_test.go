package cache_test

import (
	"os"
	"strings"
	"testing"

	"github.com/hailkomputer/kvicksand/internal/cache"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestCache(t *testing.T) {

	cache := cache.NewCache()

	_, ok := cache.Get("k1")
	if ok {
		t.Errorf("value for key should not exist")
	}

	cache.Set("k1", "v1")

	v1, ok := cache.Get("k1")
	if !ok {
		t.Errorf("value for key should exist")
	}
	if !strings.EqualFold(v1, "v1") {
		t.Errorf("value for key should match")
	}
}
