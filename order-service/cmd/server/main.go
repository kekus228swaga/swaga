package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kekus228swaga/orderflow/order-service/internal/handler"
	"github.com/kekus228swaga/orderflow/order-service/internal/kafka"
	"github.com/kekus228swaga/orderflow/order-service/internal/publisher"
	"github.com/kekus228swaga/orderflow/order-service/internal/repository"
	"github.com/kekus228swaga/orderflow/order-service/internal/service"

	"github.com/kekus228swaga/orderflow/order-service/internal/config"
	"github.com/kekus228swaga/orderflow/order-service/internal/middleware"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DBDSN)
	if err != nil {
		log.Fatalf("❌ DB error: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("❌ DB ping failed: %v", err)
	}
	log.Println("✅ Order service connected to PostgreSQL")

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Глобальный middleware (для всех роутов)
	// Пока оставим без него, чтобы healthcheck работал публично
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "order", "status": "ok"})
	})

	// Инициализация RabbitMQ Publisher
	pub, err := publisher.NewPublisher("amqp://guest:guest@rabbitmq:5672/", "order.created")
	if err != nil {
		log.Fatalf("❌ RabbitMQ publisher init failed: %v", err)
	}
	defer pub.Channel.Close()

	kafkaBrokers := []string{"kafka:9092"} // Внутри Docker сети
	kafkaTopic := "order.events"

	kafkaProducer, err := kafka.NewProducer(kafkaBrokers, kafkaTopic)
	if err != nil {
		log.Fatalf("❌ Kafka producer init failed: %v", err)
	}

	// Инициализация слоев
	orderRepo := repository.NewOrderRepo(pool)
	orderService := service.NewOrderService(orderRepo)

	// Передаём publisher в хендлер
	orderHandler := handler.NewOrderHandler(orderService, pub, kafkaProducer)

	// Роуты
	protected := r.Group("/orders")
	protected.Use(middleware.JWTAuth(cfg.JWTSecret))
	{
		protected.POST("", orderHandler.Create)
		protected.GET("/me", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok", "user_id": c.GetInt64(middleware.UserIDKey)})
		})
	}

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("🚀 Order service starting on http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(" Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println(" Order service shutting down...")

	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutCtx); err != nil {
		log.Printf("⚠️ Shutdown error: %v", err)
	}
	pool.Close()
}
