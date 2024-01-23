package repository

import (
	"database/sql"
	"net/http"

	"github.com/agusespa/ecom-be-grpc/product/internal/errors"
	"github.com/agusespa/ecom-be-grpc/product/internal/models"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (repo *ProductRepository) QueryProductById(id string) (models.ProductEntity, error) {
	row := repo.DB.QueryRow("SELECT * FROM product WHERE product_id = ?", id)
	var product models.ProductEntity

	if err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Subtitle,
		&product.Category,
		&product.Brand,
		&product.Price,
		&product.Currency,
		&product.Stock,
		&product.Description,
		&product.Sku,
	); err != nil {
		if err == sql.ErrNoRows {
			error := errors.New(err, http.StatusNotFound)
			return product, error
		}
		error := errors.New(err, http.StatusInternalServerError)
		return product, error
	}
	return product, nil
}

func (repo *ProductRepository) QueryProducts(category string, name string, brand string) ([]models.ProductEntity, error) {
	var products []models.ProductEntity

	query := "SELECT product_id, name, subtitle, category, brand, price, currency FROM product WHERE 1=1"

	var args []interface{}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}
	if name != "" {
		query += " AND name = ?"
		args = append(args, name)
	}
	if name != "" {
		query += " AND brand = ?"
		args = append(args, brand)
	}

	rows, err := repo.DB.Query(query, args...)
	if err != nil {
		error := errors.New(err, http.StatusInternalServerError)
		return nil, error
	}
	defer rows.Close()

	for rows.Next() {
		var product models.ProductEntity
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Subtitle,
			&product.Category,
			&product.Brand,
			&product.Price,
			&product.Currency,
		); err != nil {
			if err == sql.ErrNoRows {
				error := errors.New(err, http.StatusNotFound)
				return nil, error
			}
			error := errors.New(err, http.StatusInternalServerError)
			return nil, error
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		error := errors.New(err, http.StatusInternalServerError)
		return nil, error
	}
	return products, nil
}
