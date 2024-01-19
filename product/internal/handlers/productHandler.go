package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/agusespa/ecom-be-grpc/product/internal/service"
)

type ProductHandler struct {
	ProductService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}

func (h *ProductHandler) HandleAllProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		products, err := h.ProductService.GetAllProducts()
		if err != nil {
			// TODO: handle error better status code
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonBytes, err := json.Marshal(products)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(jsonBytes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	segments := strings.Split(r.URL.Path, "/")
	id := segments[len(segments)-1]

	if r.Method == "GET" {
		product, err := h.ProductService.GetProductById(id)
		if err != nil {
			// TODO: handle error better status code
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		jsonBytes, err := json.Marshal(product)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write(jsonBytes); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
