package ports

import "order/internal/application/domain"

type APIPort interface {
	PostOrder(domain.Order) (domain.Order, error)
	GetOrder(int64) (domain.Order, error)
}
