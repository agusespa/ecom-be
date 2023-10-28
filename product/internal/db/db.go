package db

import (
	"fmt"
)

type ProductEntity struct {
	ID          int64   `db:"item_id"`
	Name        string  `db:"name"`
	Brand       string  `db:"brand"`
	Subtitle    string  `db:"subtitle"`
	Category    string  `db:"category"`
	Gender      string  `db:"gender"`
	Color       string  `db:"color"`
	Price       float32 `db:"price"`
	Currency    string  `db:"currency"`
	Sku         string  `db:"sku"`
	Description string  `db:"description"`
	Stock       int32   `db:"quantity"`
}

func QueryProducts(category string, gender string, name string) ([]ProductEntity, error) {
	var products []ProductEntity

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

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("error")
	}
	defer rows.Close()

	for rows.Next() {
		var product ProductEntity
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
			return nil, fmt.Errorf("error")
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error")
	}
	return products, nil
}
