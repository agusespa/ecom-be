package main

import (
	"log"
	"net/http"

	"github.com/agusespa/ecom-be/product/internal/database"
	"github.com/agusespa/ecom-be/product/internal/handlers"
	"github.com/agusespa/ecom-be/product/internal/repository"
	"github.com/agusespa/ecom-be/product/internal/service"
)

func main() {
	db, dbErr := database.ConnectDB()
	if dbErr != nil {
		log.Fatalf("Error establishing database connection: %v", dbErr)
	}

	// TODO: get port dinamically
	port := "3004"

	productRepository := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepository)
	productHandler := handlers.NewProductHandler(productService)

	http.HandleFunc("/productapi/products", productHandler.HandleProducts)
	http.HandleFunc("/productapi/products/", productHandler.HandleProductByID)
	http.HandleFunc("/productapi/products/search", productHandler.HandleProductSearch)
	http.HandleFunc("/productapi/products/categories", productHandler.HandleProductCategories)

	log.Printf("Listening on port %v", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting the HTTP server: %v", err)
	}
}
