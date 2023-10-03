package api

import (
	"order/internal/application/domain"
	"order/internal/ports"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{
		db: db,
	}
}

func (a Application) PostOrder(order domain.Order) (domain.Order, error) {
	order, err := a.db.CreateOrder(order)
	if err != nil {
		return domain.Order{}, err
	}
	return order, nil
}

func (a Application) GetOrder(id int64) (domain.Order, error) {
	return a.db.ReadOrder(id)
}
