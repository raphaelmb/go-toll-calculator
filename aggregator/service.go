package main

import (
	"github.com/raphaelmb/go-toll-calculator/types"
	"github.com/sirupsen/logrus"
)

const basePrice = 3.15

type Aggregator interface {
	AggregateDistance(types.Distance) error
	CalculateInvoice(int) (*types.Invoice, error)
}

type Storer interface {
	Insert(types.Distance) error
	Get(int) (float64, error)
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
	logrus.WithFields(logrus.Fields{
		"obuid":    dist.OBUID,
		"distance": dist.Value,
		"unix":     dist.Unix,
	}).Info("aggregating distance")
	return i.Store.Insert(dist)
}

func (i *InvoiceAggregator) CalculateInvoice(id int) (*types.Invoice, error) {
	dist, err := i.Store.Get(id)
	if err != nil {
		return nil, err
	}
	invoice := &types.Invoice{
		OBUID:         id,
		TotalDistance: dist,
		TotalAmount:   basePrice * dist,
	}
	return invoice, nil
}
