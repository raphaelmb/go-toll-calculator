package aggrendpoint

import (
	"context"
	"time"

	"github.com/go-kit/kit/circuitbreaker"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/kit/ratelimit"
	"github.com/go-kit/log"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	aggrservice "github.com/raphaelmb/go-toll-calculator/go-kit-example/aggr_svc/aggr_service"
	"github.com/raphaelmb/go-toll-calculator/types"
	"github.com/sony/gobreaker"
	"golang.org/x/time/rate"
)

type Set struct {
	AggregateEndpoint endpoint.Endpoint
	CalculateEndpoint endpoint.Endpoint
}

func New(svc aggrservice.Service, logger log.Logger) Set {
	duration := prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
		Namespace: "toll_calculator",
		Subsystem: "aggservice",
		Name:      "request_duration_seconds",
		Help:      "Request duration is seconds",
	}, []string{"method", "success"})

	var aggregateEndpoint endpoint.Endpoint
	{
		aggregateEndpoint = MakeAggregateEndpoint(svc)
		aggregateEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Every(time.Second), 1))(aggregateEndpoint)
		aggregateEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(aggregateEndpoint)
		aggregateEndpoint = LoggingMiddleware(log.With(logger, "method", "Aggregate"))(aggregateEndpoint)
		aggregateEndpoint = InstrumentingMiddleware(duration.With("method", "Aggregate"))(aggregateEndpoint)
	}
	var calculateEndpoint endpoint.Endpoint
	{
		calculateEndpoint = MakeConcatEndpoint(svc)
		calculateEndpoint = ratelimit.NewErroringLimiter(rate.NewLimiter(rate.Limit(1), 100))(calculateEndpoint)
		calculateEndpoint = circuitbreaker.Gobreaker(gobreaker.NewCircuitBreaker(gobreaker.Settings{}))(calculateEndpoint)
		calculateEndpoint = LoggingMiddleware(log.With(logger, "method", "Calculate"))(calculateEndpoint)
		calculateEndpoint = InstrumentingMiddleware(duration.With("method", "Calculate"))(calculateEndpoint)
	}
	return Set{
		AggregateEndpoint: aggregateEndpoint,
		CalculateEndpoint: calculateEndpoint,
	}
}

type AggregateRequest struct {
	Value float64 `json:"value"`
	OBUID int     `json:"obuID"`
	Unix  int64   `json:"unix"`
}

type AggregateResponse struct {
	Err error `json:"err"`
}

type CalculateResponse struct {
	OBUID         int     `json:"obuID"`
	TotalDistance float64 `json:"totalDistance"`
	TotalAmount   float64 `json:"totalAmount"`
	Err           error   `json:"err"`
}

type CalculateRequest struct {
	OBUID int `json:"obuID"`
}

func (s Set) Aggregate(ctx context.Context, dist types.Distance) error {
	_, err := s.AggregateEndpoint(ctx, AggregateRequest{
		OBUID: dist.OBUID,
		Value: dist.Value,
		Unix:  dist.Unix,
	})
	return err
}

func (s Set) Calculate(ctx context.Context, id int) (*types.Invoice, error) {
	resp, err := s.CalculateEndpoint(ctx, CalculateRequest{
		OBUID: id})
	if err != nil {
		return nil, err
	}
	result := resp.(CalculateResponse)
	return &types.Invoice{
		OBUID:         result.OBUID,
		TotalDistance: result.TotalDistance,
		TotalAmount:   result.TotalAmount,
	}, nil
}

func MakeAggregateEndpoint(s aggrservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(AggregateRequest)
		err = s.Aggregate(ctx, types.Distance{
			OBUID: req.OBUID,
			Value: req.Value,
			Unix:  req.Unix,
		})
		return AggregateResponse{Err: err}, nil
	}
}

func MakeConcatEndpoint(s aggrservice.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(CalculateRequest)
		invoice, err := s.Calculate(ctx, req.OBUID)
		return CalculateResponse{
			Err:           err,
			OBUID:         invoice.OBUID,
			TotalDistance: invoice.TotalDistance,
			TotalAmount:   invoice.TotalAmount,
		}, nil
	}
}
