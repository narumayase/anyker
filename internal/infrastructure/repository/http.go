package repository

import (
	"anyker/config"
	"anyker/internal/domain"
	"anyker/internal/infrastructure/client"
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"net/http"
)

// ForwardRepositoryImpl implements the domain.ForwardRepository interface using an HTTP client.
type ForwardRepositoryImpl struct {
	config     config.Config
	httpClient client.HttpClient
}

// NewForwardRepository creates a new ForwardRepositoryImpl.
func NewForwardRepository(
	config config.Config,
	httpClient client.HttpClient) domain.ForwardRepository {
	return &ForwardRepositoryImpl{
		httpClient: httpClient,
		config:     config,
	}
}

// Forward forwards a message to the configured API endpoint.
func (f *ForwardRepositoryImpl) Forward(ctx context.Context, message domain.Message) error {
	headers := map[string]string{
		"X-Correlation-ID": string(message.Headers["correlation_id"]),
	}
	resp, err := f.httpClient.Post(ctx, headers, message.Content, f.config.APIEndpoint)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	log.Info().Msgf("API response status: %s", resp.Status)

	return nil
}
