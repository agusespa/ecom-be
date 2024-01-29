package grpc

import (
	"context"
	pb "order/proto"

	"order/internal/application/domain"
)

func (a Adapter) Create(ctx context.Context, request *pb.CreateOrderRequest) (pb.OrderResponse, error) {
	var orderProducts []domain.OrderProduct
	for _, orderProduct := range request.OrderProducts {
		orderProducts = append(orderProducts, domain.OrderProduct{
			ID:       orderProduct.ProductId,
			Quantity: orderProduct.Quantity,
		})
	}

	newOrder := domain.NewOrder(request.CustomerId, orderProducts)

	result, err := a.api.PlaceOrder(newOrder)
	if err != nil {
		return pb.OrderResponse{}, err
	}

	var orderProductsResponse []*pb.OrderProductResponse
	for _, orderProduct := range result.OrderProducts {
		productPrice := &pb.Price{
			Amount:   orderProduct.Price.Amount,
			Currency: orderProduct.Price.Currency,
		}
		orderProductsResponse = append(orderProductsResponse, &pb.OrderProductResponse{
			ProductId: orderProduct.ID,
			Sku:       orderProduct.Sku,
			Price:     productPrice,
			Quantity:  orderProduct.Quantity,
		})
	}

	createdOrder := pb.OrderResponse{
		OrderId:       result.ID,
		CustomerId:    result.CustomerID,
		OrderProducts: orderProductsResponse,
		Status:        result.Status,
		CreatedAt:     result.CreatedAt,
	}

	return createdOrder, nil
}
