package repository

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/agusespa/ecom-be/product/internal/errors"
	"github.com/agusespa/ecom-be/product/internal/models"
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
			error := errors.NewError(err, http.StatusNotFound)
			return product, error
		}
		error := errors.NewError(err, http.StatusInternalServerError)
		return product, error
	}
	return product, nil
}

func (repo *ProductRepository) QueryCategories() ([]string, error) {
	rows, err := repo.DB.Query("SELECT DISTINCT category FROM product")
	if err != nil {
		error := errors.NewError(err, http.StatusInternalServerError)
		return nil, error
	}
	defer rows.Close()

	var categories []string
	for rows.Next() {
		var category string
		if err := rows.Scan(&category); err != nil {
			error := errors.NewError(err, http.StatusInternalServerError)
			return nil, error
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		error := errors.NewError(err, http.StatusInternalServerError)
		return nil, error
	}

	return categories, nil
}

func (repo *ProductRepository) QueryProducts(category string, brand string) ([]models.ProductEntity, error) {
	var products []models.ProductEntity

	query := "SELECT product_id, name, subtitle, category, brand, price, currency FROM product WHERE 1=1"

	var args []interface{}

	if category != "" {
		query += " AND category = ?"
		args = append(args, category)
	}
	if brand != "" {
		query += " AND brand = ?"
		args = append(args, brand)
	}

	rows, err := repo.DB.Query(query, args...)
	if err != nil {
		error := errors.NewError(err, http.StatusInternalServerError)
		return products, error
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
			if err != sql.ErrNoRows {
				error := errors.NewError(err, http.StatusInternalServerError)
				return nil, error
			}
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		error := errors.NewError(err, http.StatusInternalServerError)
		return nil, error
	}

	return products, nil
}

func (repo *ProductRepository) SearchProducts(term string) ([]models.ProductEntity, error) {
	var products []models.ProductEntity

	query := "SELECT product_id, name, subtitle, category, brand, price, currency FROM product WHERE 1=0"

	var args []interface{}

	conditions := []string{"category", "brand", "name"}

	for _, field := range conditions {
		query += fmt.Sprintf(" OR %s LIKE ?", field)
		args = append(args, "%"+term+"%")
	}

	rows, err := repo.DB.Query(query, args...)
	if err != nil {
		error := errors.NewError(err, http.StatusInternalServerError)
		return products, error
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
			if err != sql.ErrNoRows {
				error := errors.NewError(err, http.StatusInternalServerError)
				return nil, error
			}
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		error := errors.NewError(err, http.StatusInternalServerError)
		return nil, error
	}

	return products, nil
}
