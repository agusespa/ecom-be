package models

type ProductEntity struct {
	ID          int64   `db:"product_id"`
	Name        string  `db:"name"`
	Subtitle    string  `db:"subtitle"`
	Category    string  `db:"category"`
	Brand       string  `db:"brand"`
	Price       float32 `db:"price"`
	Currency    string  `db:"currency"`
	Stock       int32   `db:"quantity"`
	Description string  `db:"description"`
	Sku         string  `db:"sku"`
}
