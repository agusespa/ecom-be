package main

import (
	"database/sql"
	"log"
	"net"
	"os"

	pb "product/proto"

	"github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
)

type ProductServiceServer struct {
	pb.UnimplementedProductServiceServer
}

func main() {
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

	defer db.Close()

	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %s", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProductServiceServer(grpcServer, &ProductServiceServer{})

	log.Println("gRPC server started at port 9000")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
