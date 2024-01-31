package handlers

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/agusespa/ecom-be/product/internal/errors"
	"github.com/agusespa/ecom-be/product/internal/payload"
	"github.com/agusespa/ecom-be/product/internal/service"
)

type ProductHandler struct {
	ProductService *service.ProductService
}

func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{ProductService: productService}
}

func (h *ProductHandler) HandleProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		err := errors.NewError(nil, http.StatusMethodNotAllowed)
		payload.WriteError(w, r, err)
		return
	}

	queryParams := r.URL.Query()

	category := queryParams.Get("category")
	decodedCategory, err := url.QueryUnescape(category)
	if err != nil {
		err := errors.NewError(err, http.StatusBadRequest)
		payload.WriteError(w, r, err)
		return
	}

	brand := queryParams.Get("brand")
	decodedBrand, err := url.QueryUnescape(brand)
	if err != nil {
		err := errors.NewError(err, http.StatusBadRequest)
		payload.WriteError(w, r, err)
		return
	}

	products, err := h.ProductService.GetProducts(decodedCategory, decodedBrand)
	if err != nil {
		payload.WriteError(w, r, err)
		return
	}

	payload.Write(w, r, products)
}

func (h *ProductHandler) HandleProductSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		err := errors.NewError(nil, http.StatusMethodNotAllowed)
		payload.WriteError(w, r, err)
		return
	}

	queryParams := r.URL.Query()

	term := queryParams.Get("term")
	decodedTerm, err := url.QueryUnescape(term)
	if err != nil {
		err := errors.NewError(err, http.StatusBadRequest)
		payload.WriteError(w, r, err)
		return
	}

	products, err := h.ProductService.GetProductsBySearchTerm(decodedTerm)
	if err != nil {
		payload.WriteError(w, r, err)
		return
	}

	payload.Write(w, r, products)
}

func (h *ProductHandler) HandleProductByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		err := errors.NewError(nil, http.StatusMethodNotAllowed)
		payload.WriteError(w, r, err)
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

func (h *ProductHandler) HandleProductCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		err := errors.NewError(nil, http.StatusMethodNotAllowed)
		payload.WriteError(w, r, err)
		return
	}

	categories, err := h.ProductService.GetCategories()
	if err != nil {
		payload.WriteError(w, r, err)
		return
	}

	payload.Write(w, r, categories)
}
