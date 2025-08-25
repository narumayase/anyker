package domain

import "context"

// MessageUseCase is the interface for the message use case.
type MessageUseCase interface {
	Forward(ctx context.Context, message *Message) error

	// Consume consumes messages from Kafka and sends them to the provided channel.
	Consume(ctx context.Context, messages chan<- *Message) error

	Close() error
}
