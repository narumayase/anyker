package main

import (
	"anyker/cmd"
	"anyker/config"
	"anyker/internal/application"
	"anyker/internal/infrastructure/client"
	"anyker/internal/infrastructure/repository"
	"github.com/rs/zerolog/log"
	"net/http"
)

func main() {
	// Load configuration
	cfg := config.Load()
	log.Info().Str("nanobot_name", cfg.NanobotName).Msg("Starting nanobot")

	// Create http forward client.
	// It's a good practice to set a timeout for HTTP clients in production.
	httpClient := &http.Client{
		Timeout: cfg.HTTPClientTimeout,
	}
	forwardHttpClient := client.NewHttpClient(httpClient, "")

	// Create repositories
	forwardRepository := repository.NewForwardRepository(cfg, forwardHttpClient)

	consumerRepository, err := repository.NewConsumer(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create consumerRepository")
	}

	// Create use case
	messageService := application.NewMessageService(cfg, forwardRepository, consumerRepository)

	cmd.Run(messageService)
}
