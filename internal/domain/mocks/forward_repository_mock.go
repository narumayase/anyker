package mocks

import (
	"anyker/internal/domain"
	"context"
	"github.com/stretchr/testify/mock"
)

type MockForwardRepository struct {
	mock.Mock
}

func (m *MockForwardRepository) Forward(ctx context.Context, message domain.Message) error {
	args := m.Called(ctx, message)
	return args.Error(0)
}
