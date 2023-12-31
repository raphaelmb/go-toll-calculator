package main

import (
	"context"
	"encoding/json"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/raphaelmb/go-toll-calculator/aggregator/client"
	"github.com/raphaelmb/go-toll-calculator/types"
	"github.com/sirupsen/logrus"
)

type KafkaConsumer struct {
	consumer         *kafka.Consumer
	isRunning        bool
	calcService      CalculatorServicer
	aggregatorClient client.Client
}

func NewKafkaConsumer(topic string, svc CalculatorServicer, client client.Client) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9093",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, err
	}

	c.SubscribeTopics([]string{topic}, nil)

	return &KafkaConsumer{
		consumer:         c,
		calcService:      svc,
		aggregatorClient: client,
	}, nil
}

func (c *KafkaConsumer) Start() {
	logrus.Info("kafka transport started")
	c.isRunning = true
	c.readMessageLoop()
}

func (c *KafkaConsumer) readMessageLoop() {
	for c.isRunning {
		msg, err := c.consumer.ReadMessage(-1)
		if err != nil {
			logrus.Errorf("kafka consume error: %s", err)
			continue
		}
		var data types.OBUData
		if err := json.Unmarshal(msg.Value, &data); err != nil {
			logrus.Errorf("json serialization error: %s", err)
			continue
		}
		distance, err := c.calcService.CalculateDistance(data)
		if err != nil {
			logrus.Errorf("calculation error: %s", err)
			continue
		}
		req := &types.AggregateRequest{
			Value: distance,
			Unix:  time.Now().Unix(),
			ObuID: int32(data.OBUID),
		}
		if err := c.aggregatorClient.Aggregate(context.Background(), req); err != nil {
			logrus.Error("aggregate error:", err)
			continue
		}
	}

}
