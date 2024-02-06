package main

import (
	"log"
	"net/http"

	"github.com/agusespa/ecom-be/customer/internal/database"
	"github.com/agusespa/ecom-be/customer/internal/handlers"
	"github.com/agusespa/ecom-be/customer/internal/repository"
	"github.com/agusespa/ecom-be/customer/internal/service"
)

func main() {
	db, dbErr := database.ConnectDB()
	if dbErr != nil {
		log.Fatalf("Error establishing database connection: %v", dbErr)
	}

	// TODO: get port dinamically
	port := "3002"

	customerRepository := repository.NewCustomerRepository(db)
	customerService := service.NewCustomerService(customerRepository)
	customerHandler := handlers.NewCustomerHandler(customerService)

	http.HandleFunc("/customerapi/customer", customerHandler.HandleCustomerByUUID)

	log.Printf("Listening on port %v", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatalf("Error starting the HTTP server: %v", err)
	}
}
