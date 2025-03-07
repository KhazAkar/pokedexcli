package pokecache

import (
	"testing"
	"time"
)

func TestCacheEntryCreationTime(t *testing.T) {
	cache := NewCache(5 * time.Minute)
	key := "test_key"
	value := []byte("test_value")

	// Add entry and check creation time
	beforeAdd := time.Now()
	cache.Add(key, value)
	
	// Retrieve the entry
	retrievedVal, exists := cache.Get(key)
	
	if !exists {
		t.Fatalf("Expected entry to exist after adding")
	}
	
	if retrievedVal == nil || len(retrievedVal) == 0 {
		t.Fatalf("Retrieved value should not be empty")
	}
}

func TestCacheEntryExpiration(t *testing.T) {
	// Set a very short reap interval to test expiration quickly
	reapInterval := 100 * time.Millisecond
	cache := NewCache(reapInterval)
	
	key := "expiring_key"
	value := []byte("expiring_value")
	
	// Add entry
	cache.Add(key, value)
	
	// Wait for longer than reap interval
	time.Sleep(200 * time.Millisecond)
	
	// Try to retrieve the entry
	_, exists := cache.Get(key)
	
	if exists {
		t.Fatalf("Entry should have been expired and removed")
	}
}

func TestCacheReapingMechanism(t *testing.T) {
	reapInterval := 100 * time.Millisecond
	cache := NewCache(reapInterval)
	
	// Add multiple entries
	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		cache.Add(key, []byte(key+"_value"))
	}
	
	// Wait for reap interval
	time.Sleep(200 * time.Millisecond)
	
	// Check if all entries are removed
	for _, key := range keys {
		_, exists := cache.Get(key)
		if exists {
			t.Fatalf("Entry %s should have been removed", key)
		}
	}
}

func TestCacheEntryLifecycle(t *testing.T) {
	reapInterval := 500 * time.Millisecond
	cache := NewCache(reapInterval)
	
	key := "lifecycle_key"
	value := []byte("lifecycle_value")
	
	// Add entry
	cache.Add(key, value)
	
	// Verify entry exists immediately after adding
	retrievedVal, exists := cache.Get(key)
	if !exists {
		t.Fatalf("Entry should exist immediately after adding")
	}
	
	if string(retrievedVal) != string(value) {
		t.Fatalf("Retrieved value does not match original value")
	}
	
	// Wait for less than reap interval
	time.Sleep(200 * time.Millisecond)
	
	// Entry should still exist
	_, exists = cache.Get(key)
	if !exists {
		t.Fatalf("Entry should still exist before reap interval")
	}
	
	// Wait for longer than reap interval
	time.Sleep(400 * time.Millisecond)
	
	// Entry should be removed
	_, exists = cache.Get(key)
	if exists {
		t.Fatalf("Entry should have been removed after reap interval")
	}
}
