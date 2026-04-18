package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/kekus228swaga/orderflow/notification-service/internal/config"
	"github.com/kekus228swaga/orderflow/notification-service/internal/worker"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	w, err := worker.New(ctx, cfg.RabbitMQURL, cfg.QueueName)
	if err != nil {
		log.Fatalf("❌ RabbitMQ connection failed: %v", err)
	}
	defer w.Channel.Close()
	defer w.Conn.Close()

	log.Println("✅ Notification service connected to RabbitMQ")

	if err := w.Start(ctx); err != nil {
		log.Fatalf("❌ Worker error: %v", err)
	}
}
