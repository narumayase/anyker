package config

import (
	"os"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		expected     string
	}{
		{
			name:         "environment variable exists",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "env_value",
			expected:     "env_value",
		},
		{
			name:         "environment variable does not exist",
			key:          "NON_EXISTENT_KEY",
			defaultValue: "default_value",
			envValue:     "",
			expected:     "default_value",
		},
		{
			name:         "empty environment variable",
			key:          "EMPTY_KEY",
			defaultValue: "default",
			envValue:     "",
			expected:     "default",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment
			os.Unsetenv(tt.key)

			// Set environment variable if provided
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			result := getEnv(tt.key, tt.defaultValue)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_DefaultValues(t *testing.T) {
	// Clean up environment variables
	envVars := []string{"KAFKA_BROKER", "KAFKA_TOPIC", "KAFKA_GROUP_ID", "API_ENDPOINT", "NANOBOT_NAME", "HTTP_CLIENT_TIMEOUT_SECONDS"}
	for _, env := range envVars {
		os.Unsetenv(env)
	}

	t.Run("config structure", func(t *testing.T) {
		config := Load()

		assert.Equal(t, "localhost:9092", config.KafkaBroker)
		assert.Equal(t, "anyker-topic", config.KafkaTopic)
		assert.Equal(t, "anyker-group", config.KafkaGroupID)
		assert.Equal(t, "http://localhost:8080/messages", config.APIEndpoint)
		assert.Equal(t, "anyker-nanobot-1", config.NanobotName)
		assert.Equal(t, 30*time.Second, config.HTTPClientTimeout)
	})
}

func TestConfig_WithEnvironmentVariables(t *testing.T) {
	// Set test environment variables
	testEnvVars := map[string]string{
		"KAFKA_BROKER":                "kafka:9092",
		"KAFKA_TOPIC":                 "test-topic",
		"KAFKA_GROUP_ID":              "test-group",
		"API_ENDPOINT":                "http://api:8000/test",
		"NANOBOT_NAME":                "test-nanobot",
		"HTTP_CLIENT_TIMEOUT_SECONDS": "60",
	}

	// Set environment variables
	for key, value := range testEnvVars {
		os.Setenv(key, value)
		defer os.Unsetenv(key)
	}

	t.Run("config with environment variables", func(t *testing.T) {
		config := Load()

		assert.Equal(t, "kafka:9092", config.KafkaBroker)
		assert.Equal(t, "test-topic", config.KafkaTopic)
		assert.Equal(t, "test-group", config.KafkaGroupID)
		assert.Equal(t, "http://api:8000/test", config.APIEndpoint)
		assert.Equal(t, "test-nanobot", config.NanobotName)
		assert.Equal(t, 60*time.Second, config.HTTPClientTimeout)
	})
}

func TestConfig_PartialEnvironmentVariables(t *testing.T) {
	// Clean up all environment variables first
	envVars := []string{"KAFKA_BROKER", "KAFKA_TOPIC", "KAFKA_GROUP_ID", "API_ENDPOINT", "NANOBOT_NAME", "HTTP_CLIENT_TIMEOUT_SECONDS"}
	for _, env := range envVars {
		os.Unsetenv(env)
	}

	// Set only some environment variables
	os.Setenv("KAFKA_BROKER", "kafka-partial:9092")
	os.Setenv("HTTP_CLIENT_TIMEOUT_SECONDS", "10")
	defer os.Unsetenv("KAFKA_BROKER")
	defer os.Unsetenv("HTTP_CLIENT_TIMEOUT_SECONDS")

	t.Run("config with partial environment variables", func(t *testing.T) {
		config := Load()

		assert.Equal(t, "kafka-partial:9092", config.KafkaBroker)
		assert.Equal(t, "anyker-topic", config.KafkaTopic)
		assert.Equal(t, "anyker-group", config.KafkaGroupID)
		assert.Equal(t, "http://localhost:8080/messages", config.APIEndpoint)
		assert.Equal(t, "anyker-nanobot-1", config.NanobotName)
		assert.Equal(t, 10*time.Second, config.HTTPClientTimeout)
	})
}

func TestSetLogLevel(t *testing.T) {
	tests := []struct {
		name     string
		logLevel string
		expected zerolog.Level
	}{
		{
			name:     "debug level",
			logLevel: "debug",
			expected: zerolog.DebugLevel,
		},
		{
			name:     "info level",
			logLevel: "info",
			expected: zerolog.InfoLevel,
		},
		{
			name:     "warn level",
			logLevel: "warn",
			expected: zerolog.WarnLevel,
		},
		{
			name:     "error level",
			logLevel: "error",
			expected: zerolog.ErrorLevel,
		},
		{
			name:     "fatal level",
			logLevel: "fatal",
			expected: zerolog.FatalLevel,
		},
		{
			name:     "panic level",
			logLevel: "panic",
			expected: zerolog.PanicLevel,
		},
		{
			name:     "uppercase level",
			logLevel: "DEBUG",
			expected: zerolog.DebugLevel,
		},
		{
			name:     "mixed case level",
			logLevel: "WaRn",
			expected: zerolog.WarnLevel,
		},
		{
			name:     "invalid level defaults to info",
			logLevel: "invalid",
			expected: zerolog.InfoLevel,
		},
		{
			name:     "empty level defaults to info",
			logLevel: "",
			expected: zerolog.InfoLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up environment
			os.Unsetenv("LOG_LEVEL")

			// Set LOG_LEVEL if provided
			if tt.logLevel != "" {
				os.Setenv("LOG_LEVEL", tt.logLevel)
				defer os.Unsetenv("LOG_LEVEL")
			}

			// Call setLogLevel
			setLogLevel()

			// Verify the global log level was set correctly
			actual := zerolog.GlobalLevel()
			assert.Equal(t, tt.expected, actual)
		})
	}
}
