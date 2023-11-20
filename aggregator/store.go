package main

import (
	"fmt"
	"log"
	"os"

	"github.com/raphaelmb/go-toll-calculator/types"
)

type MemoryStore struct {
	data map[int]float64
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[int]float64),
	}
}

func makeStore() Storer {
	storeType := os.Getenv("AGGREGATOR_STORE_TYPE")
	switch storeType {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatalf("invalid store type give: %s", storeType)
		return nil
	}
}

func (m *MemoryStore) Insert(dist types.Distance) error {
	m.data[dist.OBUID] = dist.Value
	return nil
}

func (m *MemoryStore) Get(id int) (float64, error) {
	dist, ok := m.data[id]
	if !ok {
		return 0.0, fmt.Errorf("could not find distance for OBU ID %d", id)
	}
	return dist, nil
}
