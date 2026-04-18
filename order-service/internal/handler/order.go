package handler

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kekus228swaga/orderflow/order-service/internal/domain/order"
	"github.com/kekus228swaga/orderflow/order-service/internal/kafka" // ← Новый импорт
	"github.com/kekus228swaga/orderflow/order-service/internal/middleware"
	"github.com/kekus228swaga/orderflow/order-service/internal/publisher"
	"github.com/kekus228swaga/orderflow/order-service/internal/service"
)

type OrderHandler struct {
	service         *service.OrderService
	rabbitPublisher *publisher.Publisher // RabbitMQ
	kafkaProducer   *kafka.Producer      // Kafka ← Новое поле
}

// Конструктор теперь принимает ОБА паблишера
func NewOrderHandler(
	s *service.OrderService,
	rp *publisher.Publisher,
	kp *kafka.Producer,
) *OrderHandler {
	return &OrderHandler{
		service:         s,
		rabbitPublisher: rp,
		kafkaProducer:   kp,
	}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req order.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	userID := c.GetInt64(middleware.UserIDKey)

	// 1. Сохраняем заказ в БД
	createdOrder, err := h.service.CreateOrder(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	// Готовим событие для отправки
	event := publisher.OrderEvent{
		OrderID:     createdOrder.ID,
		UserID:      userID,
		ProductName: createdOrder.ProductName,
	}

	// 2. Отправляем в RabbitMQ (асинхронно, для надёжной доставки)
	go func() {
		if err := h.rabbitPublisher.Publish(context.Background(), event); err != nil {
			log.Printf("⚠️ RabbitMQ publish error: %v", err)
		}
	}()

	// 3. Отправляем в Kafka (синхронно для простоты, можно тоже в горутину)
	// Kafka используется для стриминга/аналитики, потеря одного сообщения не критична
	h.kafkaProducer.Send(event)

	// 4. Сразу отвечаем клиенту
	c.JSON(http.StatusCreated, createdOrder)
}
