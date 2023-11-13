package client

import (
	"context"

	"github.com/raphaelmb/go-toll-calculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}
