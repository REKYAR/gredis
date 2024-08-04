package main

import (
	"sync"
	"testing"
)

func TestHashMap(t *testing.T) {
	// Initialize hash map
	hMap, err := NewHMap(1)
	if err != nil {
		t.Fatalf("Failed to create HMap: %v", err)
	}

	// Insert elements
	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		storable := &SampleStorable{key: key}
		hMap.Insert(storable)
	}

	// Lookup elements
	for _, key := range keys {
		if _, err := hMap.Lookup(key); err != nil {
			t.Errorf("Lookup failed for key %s: %v", key, err)
		}
	}

	// Test failed lookup
	if _, err := hMap.Lookup("nonexistent"); err == nil {
		t.Errorf("Expected lookup to fail for nonexistent key")
	}

	// Delete elements
	for _, key := range keys {
		if err := hMap.Delete(key); err != nil {
			t.Errorf("Delete failed for key %s: %v", key, err)
		}
	}

	// Lookup elements again to confirm deletion
	for _, key := range keys {
		if _, err := hMap.Lookup(key); err == nil {
			t.Errorf("Key %s was not deleted properly", key)
		}
	}
}

func TestHashMapConcurrent(t *testing.T) {
	// Initialize hash map
	hMap, err := NewHMap(1)
	if err != nil {
		t.Fatalf("Failed to create HMap: %v", err)
	}

	keys := []string{"key1", "key2", "key3", "key4", "key5", "key6", "key7", "key8", "key9", "key10"}
	var wg sync.WaitGroup

	// Concurrent Insertion
	wg.Add(len(keys))
	for _, key := range keys {
		go func(key string) {
			defer wg.Done()
			storable := &SampleStorable{key: key}
			hMap.Insert(storable)
		}(key)
	}
	wg.Wait() // Wait for all insertions to complete

	// Concurrent Lookup
	wg.Add(len(keys))
	for _, key := range keys {
		go func(key string) {
			defer wg.Done()
			if _, err := hMap.Lookup(key); err != nil {
				t.Errorf("Lookup failed for key %s: %v", key, err)
			}
		}(key)
	}
	wg.Wait() // Wait for all lookups to complete

	// Concurrent Deletion
	wg.Add(len(keys))
	for _, key := range keys {
		go func(key string) {
			defer wg.Done()
			if err := hMap.Delete(key); err != nil {
				t.Errorf("Delete failed for key %s: %v", key, err)
			}
		}(key)
	}
	wg.Wait() // Wait for all deletions to complete

	// Final Lookup to confirm deletion
	wg.Add(len(keys))
	for _, key := range keys {
		go func(key string) {
			defer wg.Done()
			if _, err := hMap.Lookup(key); err == nil {
				t.Errorf("Key %s was not deleted properly", key)
			}
		}(key)
	}
	wg.Wait() // Wait for all final lookups to complete
}
