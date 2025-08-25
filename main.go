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

	// Create http forward client
	httpClient := &http.Client{}
	forwardHttpClient := client.NewHttpClient(httpClient, "")

	// Create repositories
	forwardRepository := repository.NewForwardRepository(cfg, forwardHttpClient)

	consumerRepository, err := repository.NewConsumer(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create consumerRepository")
	}
	//defer consumerRepository.Close()
	//TODO hay que cerrar algo?

	// Create use case
	messageService := application.NewMessageService(forwardRepository, consumerRepository)

	cmd.Run(messageService)
}
