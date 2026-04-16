package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kekus228swaga/orderflow/order-service/internal/domain/order"
	"github.com/kekus228swaga/orderflow/order-service/internal/middleware"
	"github.com/kekus228swaga/orderflow/order-service/internal/service"
)

type OrderHandler struct {
	service *service.OrderService
}

func NewOrderHandler(service *service.OrderService) *OrderHandler {
	return &OrderHandler{service: service}
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req order.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	//  Достаем user_id, который положил туда Middleware
	userID := c.GetInt64(middleware.UserIDKey)

	createdOrder, err := h.service.CreateOrder(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, createdOrder)
}
