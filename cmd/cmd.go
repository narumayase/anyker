package cmd

import (
	"anyker/internal/domain"
	"context"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
	"syscall"
)

// Run starts the worker, which consumes messages from Kafka and forwards them.
// It also handles graceful shutdown on SIGINT or SIGTERM signals.
func Run(usecase domain.MessageUseCase) {
	// consume messages
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		log.Info().Msg("Shutting down...")
		cancel()
	}()

	messages := make(chan *domain.Message)

	go func() {
		if err := usecase.Consume(ctx, messages); err != nil {
			log.Error().Err(err).Msg("failed to consume messages")
		}
	}()
	log.Info().Msg("Worker listening to Kafka...")

	defer usecase.Close()

	// forward messages
	for message := range messages {
		if err := usecase.Forward(ctx, message); err != nil {
			log.Error().Err(err).Msg("failed to forward message")
		}
	}
	log.Info().Msg("Worker stopped.")
}
