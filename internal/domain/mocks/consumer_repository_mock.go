package mocks

import (
	"anyker/internal/domain"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockConsumerRepository struct {
	mock.Mock
}

func (m *MockConsumerRepository) Consume(ctx context.Context, messages chan<- *domain.Message) error {
	args := m.Called(ctx, messages)
	return args.Error(0)
}

func (m *MockConsumerRepository) Close() error {
	args := m.Called()
	return args.Error(0)
}
