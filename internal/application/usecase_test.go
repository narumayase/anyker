package application

import (
	"anyker/config"
	"anyker/internal/domain"
	"anyker/internal/domain/mocks"
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestMessageUsecase_Forward(t *testing.T) {
	mockForwardRepo := new(mocks.MockForwardRepository)
	cfg := config.Config{} // or initialize with specific values if necessary
	usecase := NewMessageService(cfg, mockForwardRepo, nil)

	ctx := context.Background()
	msg := domain.Message{Content: []byte("hello")}

	t.Run("success", func(t *testing.T) {
		mockForwardRepo.On("Forward", ctx, msg).Return(nil).Once()

		err := usecase.Forward(ctx, msg)

		assert.NoError(t, err)
		mockForwardRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("forward error")
		mockForwardRepo.On("Forward", ctx, msg).Return(expectedErr).Once()

		err := usecase.Forward(ctx, msg)

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockForwardRepo.AssertExpectations(t)
	})
}

func TestMessageUsecase_Forward_OriginFiltering(t *testing.T) {
	mockForwardRepo := new(mocks.MockForwardRepository)
	ctx := context.Background()

	t.Run("empty origin config - all messages processed", func(t *testing.T) {
		cfg := config.Config{Origin: ""} // Empty origin means process all
		usecase := NewMessageService(cfg, mockForwardRepo, nil)

		msg := domain.Message{
			Key:     "any-origin:routing-id",
			Content: []byte("test message"),
		}

		mockForwardRepo.On("Forward", ctx, msg).Return(nil).Once()

		err := usecase.Forward(ctx, msg)

		assert.NoError(t, err)
		mockForwardRepo.AssertExpectations(t)
	})

	t.Run("matching origin - message forwarded", func(t *testing.T) {
		cfg := config.Config{Origin: "service-a"}
		usecase := NewMessageService(cfg, mockForwardRepo, nil)

		msg := domain.Message{
			Key:     "service-a:user-123",
			Content: []byte("test message"),
		}

		mockForwardRepo.On("Forward", ctx, msg).Return(nil).Once()

		err := usecase.Forward(ctx, msg)

		assert.NoError(t, err)
		mockForwardRepo.AssertExpectations(t)
	})

	t.Run("non-matching origin - message discarded", func(t *testing.T) {
		cfg := config.Config{Origin: "service-a"}
		usecase := NewMessageService(cfg, mockForwardRepo, nil)

		msg := domain.Message{
			Key:     "service-b:user-123",
			Content: []byte("test message"),
		}

		// Should not call Forward because origin doesn't match
		// mockForwardRepo has no expectations, so if called it will fail

		err := usecase.Forward(ctx, msg)

		assert.NoError(t, err) // No error, but also not forwarded
		mockForwardRepo.AssertExpectations(t)
	})

	t.Run("key without colon - origin is entire key", func(t *testing.T) {
		cfg := config.Config{Origin: "simple-key"}
		usecase := NewMessageService(cfg, mockForwardRepo, nil)

		msg := domain.Message{
			Key:     "simple-key",
			Content: []byte("test message"),
		}

		mockForwardRepo.On("Forward", ctx, msg).Return(nil).Once()

		err := usecase.Forward(ctx, msg)

		assert.NoError(t, err)
		mockForwardRepo.AssertExpectations(t)
	})

	t.Run("key without colon - non-matching origin", func(t *testing.T) {
		cfg := config.Config{Origin: "expected-key"}
		usecase := NewMessageService(cfg, mockForwardRepo, nil)

		msg := domain.Message{
			Key:     "different-key",
			Content: []byte("test message"),
		}

		// No expectations - shouldn't be called
		err := usecase.Forward(ctx, msg)

		assert.NoError(t, err)
		mockForwardRepo.AssertExpectations(t)
	})
}

func TestMessageUsecase_getOriginAndRoutingID(t *testing.T) {
	usecase := &MessageUsecase{}

	t.Run("key with colon", func(t *testing.T) {
		origin, routingID := usecase.getOriginAndRoutingID("service-a:user-123")

		assert.Equal(t, "service-a", origin)
		assert.Equal(t, "user-123", routingID)
	})

	t.Run("key with multiple colons", func(t *testing.T) {
		origin, routingID := usecase.getOriginAndRoutingID("service-a:user-123:extra-data")

		assert.Equal(t, "service-a", origin)
		assert.Equal(t, "user-123:extra-data", routingID) // SplitN(2) keeps the rest together
	})

	t.Run("key without colon", func(t *testing.T) {
		origin, routingID := usecase.getOriginAndRoutingID("simple-key")

		assert.Equal(t, "simple-key", origin)
		assert.Equal(t, "", routingID)
	})

	t.Run("empty key", func(t *testing.T) {
		origin, routingID := usecase.getOriginAndRoutingID("")

		assert.Equal(t, "", origin)
		assert.Equal(t, "", routingID)
	})

	t.Run("key starting with colon", func(t *testing.T) {
		origin, routingID := usecase.getOriginAndRoutingID(":routing-only")

		assert.Equal(t, "", origin)
		assert.Equal(t, "routing-only", routingID)
	})

	t.Run("key ending with colon", func(t *testing.T) {
		origin, routingID := usecase.getOriginAndRoutingID("service-a:")

		assert.Equal(t, "service-a", origin)
		assert.Equal(t, "", routingID)
	})
}

func TestMessageUsecase_NewMessageService(t *testing.T) {
	mockForwardRepo := new(mocks.MockForwardRepository)
	mockConsumerRepo := new(mocks.MockConsumerRepository)
	cfg := config.Config{Origin: "test-origin"}

	usecase := NewMessageService(cfg, mockForwardRepo, mockConsumerRepo)

	assert.NotNil(t, usecase)

	// Verify that it implements the interface
	var _ domain.MessageUseCase = usecase

	// Verify that it's the correct type
	concreteUsecase, ok := usecase.(*MessageUsecase)
	assert.True(t, ok)

	// Verify that fields were assigned correctly
	assert.Equal(t, cfg, concreteUsecase.config)
	assert.Equal(t, mockForwardRepo, concreteUsecase.forwardRepository)
	assert.Equal(t, mockConsumerRepo, concreteUsecase.consumerRepository)
}

func TestMessageUsecase_Consume(t *testing.T) {
	mockConsumerRepo := new(mocks.MockConsumerRepository)
	cfg := config.Config{} // or initialize with specific values if necessary
	usecase := NewMessageService(cfg, nil, mockConsumerRepo)

	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockConsumerRepo.On("Consume", ctx, mock.AnythingOfType("chan<- *domain.Message")).Return(nil).Once()

		err := usecase.Consume(ctx, make(chan *domain.Message))

		assert.NoError(t, err)
		mockConsumerRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("consume error")
		mockConsumerRepo.On("Consume", ctx, mock.AnythingOfType("chan<- *domain.Message")).Return(expectedErr).Once()

		err := usecase.Consume(ctx, make(chan *domain.Message))

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockConsumerRepo.AssertExpectations(t)
	})

	t.Run("nil channel", func(t *testing.T) {
		mockConsumerRepo.On("Consume", ctx, mock.AnythingOfType("chan<- *domain.Message")).Return(nil).Once()

		err := usecase.Consume(ctx, nil)

		assert.NoError(t, err)
		mockConsumerRepo.AssertExpectations(t)
	})

	t.Run("cancelled context", func(t *testing.T) {
		cancelledCtx, cancel := context.WithCancel(ctx)
		cancel() // Cancel immediately

		expectedErr := context.Canceled
		mockConsumerRepo.On("Consume", cancelledCtx, mock.AnythingOfType("chan<- *domain.Message")).Return(expectedErr).Once()

		err := usecase.Consume(cancelledCtx, make(chan *domain.Message))

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockConsumerRepo.AssertExpectations(t)
	})
}

func TestMessageUsecase_Close(t *testing.T) {
	mockConsumerRepo := new(mocks.MockConsumerRepository)
	cfg := config.Config{} // or initialize with specific values if necessary
	usecase := NewMessageService(cfg, nil, mockConsumerRepo)

	t.Run("success", func(t *testing.T) {
		mockConsumerRepo.On("Close").Return(nil).Once()

		err := usecase.Close()

		assert.NoError(t, err)
		mockConsumerRepo.AssertExpectations(t)
	})

	t.Run("error", func(t *testing.T) {
		expectedErr := errors.New("close error")
		mockConsumerRepo.On("Close").Return(expectedErr).Once()

		err := usecase.Close()

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
		mockConsumerRepo.AssertExpectations(t)
	})
}

// Integration test that verifies the complete flow
func TestMessageUsecase_Integration(t *testing.T) {
	mockForwardRepo := new(mocks.MockForwardRepository)
	mockConsumerRepo := new(mocks.MockConsumerRepository)
	cfg := config.Config{Origin: "test-service"}

	usecase := NewMessageService(cfg, mockForwardRepo, mockConsumerRepo)
	ctx := context.Background()

	t.Run("complete workflow", func(t *testing.T) {
		// Setup mocks for a complete flow
		messages := make(chan *domain.Message, 1)

		// Mock consume
		mockConsumerRepo.On("Consume", ctx, mock.AnythingOfType("chan<- *domain.Message")).Return(nil)

		// Mock forward for message that should be processed
		msg := domain.Message{
			Key:     "test-service:user-123",
			Content: []byte("test content"),
		}
		mockForwardRepo.On("Forward", ctx, msg).Return(nil)

		// Mock close
		mockConsumerRepo.On("Close").Return(nil)

		// Execute consume
		err := usecase.Consume(ctx, messages)
		assert.NoError(t, err)

		// Execute forward
		err = usecase.Forward(ctx, msg)
		assert.NoError(t, err)

		// Execute close
		err = usecase.Close()
		assert.NoError(t, err)

		mockConsumerRepo.AssertExpectations(t)
		mockForwardRepo.AssertExpectations(t)
	})
}

// Benchmarks to measure performance
func BenchmarkMessageUsecase_getOriginAndRoutingID(b *testing.B) {
	usecase := &MessageUsecase{}

	b.Run("with_colon", func(b *testing.B) {
		key := "service-a:user-12345"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = usecase.getOriginAndRoutingID(key)
		}
	})

	b.Run("without_colon", func(b *testing.B) {
		key := "simple-key"
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, _ = usecase.getOriginAndRoutingID(key)
		}
	})
}

func BenchmarkMessageUsecase_Forward(b *testing.B) {
	mockForwardRepo := new(mocks.MockForwardRepository)
	cfg := config.Config{Origin: ""}
	usecase := NewMessageService(cfg, mockForwardRepo, nil)

	ctx := context.Background()
	msg := domain.Message{
		Key:     "service:user",
		Content: []byte("benchmark message"),
	}

	mockForwardRepo.On("Forward", ctx, msg).Return(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = usecase.Forward(ctx, msg)
	}
}
