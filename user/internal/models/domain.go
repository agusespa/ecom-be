package models

type Product struct {
	ID       int64          `json:"id"`
	Name     string         `json:"name"`
	Brand    string         `json:"brand"`
	Subtitle string         `json:"subtitle"`
	Category string         `json:"category"`
	Price    SimplePrice    `json:"price"`
	Details  ProductDetails `json:"details"`
}

type SimplePrice struct {
	Amount   float32 `json:"amount"`
	Currency string  `json:"currency"`
}

type ProductDetails struct {
	Stock       int32  `json:"stock"`
	Description string `json:"description"`
	Sku         string `json:"sku"`
}

func NewSimpleProduct(id int64, name string, brand string, subtitle string, category string, priceAmount float32, currency string) Product {
	price := newSimplePrice(priceAmount, currency)
	return Product{
		ID:       id,
		Name:     name,
		Brand:    brand,
		Subtitle: subtitle,
		Category: category,
		Price:    price,
	}
}

func NewProduct(id int64, name string, brand string, subtitle string, category string, priceAmount float32, currency string, stock int32, description string, sku string) Product {
	price := newSimplePrice(priceAmount, currency)
	details := newProductDetails(stock, description, sku)
	return Product{
		ID:       id,
		Name:     name,
		Brand:    brand,
		Subtitle: subtitle,
		Category: category,
		Price:    price,
		Details:  details,
	}
}

func newSimplePrice(priceAmount float32, currency string) SimplePrice {
	return SimplePrice{
		Amount:   priceAmount,
		Currency: currency,
	}
}

func newProductDetails(stock int32, description string, sku string) ProductDetails {
	return ProductDetails{
		Stock:       stock,
		Description: description,
		Sku:         sku,
	}
}
