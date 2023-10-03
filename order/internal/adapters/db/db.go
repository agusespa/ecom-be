package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"order/internal/application/domain"

	"github.com/go-sql-driver/mysql"
)

type Adapter struct {
	db *sql.DB
}

type OrderEntity struct {
	ID         int64  `db:"order_id"`
	CustomerID int64  `db:"customer_id"`
	Status     string `db:"status"`
	CreatedAt  int64  `db:"created_at"`
}

type OrderItemEntity struct {
	ID       int64   `db:"item_id"`
	Sku      string  `db:"sku"`
	Price    float32 `db:"price"`
	Currency string  `db:"currency"`
	Quantity int32   `db:"quantity"`
}

func NewAdapter(dataSourceUrl string) (*Adapter, error) {
	cfg := mysql.Config{
		User:   os.Getenv("DBUSER"),
		Passwd: os.Getenv("DBPASS"),
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		DBName: "databasename",
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatalf("Error opening database connection: %v", err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatalf("Error pinging database: %v", pingErr)
	}

	return &Adapter{db: db}, nil
}

func (a Adapter) Create(newOrder domain.Order) (domain.Order, error) {
	orderResult, err := a.db.Exec("INSERT INTO orders (customer_id, status) VALUES (?, ?)", newOrder.CustomerID, newOrder.Status)
	if err != nil {
		return domain.Order{}, fmt.Errorf("save_order: %v", err)
	}

	for _, value := range newOrder.OrderProducts {
		_, err := a.db.Exec("INSERT INTO items (sku, price, quantity) VALUES (?, ?, ?)", value.ID, value.Price, value.Quantity)
		if err != nil {
			return domain.Order{}, fmt.Errorf("save_item: %v", err)
		}
	}

	id, err := orderResult.LastInsertId()
	if err != nil {
		return domain.Order{}, fmt.Errorf("save_order: %v", err)
	}

	order, err := a.Read(id)
	if err != nil {
		return domain.Order{}, err
	}
	return order, nil
}

func (a Adapter) Read(id int64) (domain.Order, error) {
	var orderEntity OrderEntity
	row := a.db.QueryRow("SELECT * FROM orders WHERE order_id = ?", id)
	if err := row.Scan(&orderEntity.ID, &orderEntity.CustomerID, &orderEntity.Status, &orderEntity.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return domain.Order{}, fmt.Errorf("no order with id %d", id)
		}
		return domain.Order{}, fmt.Errorf("%d: %v", id, err)
	}

	var orderProducts []domain.OrderProduct
	rows, err := a.db.Query("SELECT item_id, sku, price, quantity FROM items WHERE order_id = ?", id)
	if err != nil {
		return domain.Order{}, fmt.Errorf("items by order %d: %v", id, err)
	}
	defer rows.Close()

	for rows.Next() {
		var orderItem domain.OrderProduct

		if err := rows.Scan(&orderItem.ID, &orderItem.Price, &orderItem.Quantity); err != nil {
			return domain.Order{}, fmt.Errorf("items by order %d: %v", id, err)
		}
		orderProducts = append(orderProducts, orderItem)
	}

	order := domain.Order{
		ID:            int64(orderEntity.ID),
		CustomerID:    orderEntity.CustomerID,
		Status:        orderEntity.Status,
		OrderProducts: orderProducts,
		CreatedAt:     orderEntity.CreatedAt,
	}
	return order, nil
}
