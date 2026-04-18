FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /catalog-server ./cmd/server

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /catalog-server /app/catalog-server
EXPOSE 8083
CMD ["./catalog-server"]