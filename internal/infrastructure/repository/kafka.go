package repository

import (
	"anyker/config"
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"time"

	"anyker/internal/domain"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// KafkaConsumer defines the interface for the Kafka consumer, to allow for mocking.
type KafkaConsumer interface {
	Subscribe(topic string, rebalanceCb kafka.RebalanceCb) error
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	Close() error
}

// Consumer is a Kafka consumer that implements the ConsumerRepository interface.
type Consumer struct {
	consumer KafkaConsumer
	topic    string
}

// NewConsumer creates a new Kafka consumer.
func NewConsumer(config config.Config) (domain.ConsumerRepository, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": config.KafkaBroker,
		"group.id":          config.KafkaGroupID,
		"auto.offset.reset": "latest",
		// latest para ignorar los mensajes viejos, earliest para lo contrario
		// TODO parametrizar el offset
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
			msg, err := c.consumer.ReadMessage(1000 * time.Millisecond)
			if err != nil {
				if kafkaErr, ok := err.(kafka.Error); ok && kafkaErr.Code() == kafka.ErrTimedOut {
					continue
				}
				if errors.Is(err, context.Canceled) {
					return nil // Graceful exit on context cancellation
				}
				return fmt.Errorf("failed to read message: %w", err)
			}
			headers := make(map[string]string)
			if msg.Headers != nil {
				for _, h := range msg.Headers {
					headers[h.Key] = string(h.Value)
				}
			}
			log.Debug().Msgf("headers received from Kafka message %v", headers)
			log.Debug().Msgf("key receive from Kafka message %v", string(msg.Key))
			log.Debug().Msgf("payload receive from Kafka message %v", string(msg.Value))

			messages <- &domain.Message{
				Content: msg.Value,
				Headers: headers,
				Key:     string(msg.Key),
			}
		}
	}
}

// Close closes the Kafka consumer.
func (c *Consumer) Close() error {
	return c.consumer.Close()
}
