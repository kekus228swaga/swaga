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

	"github.com/kekus228swaga/orderflow/auth-service/internal/config"
	"github.com/kekus228swaga/orderflow/auth-service/internal/handler"
	"github.com/kekus228swaga/orderflow/auth-service/internal/repository"
	"github.com/kekus228swaga/orderflow/auth-service/internal/service"
)

func main() {
	// 1. Конфиг
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}

	// 2. БД
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DBDSN)
	if err != nil {
		log.Fatalf("❌ DB error: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("❌ DB ping failed: %v", err)
	}
	log.Println("✅ PostgreSQL connected")

	// 3. Инициализация слоев (Dependency Injection)
	userRepo := repository.NewUserRepo(pool)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService, cfg.JWTSecret)

	// 4. Роутер
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Healthcheck
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Группируем роуты авторизации
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
	}

	// 5. Запуск сервера
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("🚀 Auth service starting on http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf(" Server error: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down...")

	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutCtx); err != nil {
		log.Printf("⚠️ Shutdown error: %v", err)
	}
	pool.Close()
}
