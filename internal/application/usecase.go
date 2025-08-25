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

// NewMessageService creates a new MessageUsecase.
func NewMessageService(forwardRepository domain.ForwardRepository,
	consumerRepository domain.ConsumerRepository) domain.MessageUseCase {
	return &MessageUsecase{
		forwardRepository:  forwardRepository,
		consumerRepository: consumerRepository,
	}
}

// Forward forwards a message to the configured API endpoint.
func (u *MessageUsecase) Forward(ctx context.Context, message *domain.Message) error {
	return u.forwardRepository.Forward(ctx, message)
}

func (u *MessageUsecase) Consume(ctx context.Context, messages chan<- *domain.Message) error {
	return u.consumerRepository.Consume(ctx, messages)
}

func (u *MessageUsecase) Close() error {
	return u.consumerRepository.Close()
}
