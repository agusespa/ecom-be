package handlers

import (
	"net/http"

	"github.com/agusespa/ecom-be-grpc/product/internal/service"
)

type ProductHandler struct {
	ProductService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}

func (handler *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	// Implement handler logic to get product by ID from the use case
	// Parse request parameters, call use case method, and write response
}
