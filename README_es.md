# Anyker

Este proyecto proporciona un worker/nanobot que consume mensajes de un tópico de Kafka y los reenvía a un endpoint de API configurado.

### CARACTERÍSTICAS

*   Consume mensajes de un tópico de Kafka.
*   Reenvía mensajes a un endpoint de API configurado.
*   Escalable y extensible.

### PREREQUISITOS

*   Go 1.21 o superior
*   Kafka Broker
*   Docker (opcional)

### 🚀 INSTALACIÓN

1.  Instalar dependencias:
    ```sh
    go mod tidy
    ```
2.  Configurar variables de entorno:
    ```sh
    cp env.example .env
    # Edita .env con los valores que se describen a continuación.
    ```
3.  Ejecutar la aplicación:
    ```sh
    go run main.go
    ```

### 🔧 CONFIGURACIÓN

#### VARIABLES DE ENTORNO

Crea un archivo `.env` basado en `env.example`:

*   `KAFKA_BROKER`: Dirección del broker de Kafka.
*   `KAFKA_TOPIC`: Tópico de Kafka del que consumir los mensajes.
*   `KAFKA_GROUP_ID`: ID del grupo de consumidores de Kafka.
*   `API_ENDPOINT`: Endpoint de la API a la que reenviar los mensajes.
*   `NANOBOT_NAME`: Nombre de la instancia del nanobot.
*   `LOG_LEVEL`: Nivel de log (`debug`, `info`, `warn`, `error`, `fatal`, `panic` - por defecto: `info`)

### 📡 ENDPOINTS

Este proyecto no expone ningún endpoint. Consume mensajes de un tópico de Kafka y los reenvía a un endpoint de API configurado.

### 🎗️ ARQUITECTURA

Este proyecto sigue los principios de Clean Architecture:

*   **Domain**: Entidades, interfaces de repositorio y casos de uso
*   **Application**: Implementación de los casos de uso
*   **Infrastructure**: Implementaciones de los repositorios del consumidor de Kafka y del cliente HTTP
*   **Interfaces**: Comandos y manejadores de CLI

### 📁 ESTRUCTURA DEL PROYECTO

```
anyker/
├── cmd/                  # Puntos de entrada de la aplicación
│   └── root.go           # Comando principal
├── internal/             # Código específico del proyecto
│   ├── application/      # Casos de uso
│   ├── config/           # Configuración
│   ├── domain/           # Entidades e interfaces de dominio
│   └── infrastructure/   # Implementaciones de repositorios
│       ├── client/       # Cliente HTTP
│       └── repository/   # Consumidor de Kafka
├── main.go               # Punto de entrada principal
├── go.mod                # Dependencias de Go
├── README_es.md          # README en español
└── README.md             # Este archivo
```

### 🧪 PRUEBAS

#### EJECUTAR PRUEBAS

Para ejecutar todas las pruebas:

```sh
go test ./...
```

Para ejecutar las pruebas con salida detallada:

```sh
go test -v ./...
```

Para ejecutar las pruebas de un paquete específico:

```sh
go test ./internal/config/
```

#### COBERTURA DE PRUEBAS

Para comprobar la cobertura de las pruebas (excluyendo los mocks):

```sh
# Generar informe de cobertura
go test -coverprofile=coverage.out ./...

# Ver el informe de cobertura en la terminal
go tool cover -func=coverage.out

# Generar informe de cobertura HTML
go tool cover -html=coverage.out -o coverage.html

# Ver la cobertura excluyendo los mocks
go test -coverprofile=coverage.out ./... && \
go tool cover -func=coverage.out | grep -v "mocks"
```

### BACKLOG

*   Pruebas unitarias
*   Pruebas de integración
*   Añadir más brokers de mensajería (por ejemplo, RabbitMQ, NATS)
*   Documentación de la API con Swagger

### ACERCA DE

Un worker/nanobot que consume mensajes de un tópico de Kafka y los reenvía a un endpoint de API configurado.
