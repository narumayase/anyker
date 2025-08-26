package mocks

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/mock"
)

type MockKafkaConsumer struct {
	mock.Mock
}

func (m *MockKafkaConsumer) Subscribe(topic string, rebalanceCb kafka.RebalanceCb) error {
	args := m.Called(topic, rebalanceCb)
	return args.Error(0)
}

func (m *MockKafkaConsumer) ReadMessage(timeout time.Duration) (*kafka.Message, error) {
	args := m.Called(timeout)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*kafka.Message), args.Error(1)
}

func (m *MockKafkaConsumer) Close() error {
	args := m.Called()
	return args.Error(0)
}
