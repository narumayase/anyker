package application

import (
	"anyker/internal/domain"
	"context"
)

// MessageUsecase is the implementation of the MessageUseCase.
type MessageUsecase struct {
	forwardRepository  domain.ForwardRepository
	consumerRepository domain.ConsumerRepository
}

// NewMessageService creates a new MessageUsecase with the given repositories.
func NewMessageService(forwardRepository domain.ForwardRepository,
	consumerRepository domain.ConsumerRepository) domain.MessageUseCase {
	return &MessageUsecase{
		forwardRepository:  forwardRepository,
		consumerRepository: consumerRepository,
	}
}

// Forward forwards a message using the forward repository.
func (u *MessageUsecase) Forward(ctx context.Context, message *domain.Message) error {
	return u.forwardRepository.Forward(ctx, message)
}

// Consume consumes messages from the consumer repository.
func (u *MessageUsecase) Consume(ctx context.Context, messages chan<- *domain.Message) error {
	return u.consumerRepository.Consume(ctx, messages)
}

// Close closes the underlying consumer repository.
func (u *MessageUsecase) Close() error {
	return u.consumerRepository.Close()
}
