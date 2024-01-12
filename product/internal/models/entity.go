package models

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
