package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/raphaelmb/go-toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type MetricsMiddleware struct {
	reqCounterAgg  prometheus.Counter
	reqCounterCalc prometheus.Counter
	errCounterAgg  prometheus.Counter
	errCounterCalc prometheus.Counter
	reqLatencyAgg  prometheus.Histogram
	reqLatencyCalc prometheus.Histogram
	next           Aggregator
}

func NewMetricsMiddleware(next Aggregator) *MetricsMiddleware {
	errCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "aggregate",
	})
	errCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_error_counter",
		Name:      "calculate",
	})
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate",
	})
	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "calculate",
	})
	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "calculate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	return &MetricsMiddleware{
		next:           next,
		reqCounterAgg:  reqCounterAgg,
		reqCounterCalc: reqCounterCalc,
		reqLatencyAgg:  reqLatencyAgg,
		reqLatencyCalc: reqLatencyCalc,
		errCounterAgg:  errCounterAgg,
		errCounterCalc: errCounterCalc,
	}
}

func (m *MetricsMiddleware) AggregateDistance(dist types.Distance) (err error) {
	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(time.Since(start).Seconds())
		m.reqCounterAgg.Inc()
		if err != nil {
			m.errCounterAgg.Inc()
		}
	}(time.Now())
	err = m.next.AggregateDistance(dist)
	return
}

func (m *MetricsMiddleware) CalculateInvoice(id int) (invoice *types.Invoice, err error) {
	defer func(start time.Time) {
		m.reqLatencyCalc.Observe(time.Since(start).Seconds())
		m.reqCounterCalc.Inc()
		if err != nil {
			m.errCounterCalc.Inc()
		}
	}(time.Now())
	invoice, err = m.next.CalculateInvoice(id)
	return
}

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
