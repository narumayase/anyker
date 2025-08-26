package repository

import (
	"anyker/config"
	"anyker/internal/domain"
	clientmocks "anyker/internal/infrastructure/client/mocks"
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForwardRepositoryImpl_Forward(t *testing.T) {
	mockHTTPClient := new(clientmocks.MockHTTPClient)
	cfg := config.Config{APIEndpoint: "http://localhost:8080"}
	repo := NewForwardRepository(cfg, mockHTTPClient)

	ctx := context.Background()
	msg := &domain.Message{Content: []byte(`{"key":"value"}`)}

	t.Run("success", func(t *testing.T) {
		mockResponse := clientmocks.CreateMockResponse(http.StatusOK, `{"status":"ok"}`)
		mockHTTPClient.On("Post", ctx, msg.Content, cfg.APIEndpoint).Return(mockResponse, nil).Once()

		err := repo.Forward(ctx, msg)

		assert.NoError(t, err)
		mockHTTPClient.AssertExpectations(t)
	})

	t.Run("http client error", func(t *testing.T) {
		expectedErr := errors.New("http client error")
		mockHTTPClient.On("Post", ctx, msg.Content, cfg.APIEndpoint).Return(nil, expectedErr).Once()

		err := repo.Forward(ctx, msg)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockHTTPClient.AssertExpectations(t)
	})

	t.Run("unexpected status code", func(t *testing.T) {
		mockResponse := clientmocks.CreateMockResponse(http.StatusInternalServerError, `{"error":"internal server error"}`)
		mockHTTPClient.On("Post", ctx, msg.Content, cfg.APIEndpoint).Return(mockResponse, nil).Once()

		err := repo.Forward(ctx, msg)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "unexpected status code: 500")
		mockHTTPClient.AssertExpectations(t)
	})
}
