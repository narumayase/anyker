# Etapa 1: build
FROM golang:1.22-alpine AS builder
RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o anyker .

# Etapa 2: imagen final
FROM alpine:latest
WORKDIR /app

COPY --from=builder /app/anyker .

EXPOSE 8082
CMD ["./anyker"]