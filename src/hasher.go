package main

import (
	"hash/maphash"
	"sync"
)

type Hasher struct {
	mx     sync.Mutex
	hasher maphash.Hash
}

func (h *Hasher) Hash(key string) uint {
	h.mx.Lock()
	defer h.mx.Unlock()
	h.hasher.Reset()
	h.hasher.WriteString(key)
	return uint(h.hasher.Sum64())
}

func (h *Hasher) CloneHasher() *Hasher {
	hs := maphash.Hash{}
	hs.SetSeed(h.hasher.Seed())
	return &Hasher{
		hasher: hs,
	}
}

func (h *Hasher) Init() {
	seed := maphash.MakeSeed()
	h.hasher.SetSeed(seed)
}

func NewHasher() *Hasher {
	return &Hasher{}
}
