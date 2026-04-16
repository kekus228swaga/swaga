package order

import "context"

type Repository interface {
	Create(ctx context.Context, order *Order) error
}
