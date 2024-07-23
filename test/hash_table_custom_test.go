package main_test

import (
	"reflect"
	"testing"

	main "github.com/REKYAR/gredis"
)

// TestInsert tests the Insert function of the hash table
func TestInsert(t *testing.T) {
	ht := main.NewHashTable() // Assuming a constructor function NewHashTable
	ht.Insert("key1", "value1")
	expected := "value1"
	if got := ht.Get("key1"); got != expected {
		t.Errorf("Insert(\"key1\", \"value1\") = %v; want %v", got, expected)
	}
}

// TestRemove tests the Remove function of the hash table
func TestRemove(t *testing.T) {
	ht := main.NewHashTable() // Assuming a constructor function NewHashTable
	ht.Insert("key1", "value1")
	ht.Remove("key1")
	if got := ht.Get("key1"); got != nil {
		t.Errorf("Remove(\"key1\") = %v; want %v", got, nil)
	}
}

// TestInsertAndRemove tests the combination of Insert and Remove functions
func TestInsertAndRemove(t *testing.T) {
	tests := []struct {
		name   string
		key    string
		value  string
		remove string
		want   interface{}
	}{
		{"Insert and remove existing key", "key1", "value1", "key1", nil},
		{"Insert and remove non-existing key", "key2", "value2", "key3", nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ht := main.NewHashTable() // Assuming a constructor function NewHashTable
			ht.Insert(tt.key, tt.value)
			ht.Remove(tt.remove)
			if got := ht.Get(tt.key); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("After Insert(%q, %q) and Remove(%q), Get(%q) = %v; want %v", tt.key, tt.value, tt.remove, tt.key, got, tt.want)
			}
		})
	}
}
