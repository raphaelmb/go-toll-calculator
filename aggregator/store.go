package main

import "github.com/raphaelmb/go-toll-calculator/types"

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func (m *MemoryStore) Insert(dist types.Distance) error {
	m.data[dist.OBUID] = dist.Value
	return nil
}
