package order

import "time"

type Order struct {
	ID          int64     `json:"id"`
	UserID      int64     `json:"user_id"`
	ProductName string    `json:"product_name"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateOrderRequest struct {
	ProductName string `json:"product_name" binding:"required"`
}
