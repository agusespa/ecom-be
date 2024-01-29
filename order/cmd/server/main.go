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

	// TODO: get port dinamically
	port := "8080"

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	http.HandleFunc("/products", productHandler.HandleProducts)
	http.HandleFunc("/products/", productHandler.HandleProductByID)
	http.HandleFunc("/products/search", productHandler.HandleProductSearch)
	http.HandleFunc("/products/categories", productHandler.HandleProductCategories)

	log.Printf("Listening on port %v", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting the HTTP server: %v", err)
	}
}
