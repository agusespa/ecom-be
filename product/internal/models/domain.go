package models

type Product struct {
	ID       int64
	Name     string
	Brand    string
	Subtitle string
	Category string
	Gender   string
	Color    string
	Price    Price
}

type Price struct {
	Amount   float64
	Currency string
}

type ProductDetails struct {
	Sku         string
	Description string
	Stock       int32
}
