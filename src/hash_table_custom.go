// Purpose: This file contains the implementation of a custom hash table data structure.
//we need segmentation not to lock the whole array for write operations
//we need to be able to route trhrough two levels of array, hash mod size1, then hash mod size2
//arrays not resizable because modulo on a fixed hash

package main

import (
	"errors"
	"sync"
)

const DEFAULT_HASHTAB_SIZE = 64
const DEFAULT_HASHTAB_WORKER_NUMBER = 512

type HMap struct {
	hasher *Hasher
	hts    []*HashTab // fixed amt of node arrays
	size   uint       // cnt of node arrays
}

type HashTab struct {
	hasher *Hasher
	rwmx   sync.RWMutex
	nodes  []*HashTabNode // node array
	size   uint
}

type HashTabNode struct {
	members []*Storable // array of Storables we don't expect to resize it frequently, will be appended to on demand
}

func NewHMap(HashMapWorkerNuber uint, HashMapSize uint) (*HMap, error) {
	if HashMapWorkerNuber == 0 {
		HashMapWorkerNuber = DEFAULT_HASHTAB_WORKER_NUMBER

	}
	if HashMapSize == 0 {
		HashMapSize = DEFAULT_HASHTAB_SIZE
	}

	h_ptr := &HMap{
		hts:    make([]*HashTab, HashMapWorkerNuber),
		size:   HashMapWorkerNuber,
		hasher: NewHasher(),
	}
	h_ptr.hasher.Init()
	for i := uint(0); i < HashMapWorkerNuber; i++ {
		ht, err := newHashtab(HashMapSize)
		if err != nil {
			return nil, err
		}
		if i == 0 {
			ht.hasher.Init()
		} else {
			ht.hasher = h_ptr.hts[0].hasher.CloneHasher()
		}

		h_ptr.hts[i] = ht
	}
	return h_ptr, nil
}

func newHashtab(size uint) (*HashTab, error) {
	// Corrected the condition to check if size is not a power of two
	if size == 0 || (size&(size-1)) != 0 {
		return nil, errors.New("hash table size must be a power of two")
	}
	h_ptr := &HashTab{
		nodes:  make([]*HashTabNode, size),
		size:   size,
		hasher: NewHasher(),
	}
	for i := uint(0); i < size; i++ {
		h_ptr.nodes[i] = &HashTabNode{
			members: make([]*Storable, 0),
		}

	}
	return h_ptr, nil
}

func (hmap *HMap) Insert(storable Storable) {
	idx := hmap.hasher.Hash(storable.GetKey()) % hmap.size
	tab := hmap.hts[idx]
	tab.insert(storable)
}

func (tab *HashTab) insert(storable Storable) {
	idx := tab.hasher.Hash(storable.GetKey()) % tab.size
	node := tab.nodes[idx]
	// Properly link the new node into the chain
	tab.rwmx.Lock()
	node.members = append(node.members, &storable)
	tab.rwmx.Unlock()
}

func (tab *HMap) Lookup(key string) (*HashTabNode, error) {
	idx := tab.hasher.Hash(key) % tab.size
	return tab.hts[idx].lookup(key)

}

func (tab *HashTab) lookup(key string) (*HashTabNode, error) {
	idx := tab.hasher.Hash(key) % tab.size
	node := tab.nodes[idx]
	tab.rwmx.RLock()
	defer tab.rwmx.RUnlock()
	for _, member := range node.members {
		if (*member).GetKey() == key {
			return node, nil
		}
	}
	return nil, errors.New("key not found")
}

func (tab *HMap) Delete(key string) error {
	idx := tab.hasher.Hash(key) % tab.size
	return tab.hts[idx].delete(key)
}

func (tab *HashTab) delete(key string) error {
	idx := tab.hasher.Hash(key) % tab.size
	node := tab.nodes[idx]
	tab.rwmx.Lock()
	defer tab.rwmx.Unlock()
	for i, member := range node.members {
		if (*member).GetKey() == key {
			node.members = append(node.members[:i], node.members[i+1:]...)
			return nil
		}
	}
	return errors.New("key not found")
}
