package pokecache

import (
	"testing"
	"time"
)

func TestCacheAddGet(t *testing.T) {
	cache := NewCache(200 * time.Millisecond)

	key := "test-key"
	val := []byte("hello")

	cache.Add(key, val)

	got, ok := cache.Get(key)
	if !ok {
		t.Fatalf("expected entry to exist in cache")
	}

	if string(got) != string(val) {
		t.Fatalf("want %s, got %s", val, got)
	}
}

func TestCacheReap(t *testing.T) {
	cache := NewCache(100 * time.Millisecond)

	cache.Add("old", []byte("data"))

	time.Sleep(250 * time.Millisecond)

	_, ok := cache.Get("old")
	if ok {
		t.Fatalf("expected entry to be reaped, but it still exists")
	}
}
