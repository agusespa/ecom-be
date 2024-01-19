package service

import (
	"github.com/agusespa/ecom-be-grpc/product/internal/models"
	"github.com/agusespa/ecom-be-grpc/product/internal/repository"
)

type ProductService struct {
	ProductRepo *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{ProductRepo: productRepo}
}

// TODO pagination
func (ps *ProductService) GetAllProducts() ([]models.Product, error) {
	productEntities, err := ps.ProductRepo.QueryProducts("", "", "")
	var mappedProducts []models.Product
	for _, entity := range productEntities {
		mappedProduct := models.NewSimpleProduct(
			entity.ID,
			entity.Name,
			entity.Brand,
			entity.Subtitle,
			entity.Category,
			entity.Price,
			entity.Currency,
		)
		mappedProducts = append(mappedProducts, mappedProduct)
	}
	return mappedProducts, err
}

func (ps *ProductService) GetProductById(id string) (models.Product, error) {
	entity, err := ps.ProductRepo.QueryProductById(id)
	mappedProduct := models.NewProduct(
		entity.ID,
		entity.Name,
		entity.Brand,
		entity.Subtitle,
		entity.Category,
		entity.Price,
		entity.Currency,
		entity.Stock,
		entity.Description,
		entity.Sku,
	)
	return mappedProduct, err
}
