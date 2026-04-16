package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kekus228swaga/orderflow/order-service/internal/domain/order"
)

type OrderRepo struct {
	pool *pgxpool.Pool
}

func NewOrderRepo(pool *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{pool: pool}
}

func (r *OrderRepo) Create(ctx context.Context, ord *order.Order) error {
	// Статус по умолчанию "new"
	ord.Status = "new"

	query := `
		INSERT INTO orders (user_id, product_name, status) 
		VALUES ($1, $2, $3) 
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(ctx, query, ord.UserID, ord.ProductName, ord.Status).
		Scan(&ord.ID, &ord.CreatedAt)

	if err != nil {
		return err
	}
	return nil
}
