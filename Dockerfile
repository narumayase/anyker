# anyker/Dockerfile
FROM golang:1.24-bullseye AS builder

RUN apt-get update && apt-get install -y \
    git \
    pkg-config \
    librdkafka-dev \
    build-essential \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o anyker .

FROM debian:bullseye-slim
WORKDIR /app
COPY --from=builder /app/anyker .

RUN apt-get update && apt-get install -y librdkafka1 && rm -rf /var/lib/apt/lists/*

CMD ["./anyker"]