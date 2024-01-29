package ports

import "order/internal/application/domain"

type DBPort interface {
	ReadOrder(int64) (domain.Order, error)
	CreateOrder(domain.Order) (domain.Order, error)
}
