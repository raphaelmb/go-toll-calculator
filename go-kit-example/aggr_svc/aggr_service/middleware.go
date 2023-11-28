package aggrservice

import (
	"context"

	"github.com/raphaelmb/go-toll-calculator/types"
)

type Middleware func(Service) Service

type loggingMiddleware struct {
	next Service
}

func newLoggingMiddelware() Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			next: next,
		}
	}
}

func (mw loggingMiddleware) Aggregate(_ context.Context, dist types.Distance) error {
	return nil
}

func (mw loggingMiddleware) Calculate(_ context.Context, dist int) (*types.Invoice, error) {
	return nil, nil
}

type instrumentationMiddleware struct {
	next Service
}

func newinstrumentationMiddelware() Middleware {
	return func(next Service) Service {
		return instrumentationMiddleware{
			next: next,
		}
	}
}

func (mw instrumentationMiddleware) Aggregate(_ context.Context, dist types.Distance) error {
	return nil
}

func (mw instrumentationMiddleware) Calculate(_ context.Context, id int) (*types.Invoice, error) {
	return nil, nil
}
