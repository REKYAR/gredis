// package main

// import "errors"

// const resizing_work = 128

// type HMap struct {
// 	ht1  HashTab
// 	ht2  HashTab
// 	size uint
// }

// type HashTab struct {
// 	nodes []HashTabNode
// 	mask  uint
// 	size  uint
// }

// type HashTabNode struct {
// 	hcode       uint
// 	next        *HashTabNode
// 	initialised bool
// }

// func NewHashtab(size uint) (*HashTab, error) {
// 	if (size > 0) && (((size - 1) & size) != 0) {
// 		return nil, errors.New("hash tabe size must be a power of two")
// 	}
// 	h_ptr := new(HashTab)
// 	h_ptr.nodes = make([]HashTabNode, size)
// 	h_ptr.mask = size - 1
// 	h_ptr.size = 0
// 	return h_ptr, nil
// }

// func (tab HashTab) Insert(node HashTabNode) {
// 	pos := node.hcode & tab.mask
// 	next := tab.nodes[pos]
// 	node.next = &next
// 	node.initialised = true
// 	tab.nodes[pos] = node
// 	tab.size++
// }

// // cmp true if nodes identical
// func (tab HashTab) Lookup(node HashTabNode, cmp func(a HashTabNode, b HashTabNode) bool) **HashTabNode {
// 	pos := node.hcode & tab.mask
// 	from := &(tab.nodes[pos])
// 	//check if first is null
// 	if !tab.nodes[pos].initialised {
// 		return nil
// 	}
// 	//check if from->next (the one we actually want) is null
// 	for from.next != nil {
// 		if cmp(node, *from.next) {
// 			//return address of pointer to next
// 			return &(from.next)
// 		}
// 		from = from.next
// 	}
// 	return nil
// }

// func (tab HashTab) Detach(node **HashTabNode) *HashTabNode {
// 	detached := *node
// 	*node = (*node).next
// 	tab.size--
// 	return detached
// }

// func (hmap HMap) Lookup(node HashTabNode, cmp func(a HashTabNode, b HashTabNode) bool) *HashTabNode {
// 	//resize

// 	from := hmap.ht1.Lookup(node, cmp)
// 	if from == nil {
// 		from = hmap.ht2.Lookup(node, cmp)
// 	}
// 	if from == nil {
// 		return nil
// 	}
// 	return *from
// }

// func (hmap HMap) HelpResizing() {
// 	if hmap.ht2 == nil {
// 		return
// 	}
// }

package main

import "errors"

const DEFAULT_HASHTAB_SIZE = 128
const DEFAULT_HASHTAB_WORKER_NUMBER = 128

type HMap struct {
	hts  []*HashTab // fixed amt of node arrays
	size uint       // cnt of node arrays
}

type HashTab struct {
	nodes []*HashTabNode // node array
	mask  uint
	size  uint
}

type HashTabNode struct {
	members     []Storable // array of Storables we don't expect to resize it frequently, will be appended to on demand
	hcode       uint
	next        *HashTabNode
	initialised bool
}

func NewHMap(HashMapWorkerNuber uint, HashMapSize uint) (*HMap, error) {
	if HashMapWorkerNuber == 0 {
		HashMapWorkerNuber = DEFAULT_HASHTAB_WORKER_NUMBER

	}
	if HashMapSize == 0 {
		HashMapSize = DEFAULT_HASHTAB_SIZE
	}

	h_ptr := &HMap{
		hts:  make([]*HashTab, HashMapWorkerNuber),
		size: HashMapWorkerNuber,
	}
	for i := uint(0); i < HashMapWorkerNuber; i++ {
		ht, err := newHashtab(HashMapSize)
		if err != nil {
			return nil, err
		}
		for j := uint(0); j < HashMapSize; j++ {
			ht.nodes[j] = &HashTabNode{
				members: make([]Storable, 0),
			}
			h_ptr.hts[i] = ht
		}
	}
	return h_ptr, nil
}

func newHashtab(size uint) (*HashTab, error) {
	// Corrected the condition to check if size is not a power of two
	if size == 0 || (size&(size-1)) != 0 {
		return nil, errors.New("hash table size must be a power of two")
	}
	h_ptr := &HashTab{
		nodes: make([]*HashTabNode, size),
		mask:  size - 1,
		size:  0,
	}
	return h_ptr, nil
}

func (tab *HashTab) insert(node *HashTabNode) {
	pos := node.hcode & tab.mask
	// Properly link the new node into the chain
	node.next = tab.nodes[pos]
	tab.nodes[pos] = node
	tab.size++
}

func (tab *HashTab) lookup(node HashTabNode, cmp func(a, b HashTabNode) bool) *HashTabNode {
	pos := node.hcode & tab.mask
	for current := tab.nodes[pos]; current != nil; current = current.next {
		if cmp(node, *current) {
			return current
		}
	}
	return nil
}

func (tab *HashTab) Detach(node **HashTabNode) {
	if *node != nil {
		detached := *node
		*node = (*node).next
		tab.size--
		*node = detached.next // Correctly detach the node
	}
}

func (hmap *HMap) HelpResizing() {
	// Placeholder for resizing logic
	// This should include allocating a new HashTab with a larger size,
	// and rehashing all existing entries into it.
	newsize := hmap.size * 2
	if hmap.ht2 == nil {
		ht2, err := newHashtab(newsize)
		if err != nil {
			panic(err)
		}
		ht2.nodes = hmap.ht1.nodes
		hmap.ht1.nodes = nil

		hmap.ht2 = ht2
	}
}

func (hmap *HMap) Lookup(node HashTabNode, cmp func(a, b HashTabNode) bool) *HashTabNode {
	// Adjusted to directly use pointers for ht1 and ht2
	var from *HashTabNode
	if hmap.ht1 != nil {
		from = hmap.ht1.Lookup(node, cmp)
	}
	if from == nil && hmap.ht2 != nil {
		from = hmap.ht2.Lookup(node, cmp)
	}
	return from
}

func (hmap *HMap) selectActive() *HashTab {
	if hmap.ht1 == nil {
		return hmap.ht2
	}
	return hmap.ht1
}
