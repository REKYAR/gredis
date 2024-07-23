package main

import (
	"fmt"
	"log"
	"sync"
	"testing"
)

func TestHashMap(t *testing.T) {
	// Initialize hash map
	hMap, err := NewHMap(0, 0)
	if err != nil {
		log.Fatalf("Failed to create HMap: %v", err)
	}

	// Insert elements
	keys := []string{"key1", "key2", "key3"}
	for _, key := range keys {
		storable := &SampleStorable{key: key}
		hMap.Insert(storable)
	}

	// Lookup elements
	for _, key := range keys {
		node, err := hMap.Lookup(key)
		if err != nil {
			log.Printf("Lookup failed for key %s: %v", key, err)
		} else {
			fmt.Printf("Found key %s in node: %v\n", key, node)
		}
	}

	// Delete elements
	for _, key := range keys {
		err := hMap.Delete(key)
		if err != nil {
			log.Printf("Delete failed for key %s: %v", key, err)
		} else {
			fmt.Printf("Deleted key %s\n", key)
		}
	}

	// Lookup elements again to confirm deletion
	for _, key := range keys {
		_, err := hMap.Lookup(key)
		if err != nil {
			fmt.Printf("Confirmed deletion of key %s\n", key)
		} else {
			log.Printf("Key %s was not deleted properly", key)
		}
	}
}

func TestHashMapConcurrent(t *testing.T) {
	// Initialize hash map
	hMap, err := NewHMap(0, 0)
	if err != nil {
		log.Fatalf("Failed to create HMap: %v", err)
	}

	keys := []string{"key1", "key2", "key3"}
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
			node, err := hMap.Lookup(key)
			if err != nil {
				log.Printf("Lookup failed for key %s: %v", key, err)
			} else {
				fmt.Printf("Found key %s in node: %v\n", key, node)
			}
		}(key)
	}
	wg.Wait() // Wait for all lookups to complete

	// Concurrent Deletion
	wg.Add(len(keys))
	for _, key := range keys {
		go func(key string) {
			defer wg.Done()
			err := hMap.Delete(key)
			if err != nil {
				log.Printf("Delete failed for key %s: %v", key, err)
			} else {
				fmt.Printf("Deleted key %s\n", key)
			}
		}(key)
	}
	wg.Wait() // Wait for all deletions to complete

	// Final Lookup to confirm deletion
	wg.Add(len(keys))
	for _, key := range keys {
		go func(key string) {
			defer wg.Done()
			_, err := hMap.Lookup(key)
			if err != nil {
				fmt.Printf("Confirmed deletion of key %s\n", key)
			} else {
				log.Printf("Key %s was not deleted properly", key)
			}
		}(key)
	}
	wg.Wait() // Wait for all final lookups to complete
}
