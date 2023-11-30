package aggrservice

import (
	"context"
	"time"

	"github.com/go-kit/log"
	"github.com/raphaelmb/go-toll-calculator/types"
)

type Middleware func(Service) Service

type loggingMiddleware struct {
	log  log.Logger
	next Service
}

func newLoggingMiddelware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return loggingMiddleware{
			next: next,
			log:  logger,
		}
	}
}

func (mw loggingMiddleware) Aggregate(ctx context.Context, dist types.Distance) (err error) {
	defer func(start time.Time) {
		mw.log.Log("method", "Aggregate", "took", time.Since(start), "obu", dist.OBUID, "distance", dist.Value, "err", err)
	}(time.Now())
	err = mw.next.Aggregate(ctx, dist)
	return err
}

func (mw loggingMiddleware) Calculate(ctx context.Context, id int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		mw.log.Log("method", "Calculate", "took", time.Since(start), "id", id, "invoice", inv, "err", err)
	}(time.Now())
	inv, err = mw.next.Calculate(ctx, id)
	return
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

func (mw instrumentationMiddleware) Aggregate(ctx context.Context, dist types.Distance) error {
	return mw.next.Aggregate(ctx, dist)
}

func (mw instrumentationMiddleware) Calculate(ctx context.Context, id int) (*types.Invoice, error) {
	return mw.next.Calculate(ctx, id)
}
