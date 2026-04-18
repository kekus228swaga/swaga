package service

import (
	"context"
	"time"

	"github.com/kekus228swaga/orderflow/catalog-service/internal/domain/product"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductService struct {
	repo product.Repository
}

func NewProductService(repo product.Repository) *ProductService {
	return &ProductService{repo: repo}
}

func (s *ProductService) Create(ctx context.Context, req product.CreateProductRequest) (*product.Product, error) {
	p := &product.Product{
		Name:       req.Name,
		Price:      req.Price,
		Attributes: req.Attributes,
		CreatedAt:  primitive.NewDateTimeFromTime(time.Now()),
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *ProductService) GetByID(ctx context.Context, id primitive.ObjectID) (*product.Product, error) {
	return s.repo.GetByID(ctx, id)
}
