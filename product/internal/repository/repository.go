package repository

import (
	"database/sql"
	"fmt"

	"github.com/agusespa/ecom-be-grpc/product/internal/models"
)

type ProductRepository struct {
	DB *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{DB: db}
}

func (repo *ProductRepository) QueryProductById(id int32) (models.ProductEntity, error) {
	row := repo.DB.QueryRow("SELECT * FROM product WHERE id = ?", id)
	product, err := scanProduct(row)
	if err != nil {
		return product, err
	}
	return product, nil
}

func scanProduct(row *sql.Row) (models.ProductEntity, error) {
	var product models.ProductEntity
	if err := row.Scan(
		&product.ID,
		&product.Name,
		&product.Brand,
		&product.Subtitle,
		&product.Category,
		&product.Gender,
		&product.Color,
		&product.Price,
		&product.Currency,
	); err != nil {
		if err == sql.ErrNoRows {
			return product, fmt.Errorf("product not found")
		}
		return product, fmt.Errorf("error")
	}
	return product, nil
}

func (repo *ProductRepository) QueryProducts(category string, gender string, name string) ([]models.ProductEntity, error) {
	var products []models.ProductEntity

	query := "SELECT * FROM product WHERE 1=1"

	var args []interface{}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}
	if gender != "" {
		query += " AND gender = ?"
		args = append(args, gender)
	}
	if name != "" {
		query += " AND name = ?"
		args = append(args, name)
	}

	rows, err := repo.DB.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error")
	}
	defer rows.Close()

	for rows.Next() {
		var product models.ProductEntity
		if err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Brand,
			&product.Subtitle,
			&product.Category,
			&product.Gender,
			&product.Color,
			&product.Price,
			&product.Currency,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("product not found")
			}
			return nil, fmt.Errorf("error")
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error")
	}
	return products, nil
}
