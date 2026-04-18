package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kekus228swaga/orderflow/order-service/internal/domain/order"
	"github.com/kekus228swaga/orderflow/order-service/internal/middleware"
	"github.com/kekus228swaga/orderflow/order-service/internal/publisher"
	"github.com/kekus228swaga/orderflow/order-service/internal/service"
)

type OrderHandler struct {
	service   *service.OrderService
	publisher *publisher.Publisher // ← НОВОЕ ПОЛЕ
}

// Обновляем конструктор: теперь принимаем publisher
func NewOrderHandler(s *service.OrderService, p *publisher.Publisher) *OrderHandler {
	return &OrderHandler{service: s, publisher: p}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req order.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID := c.GetInt64(middleware.UserIDKey)

	// 1. Сохраняем в БД
	createdOrder, err := h.service.CreateOrder(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	// 2. Отправляем событие в RabbitMQ АСИНХРОННО
	event := publisher.OrderEvent{
		OrderID:     createdOrder.ID,
		UserID:      userID,
		ProductName: createdOrder.ProductName,
	}
	go func() {
		// Используем фоновый контекст, т.к. HTTP-запрос может завершиться раньше
		if err := h.publisher.Publish(context.Background(), event); err != nil {
			log.Printf("⚠️ RabbitMQ publish error: %v", err)
		}
	}()

	// 3. Сразу отвечаем клиенту
	c.JSON(http.StatusCreated, createdOrder)
}
