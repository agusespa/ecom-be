package main

import (
	"context"
	"fmt"
	"log"

	pb "bookstore/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewBookServiceClient(conn)

	req := &pb.BookRequest{
		Id: 248,
	}

	res, err := client.GetBook(context.Background(), req)
	if err != nil {
		log.Fatalf("Add RPC failed: %v", err)
	}

	fmt.Printf("Result: %s\n", res.Title)
}
