package main

import (
	"log"

	"github.com/raphaelmb/go-toll-calculator/aggregator/client"
)

const (
	kafkaTopic         = "obudata"
	aggregatorEndpoint = "http://127.0.0.1:3000"
)

func main() {
	svc := NewCalculatorService()
	svc = NewLogMiddleware(svc)
	client := client.NewHTTPClient(aggregatorEndpoint)
	// client, err := client.NewGRPCClient(aggregatorEndpoint)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	kafkaConsumer, err := NewKafkaConsumer(kafkaTopic, svc, client)
	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
