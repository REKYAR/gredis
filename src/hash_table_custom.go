// Purpose: This file contains the implementation of a custom hash table data structure.
//we need segmentation not to lock the whole array for write operations
//we need to be able to route trhrough two levels of array, hash mod size1, then hash mod size2
//arrays not resizable because modulo on a fixed hash

package main

import (
	"errors"
	"sync"
)

//TODO: add default HMapNode size to preallocate memory

const DEFAULT_HMAP_SIZE = 32_768

// type HMap struct {
// 	hasher *Hasher
// 	hts    []*HashTab // fixed amt of node arrays
// 	size   uint       // cnt of node arrays
// }

type HMap struct {
	hasher *Hasher
	hts    []*HMapNode // fixed amt of node arrays
	size   uint        // cnt of node arrays
}

type HMapNode struct {
	rwmx    sync.RWMutex
	members []*Storable // array of Storables we don't expect to resize it frequently, will be appended to on demand
}

func NewHMap(HashMapSize uint) (*HMap, error) {

	if HashMapSize == 0 {
		HashMapSize = DEFAULT_HMAP_SIZE
	}

	h_ptr := &HMap{
		hts:    make([]*HMapNode, HashMapSize),
		size:   HashMapSize,
		hasher: NewHasher(),
	}
	h_ptr.hasher.Init()
	for i := uint(0); i < HashMapSize; i++ {
		ht, err := newHmapNode()
		if err != nil {
			return nil, err
		}
		h_ptr.hts[i] = ht
	}
	return h_ptr, nil
}

func newHmapNode() (*HMapNode, error) {
	h_ptr := &HMapNode{
		members: make([]*Storable, 0),
	}
	return h_ptr, nil
}

func (hmap *HMap) Insert(storable Storable) {
	idx := hmap.hasher.Hash(storable.GetKey()) % hmap.size
	tab := hmap.hts[idx]
	tab.rwmx.Lock()
	tab.members = append(tab.members, &storable)
	tab.rwmx.Unlock()
}

func (hmap *HMap) Lookup(key string) (*Storable, error) {
	idx := hmap.hasher.Hash(key) % hmap.size
	tab := hmap.hts[idx]
	tab.rwmx.RLock()
	defer tab.rwmx.RUnlock()
	for _, member := range tab.members {
		if (*member).GetKey() == key {
			return member, nil
		}
	}
	return nil, errors.New("key not found")
}

func (htab *HMap) Delete(key string) error {
	idx := htab.hasher.Hash(key) % htab.size
	tab := htab.hts[idx]
	tab.rwmx.Lock()
	defer tab.rwmx.Unlock()
	for i, member := range tab.members {
		if (*member).GetKey() == key {
			tab.members = append(tab.members[:i], tab.members[i+1:]...)
			return nil
		}

	}
	return errors.New("key not found")
}
