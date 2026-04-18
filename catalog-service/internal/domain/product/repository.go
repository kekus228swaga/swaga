package product

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository interface {
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Product, error)
}
