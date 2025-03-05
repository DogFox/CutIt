package cache

import (
	"testing"
)

func TestLRUCache(t *testing.T) {
	cache := NewCache(2)

	cache.Put("key1", "value1")
	if val, ok := cache.Get("key1"); !ok || val != "value1" {
		t.Errorf("expected value1, got %v", val)
	}

	cache.Put("key2", "value2")
	if val, ok := cache.Get("key2"); !ok || val != "value2" {
		t.Errorf("expected value2, got %v", val)
	}

	cache.Put("key3", "value3")
	if _, ok := cache.Get("key1"); ok {
		t.Errorf("expected key1 to be evicted")
	}
}
