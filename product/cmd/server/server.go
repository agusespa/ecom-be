package main

import (
	"context"
	"log"
	"net"

	pb "bookstore/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedBookServiceServer
}

func (s *server) GetBook(ctx context.Context, req *pb.BookRequest) (*pb.Book, error) {
	return &pb.Book{Id: 248, Title: "El Aleph", Isbn: 987534831}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterBookServiceServer(grpcServer, &server{})

	log.Println("gRPC server started at port 9000")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
