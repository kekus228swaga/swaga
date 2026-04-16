# Сборка
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /auth-server ./cmd/server

# Рантайм
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /auth-server /app/auth-server
EXPOSE 8081
CMD ["./auth-server"]