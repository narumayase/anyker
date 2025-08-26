package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
)

// Config holds the application configuration.
type Config struct {
	KafkaBroker  string
	KafkaTopic   string
	KafkaGroupID string
	APIEndpoint  string
	NanobotName  string
	LogLevel     string
}

// Load loads configuration from environment variables or an .env file
func Load() Config {
	// Load .env file if it exists (ignore error if file doesn't exist)
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading .env file: %v", err)
	}
	setLogLevel()

	return Config{
		KafkaBroker:  getEnv("KAFKA_BROKER", "localhost:9092"),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "anyker-topic"),
		KafkaGroupID: getEnv("KAFKA_GROUP_ID", "anyker-group"),
		APIEndpoint:  getEnv("API_ENDPOINT", "http://localhost:8080/messages"),
		NanobotName:  getEnv("NANOBOT_NAME", "anyker-nanobot-1"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// setLogLevel sets the log level defined in LOG_LEVEL environment variable
func setLogLevel() {
	levels := map[string]zerolog.Level{
		"debug": zerolog.DebugLevel,
		"info":  zerolog.InfoLevel,
		"warn":  zerolog.WarnLevel,
		"error": zerolog.ErrorLevel,
		"fatal": zerolog.FatalLevel,
		"panic": zerolog.PanicLevel,
	}
	levelEnv := strings.ToLower(getEnv("LOG_LEVEL", "info"))

	level, ok := levels[levelEnv]
	if !ok {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)
}
