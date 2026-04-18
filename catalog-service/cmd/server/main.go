package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/kekus228swaga/orderflow/catalog-service/internal/config"
	"github.com/kekus228swaga/orderflow/catalog-service/internal/handler"
	"github.com/kekus228swaga/orderflow/catalog-service/internal/repository"
	"github.com/kekus228swaga/orderflow/catalog-service/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("❌ Config error: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// MongoDB подключение
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("❌ MongoDB connection failed: %v", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ MongoDB ping failed: %v", err)
	}
	log.Println("✅ MongoDB connected")

	db := client.Database(cfg.DBName)
	productRepo := repository.NewProductRepo(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"service": "catalog", "status": "ok"})
	})

	api := r.Group("/products")
	{
		api.POST("", productHandler.Create)
		api.GET("/:id", productHandler.GetByID)
	}

	server := &http.Server{Addr: ":" + cfg.Port, Handler: r}

	go func() {
		log.Printf("🚀 Catalog service starting on http://localhost:%s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server error: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("🛑 Catalog service shutting down...")

	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	server.Shutdown(shutCtx)
}
