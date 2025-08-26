package repository

import (
	"anyker/config"
	"context"
	"fmt"

	"anyker/internal/domain"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// Consumer is a Kafka consumer that implements the ConsumerRepository interface.
type Consumer struct {
	consumer *kafka.Consumer
	topic    string
}

// NewConsumer creates a new Kafka consumer.
func NewConsumer(config config.Config) (domain.ConsumerRepository, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBroker,
		"group.id":          config.KafkaGroupID,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &Consumer{
		consumer: c,
		topic:    config.KafkaTopic,
	}, nil
}

// Consume consumes messages from Kafka and sends them to the provided channel.
func (c *Consumer) Consume(ctx context.Context, messages chan<- *domain.Message) error {
	defer close(messages)

	err := c.consumer.Subscribe(c.topic, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			msg, err := c.consumer.ReadMessage(1000)
			if err != nil {
				if err.(kafka.Error).Code() == kafka.ErrTimedOut {
					continue
				}
				return fmt.Errorf("failed to read message: %w", err)
			}
			messages <- &domain.Message{Content: msg.Value}
		}
	}
}

// Close closes the Kafka consumer.
func (c *Consumer) Close() error {
	return c.consumer.Close()
}
