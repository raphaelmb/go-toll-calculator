package main

import (
	"time"

	"github.com/raphaelmb/go-toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type LogMiddleware struct {
	next Aggregator
}

func NewLogMiddleware(next Aggregator) *LogMiddleware {
	return &LogMiddleware{
		next: next,
	}
}

func (m *LogMiddleware) AggregateDistance(dist types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took": time.Since(start),
			"err":  err,
		}).Info("AggregateDistance")
	}(time.Now())
	err = m.next.AggregateDistance(dist)
	return
}

func (m *LogMiddleware) CalculateInvoice(id int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		var distance, amount float64
		if invoice != nil {
			distance = invoice.TotalDistance
			amount = invoice.TotalAmount
		}
		logrus.WithFields(logrus.Fields{
			"took":          time.Since(start),
			"err":           err,
			"obuID":         id,
			"totalDistance": distance,
			"totalAmount":   amount,
		}).Info("CalculateInvoice")
	}(time.Now())
	invoice, err = m.next.CalculateInvoice(id)
	return
}
