package service

import (
	"context"

	"github.com/kekus228swaga/orderflow/order-service/internal/domain/order"
)

type OrderService struct {
	repo order.Repository
}

func NewOrderService(repo order.Repository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) CreateOrder(ctx context.Context, userID int64, req order.CreateOrderRequest) (*order.Order, error) {
	newOrder := &order.Order{
		UserID:      userID,
		ProductName: req.ProductName,
	}

	if err := s.repo.Create(ctx, newOrder); err != nil {
		return nil, err
	}

	return newOrder, nil
}
