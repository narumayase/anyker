# Anyker

This project provides a worker/nanobot that consumes messages from a Kafka topic and forwards them to a configured API endpoint.

### FEATURES

*   Consumes messages from a Kafka topic.
*   Forwards messages to a configured API endpoint.
*   Scalable and extensible.

### PREREQUISITES

*   Go 1.21 or higher
*   Kafka Broker
*   Docker (optional)

### 🚀 INSTALLATION

1.  Install dependencies:
    ```sh
    go mod tidy
    ```
2.  Configure environment variables:
    ```sh
    cp env.example .env
    # Edit .env with the values described below.
    ```
3.  Run the application:
    ```sh
    go run main.go
    ```

### 🔧 CONFIGURATION

#### ENVIRONMENT VARIABLES

Create a `.env` file based on `env.example`:

*   `KAFKA_BROKER`: Kafka broker address.
*   `KAFKA_TOPIC`: Kafka topic to consume messages from.
*   `KAFKA_GROUP_ID`: Kafka consumer group ID.
*   `API_ENDPOINT`: API endpoint to forward messages to.
*   `NANOBOT_NAME`: Name of the nanobot instance.
*   `LOG_LEVEL`: Log level (`debug`, `info`, `warn`, `error`, `fatal`, `panic` - default: `info`)

### 📡 ENDPOINTS

This project does not expose any endpoints. It consumes messages from a Kafka topic and forwards them to a configured API endpoint.

### 🎗️ ARCHITECTURE

This project follows Clean Architecture principles:

*   **Domain**: Entities, repository interfaces, and use cases
*   **Application**: Implementation of use cases
*   **Infrastructure**: Kafka consumer and HTTP client repository implementations
*   **Interfaces**: CLI commands and handlers

### 📁 PROJECT STRUCTURE

```
anyker/
├── cmd/                  # Application entry points
│   └── root.go           # Main command
├── internal/             # Project-specific code
│   ├── application/      # Use cases
│   ├── config/           # Configuration
│   ├── domain/           # Domain entities and interfaces
│   └── infrastructure/   # Repository implementations
│       ├── client/       # HTTP client
│       └── repository/   # Kafka consumer
├── main.go               # Main entry point
├── go.mod                # Go dependencies
├── README_es.md          # README in spanish
└── README.md             # This file
```

### 🧪 TESTING

#### RUNNING TESTS

To run all tests:

```sh
go test ./...
```

To run tests with verbose output:

```sh
go test -v ./...
```

To run tests for a specific package:

```sh
go test ./internal/config/
```

#### TEST COVERAGE

To check test coverage (excluding mocks):

```sh
# Generate coverage report
go test -coverprofile=coverage.out ./...

# View coverage report in terminal
go tool cover -func=coverage.out

# Generate HTML coverage report
go tool cover -html=coverage.out -o coverage.html

# View coverage excluding mocks
go test -coverprofile=coverage.out ./... && \
go tool cover -func=coverage.out | grep -v "mocks"
```

### BACKLOG

*   Unit Tests
*   Integration tests
*   Add more message brokers (e.g., RabbitMQ, NATS)
*   API documentation with Swagger

### ABOUT

A worker/nanobot that consumes messages from a Kafka topic and forwards them to a configured API endpoint.
