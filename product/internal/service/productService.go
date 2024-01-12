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

func (ps *ProductService) GetProductById(id int32) (models.Product, error) {
	productEntity, err := ps.ProductRepo.QueryProductById(id)
	convertedProduct := convertProductEntityToProduct(productEntity)
	return convertedProduct, err
}

func (ps *ProductService) GetAllProducts() ([]models.Product, error) {
	productEntities, err := ps.ProductRepo.QueryProducts("", "", "")
	var convertedProducts []models.Product
	for _, entity := range productEntities {
		convertedProducts = append(convertedProducts, convertProductEntityToProduct(entity))
	}
	return convertedProducts, err
}

func convertProductEntityToProduct(entity models.ProductEntity) models.Product {
	price := models.Price{
		Amount:   float64(entity.Price),
		Currency: entity.Currency,
	}
	product := models.Product{
		ID:       entity.ID,
		Name:     entity.Name,
		Brand:    entity.Brand,
		Subtitle: entity.Subtitle,
		Category: entity.Category,
		Gender:   entity.Gender,
		Color:    entity.Color,
		Price:    price,
	}
	return product
}
