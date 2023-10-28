package rpc

import (
	"context"

	server "product/cmd/main"
	"product/internal/db"
	pb "product/proto"
)

func (s *server.ProductServiceServer) SearchProducts(ctx context.Context, req *pb.ProductSearchRequest) (*pb.ProductsResponse, error) {
	category := req.GetCategory()
	gender := req.GetGender()
	name := req.GetName()

	products, err := db.QueryProducts(category, gender, name)
	if err != nil {
		return nil, err
	}

	productsResponse := make([]*pb.ProductsResponse, 0)
	for _, product := range products {
		productsResponse = append(products, &pb.ProductResponse{
			Id:       product.ID,
			Name:     product.Name,
			Category: product.Category,
			Color:    product.Color,
			// TODO Add other fields as needed
		})
	}

	return &pb.ProductResponse{Product: productsResponse}, nil
}
