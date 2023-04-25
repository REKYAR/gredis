package main

import "errors"

const resizing_work = 128

type HMap struct {
	ht1  HashTab
	ht2  HashTab
	size uint
}

type HashTab struct {
	nodes []HashTabNode
	mask  uint
	size  uint
}

type HashTabNode struct {
	hcode       uint
	next        *HashTabNode
	initialised bool
}

func NewHashtab(size uint) (*HashTab, error) {
	if (size > 0) && (((size - 1) & size) == 0) {
		return nil, errors.New("hash tabe size must be a power of two")
	}
	h_ptr := new(HashTab)
	h_ptr.nodes = make([]HashTabNode, size)
	h_ptr.mask = size - 1
	h_ptr.size = 0
	return h_ptr, nil
}

func (tab HashTab) Insert(node HashTabNode) {
	pos := node.hcode & tab.mask
	next := tab.nodes[pos]
	node.next = &next
	node.initialised = true
	tab.nodes[pos] = node
	tab.size++
}

// cmp true if nodes identical
func (tab HashTab) Lookup(node HashTabNode, cmp func(a HashTabNode, b HashTabNode) bool) **HashTabNode {
	pos := node.hcode & tab.mask
	from := &(tab.nodes[pos])
	//check if first is null
	if !tab.nodes[pos].initialised {
		return nil
	}
	//check if from->next (the one we actually want) is null
	for from.next != nil {
		if cmp(node, *from.next) {
			//return address of pointer to next
			return &(from.next)
		}
		from = from.next
	}
	return nil
}

func (tab HashTab) Detach(node **HashTabNode) *HashTabNode {
	detached := *node
	*node = (*node).next
	tab.size--
	return detached
}

func (hmap HMap) Lookup(node HashTabNode, cmp func(a HashTabNode, b HashTabNode) bool) *HashTabNode {
	//resize

	from := hmap.ht1.Lookup(node, cmp)
	if from == nil {
		from = hmap.ht2.Lookup(node, cmp)
	}
	if from == nil {
		return nil
	}
	return *from
}

func (hmap HMap) HelpResizing() {
	if hmap.ht2 == nil {
		return
	}
}