package repository

import (
	"anyker/config"
	"anyker/internal/domain"
	"anyker/internal/infrastructure/repository/mocks"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewConsumer(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		cfg := config.Config{
			KafkaBroker:  "localhost:9092",
			KafkaGroupID: "test-group",
			KafkaTopic:   "test-topic",
		}
		// Since kafka.NewConsumer is a concrete function, we can't mock it directly.
		// We'll assume it works correctly for this test and focus on the NewConsumer wrapper.
		// A more advanced test might involve dependency injection for kafka.NewConsumer.
		consumer, err := NewConsumer(cfg)
		assert.NoError(t, err)
		assert.NotNil(t, consumer)
	})

	// Note: Testing the failure of kafka.NewConsumer directly is hard without
	// mocking the kafka package itself. This test focuses on the happy path
	// and assumes the underlying kafka library behaves as expected.
}

func TestConsumer_Consume(t *testing.T) {
	t.Run("successful message consumption", func(t *testing.T) {
		mockKafkaConsumer := new(mocks.KafkaConsumer)
		consumer := &Consumer{
			consumer: mockKafkaConsumer,
			topic:    "test-topic",
		}
		messagesChan := make(chan *domain.Message, 10)
		ctx, cancel := context.WithCancel(context.Background())

		defer mockKafkaConsumer.AssertExpectations(t)
		defer cancel()

		mockKafkaConsumer.On("Subscribe", "test-topic", mock.Anything).Return(nil).Once()

		// Simulate two messages, then a timeout, then context cancellation
		msg1 := &kafka.Message{
			Value:   []byte("message1"),
			Headers: []kafka.Header{{Key: "correlation_id", Value: []byte("123")}},
			Key:     []byte("key1"),
		}
		msg2 := &kafka.Message{
			Value:   []byte("message2"),
			Headers: []kafka.Header{{Key: "correlation_id", Value: []byte("456")}},
			Key:     []byte("key2"),
		}

		mockKafkaConsumer.On("ReadMessage", mock.AnythingOfType("time.Duration")).Return(msg1, nil).Once()
		mockKafkaConsumer.On("ReadMessage", mock.AnythingOfType("time.Duration")).Return(msg2, nil).Once()
		mockKafkaConsumer.On("ReadMessage", mock.AnythingOfType("time.Duration")).Return(nil, kafka.NewError(kafka.ErrTimedOut, "Local: Timed out", false)).Once()
		mockKafkaConsumer.On("ReadMessage", mock.AnythingOfType("time.Duration")).Return(nil, context.Canceled).Maybe() // For graceful exit

		go func() {
			err := consumer.Consume(ctx, messagesChan)
			assert.NoError(t, err) // Consume should exit gracefully on context cancellation
		}()

		// Read messages from the channel
		receivedMsg1 := <-messagesChan
		assert.Equal(t, "message1", string(receivedMsg1.Content))
		assert.Equal(t, "123", string(receivedMsg1.Headers["correlation_id"]))
		assert.Equal(t, "key1", receivedMsg1.Key)

		receivedMsg2 := <-messagesChan
		assert.Equal(t, "message2", string(receivedMsg2.Content))
		assert.Equal(t, "456", string(receivedMsg2.Headers["correlation_id"]))
		assert.Equal(t, "key2", receivedMsg2.Key)

		// Give some time for the consumer to process timeout and context cancellation
		time.Sleep(100 * time.Millisecond)
		cancel()                           // Cancel the context to stop the consumer loop
		time.Sleep(100 * time.Millisecond) // Give time for goroutine to finish
	})

	t.Run("subscribe error", func(t *testing.T) {
		mockKafkaConsumer := new(mocks.KafkaConsumer)
		consumer := &Consumer{
			consumer: mockKafkaConsumer,
			topic:    "test-topic",
		}
		messagesChan := make(chan *domain.Message, 10)
		ctx, cancel := context.WithCancel(context.Background())

		defer mockKafkaConsumer.AssertExpectations(t)
		defer cancel()

		expectedErr := errors.New("failed to subscribe")
		mockKafkaConsumer.On("Subscribe", "test-topic", mock.Anything).Return(expectedErr).Once()

		err := consumer.Consume(ctx, messagesChan)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to subscribe")
	})

	t.Run("read message error (non-timeout)", func(t *testing.T) {
		mockKafkaConsumer := new(mocks.KafkaConsumer)
		consumer := &Consumer{
			consumer: mockKafkaConsumer,
			topic:    "test-topic",
		}
		messagesChan := make(chan *domain.Message, 10)
		ctx, cancel := context.WithCancel(context.Background())

		defer mockKafkaConsumer.AssertExpectations(t)
		defer cancel()

		mockKafkaConsumer.On("Subscribe", "test-topic", mock.Anything).Return(nil).Once()
		expectedErr := errors.New("kafka read error")
		mockKafkaConsumer.On("ReadMessage", mock.AnythingOfType("time.Duration")).Return(nil, expectedErr).Once()

		err := consumer.Consume(ctx, messagesChan)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to read message")
		assert.Contains(t, err.Error(), expectedErr.Error())
	})

	t.Run("message with no headers", func(t *testing.T) {
		mockKafkaConsumer := new(mocks.KafkaConsumer)
		consumer := &Consumer{
			consumer: mockKafkaConsumer,
			topic:    "test-topic",
		}
		messagesChan := make(chan *domain.Message, 10)
		ctx, cancelFunc := context.WithCancel(context.Background())

		defer mockKafkaConsumer.AssertExpectations(t)
		defer cancelFunc()

		mockKafkaConsumer.On("Subscribe", "test-topic", mock.Anything).Return(nil).Once()
		msg := &kafka.Message{
			Value: []byte("no_headers_message"),
			Key:   []byte("key_no_headers"),
		}
		mockKafkaConsumer.On("ReadMessage", mock.AnythingOfType("time.Duration")).Return(msg, nil).Once()
		mockKafkaConsumer.On("ReadMessage", mock.AnythingOfType("time.Duration")).Return(nil, context.Canceled).Maybe()

		go func() {
			err := consumer.Consume(ctx, messagesChan)
			assert.NoError(t, err)
		}()

		receivedMsg := <-messagesChan
		assert.Equal(t, "no_headers_message", string(receivedMsg.Content))
		assert.Empty(t, receivedMsg.Headers)
		assert.Equal(t, "key_no_headers", receivedMsg.Key)

		time.Sleep(100 * time.Millisecond)
		cancelFunc()
		time.Sleep(100 * time.Millisecond)
	})
}

func TestConsumer_Close(t *testing.T) {
	mockKafkaConsumer := new(mocks.KafkaConsumer)
	consumer := &Consumer{
		consumer: mockKafkaConsumer,
	}

	t.Run("successful close", func(t *testing.T) {
		defer mockKafkaConsumer.AssertExpectations(t)
		mockKafkaConsumer.On("Close").Return(nil).Once()

		err := consumer.Close()
		assert.NoError(t, err)
	})

	t.Run("close error", func(t *testing.T) {
		defer mockKafkaConsumer.AssertExpectations(t)
		mockKafkaConsumer.ExpectedCalls = nil // Clear previous expectations

		expectedErr := errors.New("failed to close consumer")
		mockKafkaConsumer.On("Close").Return(expectedErr).Once()

		err := consumer.Close()
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
