package main

type Storable interface {
	GetKey() string
}

// Sample implementation of Storable interface for testing
type SampleStorable struct {
	key string
}

func (s SampleStorable) GetKey() string {
	return s.key
}
