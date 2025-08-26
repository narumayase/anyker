package application

import (
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
	usecase := NewMessageService(mockForwardRepo, nil)

	ctx := context.Background()
	msg := &domain.Message{Content: []byte("hello")}

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

func TestMessageUsecase_Consume(t *testing.T) {
	mockConsumerRepo := new(mocks.MockConsumerRepository)
	usecase := NewMessageService(nil, mockConsumerRepo)

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
}

func TestMessageUsecase_Close(t *testing.T) {
	mockConsumerRepo := new(mocks.MockConsumerRepository)
	usecase := NewMessageService(nil, mockConsumerRepo)

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