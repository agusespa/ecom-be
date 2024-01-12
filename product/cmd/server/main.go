package main

import (
	"log"
	"net/http"

	"github.com/agusespa/ecom-be-grpc/product/internal/database"
	"github.com/agusespa/ecom-be-grpc/product/internal/handlers"
	"github.com/agusespa/ecom-be-grpc/product/internal/repository"
	"github.com/agusespa/ecom-be-grpc/product/internal/service"
)

func main() {
	db, dbErr := database.ConnectDB()
	if dbErr != nil {
		log.Fatalf("Error establishing database connection: %v", dbErr)
	}

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	http.HandleFunc("/products/{id}", productHandler.HandleProductByID)
	// Define other routes for CRUD operations

	log.Println("Starting the HTTP server")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting the HTTP server: %v", err)
	}
}
