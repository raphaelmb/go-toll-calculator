package main

import (
	"context"
	"log"
	"time"

	"github.com/raphaelmb/go-toll-calculator/aggregator/client"
	"github.com/raphaelmb/go-toll-calculator/types"
)

func main() {
	c, err := client.NewGRPCClient(":3001")
	if err != nil {
		log.Fatal(err)
	}
	if err := c.Aggregate(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 58.77,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		log.Fatal(err)
	}
}
