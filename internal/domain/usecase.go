package domain

import "context"

// MessageUseCase is the interface for the message use case.
type MessageUseCase interface {
	// Forward forwards a message to a downstream service.
	Forward(ctx context.Context, message Message) error

	// Consume consumes messages from Kafka and sends them to the provided channel.
	Consume(ctx context.Context, messages chan<- *Message) error

	// Close closes the use case and its underlying resources.
	Close() error
}
