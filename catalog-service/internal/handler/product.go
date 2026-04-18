package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kekus228swaga/orderflow/catalog-service/internal/domain/product"
	"github.com/kekus228swaga/orderflow/catalog-service/internal/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductHandler struct {
	service *service.ProductService
}

func NewProductHandler(s *service.ProductService) *ProductHandler {
	return &ProductHandler{service: s}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req product.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	created, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product"})
		return
	}

	c.JSON(http.StatusCreated, created)
}

func (h *ProductHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	objID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product id"})
		return
	}

	product, err := h.service.GetByID(c.Request.Context(), objID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}
