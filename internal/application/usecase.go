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
	if u.config.Origin == "" {
		log.Debug().Msg("all messages will be read")
	} else {
		origin, _ := u.getOriginAndRoutingID(message.Key)
		if u.config.Origin != origin {
			log.Debug().Msgf("message origin: %s discarded", origin)
			return nil
		}
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

// getOriginAndRoutingID extracts the origin and routing ID from a given key.
// If the key does not contain a colon, the entire key is treated as the origin, and routingId remains empty.
func (u *MessageUsecase) getOriginAndRoutingID(key string) (origin string, routingId string) {
	keyParts := strings.SplitN(key, ":", 2)

	if len(keyParts) == 2 {
		origin = keyParts[0]
		routingId = keyParts[1]
	} else {
		// the entire key is taken as origin, userID remains empty
		origin = keyParts[0]
	}
	return
}
