package repository

import (
	"context"

	"github.com/kekus228swaga/orderflow/catalog-service/internal/domain/product"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ProductRepo struct {
	collection *mongo.Collection
}

func NewProductRepo(db *mongo.Database) *ProductRepo {
	return &ProductRepo{collection: db.Collection("products")}
}

func (r *ProductRepo) Create(ctx context.Context, p *product.Product) error {
	result, err := r.collection.InsertOne(ctx, p)
	if err != nil {
		return err
	}
	p.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ProductRepo) GetByID(ctx context.Context, id primitive.ObjectID) (*product.Product, error) {
	var p product.Product
	// ✅ Используем bson.M вместо структуры
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
