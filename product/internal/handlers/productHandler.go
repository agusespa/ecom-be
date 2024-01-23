package handlers

import (
	"net/http"
	"strings"

	"github.com/agusespa/ecom-be-grpc/product/internal/payload"
	"github.com/agusespa/ecom-be-grpc/product/internal/service"
)

type ProductHandler struct {
	ProductService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}

func (h *ProductHandler) HandleAllProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	products, err := h.ProductService.GetAllProducts()
	if err != nil {
		payload.WriteError(w, r, err)
		return
	}

	payload.Write(w, r, products)
}

func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	segments := strings.Split(r.URL.Path, "/")
	id := segments[len(segments)-1]

	product, err := h.ProductService.GetProductById(id)
	if err != nil {
		payload.WriteError(w, r, err)
		return
	}

	payload.Write(w, r, product)
}
