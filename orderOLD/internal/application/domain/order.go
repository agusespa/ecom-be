package domain

type Order struct {
	ID            int64          `json:"id"`
	CustomerID    int64          `json:"customer_id"`
	Status        string         `json:"status"`
	OrderProducts []OrderProduct `json:"order_items"`
	CreatedAt     int64          `json:"created_at"`
}

type OrderProduct struct {
	ID       int64  `json:"id"`
	Sku      string `json:"sku"`
	Price    Price  `json:"price"`
	Quantity int32  `json:"quantity"`
}

type Price struct {
	Amount   float32 `json:"amount"`
	Currency string  `json:"currency"`
}
