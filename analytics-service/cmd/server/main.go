package main

import (
	"context"
	"log"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kekus228swaga/orderflow/analytics-service/internal/config"
	"github.com/kekus228swaga/orderflow/analytics-service/internal/consumer"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	brokers := strings.Split(cfg.KafkaBrokers, ",")
	analytics, err := consumer.New(brokers, cfg.GroupID, cfg.Topic)
	if err != nil {
		log.Fatalf("❌ Kafka consumer init failed: %v", err)
	}

	log.Println("✅ Analytics service started, listening to Kafka...")
	analytics.Start(ctx)
}
