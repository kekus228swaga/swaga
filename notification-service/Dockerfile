# --- BUILD STAGE ---
FROM golang:1.25-alpine AS builder
WORKDIR /app

# Кэшируем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем код и собираем воркер
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /notification-worker ./cmd/server

# --- RUN STAGE ---
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /notification-worker /app/notification-worker

# Порт не нужен, так как это консьюмер (слушает очередь, не принимает HTTP)
CMD ["./notification-worker"]