package domain

import "context"

// ConsumerRepository defines the interface for consuming messages.
type ConsumerRepository interface {
	// Consume starts consuming messages and sends them to the provided channel.
	// This method should block until the context is cancelled or an error occurs.
	Consume(ctx context.Context, messages chan<- *Message) error
	// Close closes the consumer connection.
	Close() error
}

// ForwardRepository defines the interface for forwarding messages.
type ForwardRepository interface {
	// Forward forwards a message to a downstream service.
	Forward(ctx context.Context, message Message) error
}
