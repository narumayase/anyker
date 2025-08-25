package domain

import "context"

type ConsumerRepository interface {
	Consume(ctx context.Context, messages chan<- *Message) error
	Close() error
}

type ForwardRepository interface {
	Forward(ctx context.Context, message *Message) error
}
