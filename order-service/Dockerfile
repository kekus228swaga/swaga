FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /order-server ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /order-server /app/order-server
EXPOSE 8082
CMD ["./order-server"]