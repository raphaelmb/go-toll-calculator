package main

import (
	"fmt"

	"github.com/raphaelmb/go-toll-calculator/types"
)

type Aggregator interface {
	AggregateDistance(types.Distance) error
}

type Storer interface {
	Insert(types.Distance) error
}

type InvoiceAggregator struct {
	Store Storer
}

func NewInvoiceAggregator(store Storer) Aggregator {
	return &InvoiceAggregator{
		Store: store,
	}
}

func (i *InvoiceAggregator) AggregateDistance(dist types.Distance) error {
	fmt.Println("processing and inserting distance in storage", dist)
	return i.Store.Insert(dist)
}
