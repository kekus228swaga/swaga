package product

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name       string             `json:"name" bson:"name"`
	Price      float64            `json:"price" bson:"price"`
	Attributes map[string]any     `json:"attributes" bson:"attributes,omitempty"` // Гибкие поля
	CreatedAt  primitive.DateTime `json:"created_at" bson:"created_at,omitempty"`
}

type CreateProductRequest struct {
	Name       string         `json:"name" binding:"required"`
	Price      float64        `json:"price" binding:"required,gt=0"`
	Attributes map[string]any `json:"attributes"`
}
