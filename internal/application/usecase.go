package application

import (
	"anyker/config"
	"anyker/internal/domain"
	"context"
	"github.com/rs/zerolog/log"
	"strings"
)

// MessageUsecase is the implementation of the MessageUseCase.
type MessageUsecase struct {
	forwardRepository  domain.ForwardRepository
	consumerRepository domain.ConsumerRepository
	config             config.Config
}

// NewMessageService creates a new MessageUsecase with the given repositories.
func NewMessageService(
	config config.Config,
	forwardRepository domain.ForwardRepository,
	consumerRepository domain.ConsumerRepository) domain.MessageUseCase {
	return &MessageUsecase{
		forwardRepository:  forwardRepository,
		consumerRepository: consumerRepository,
		config:             config,
	}
}

// Forward forwards a message using the forward repository.
func (u *MessageUsecase) Forward(ctx context.Context, message domain.Message) error {
	origin, _ := u.getOriginAndRoutingID(message.Key)
	if u.config.Origin != origin {
		log.Debug().Msgf("message origin: %s discarded", origin)
		return nil
	}
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

// getOriginAndRoutingID
func (u *MessageUsecase) getOriginAndRoutingID(key string) (origin string, routingId string) {
	keyParts := strings.SplitN(key, ":", 2)

	if len(keyParts) == 2 {
		origin = keyParts[0]
		routingId = keyParts[1]
	} else {
		// toda la key la tomamos como origin, userID queda vacío
		origin = keyParts[0]
	}
	return
}
