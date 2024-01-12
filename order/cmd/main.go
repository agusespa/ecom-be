package main

import (
	"log"
	"order/internal/adapters/db"
	"order/internal/application/api"
)

func main() {
	dbAdapter, err := db.NewAdapter()
	if err != nil {
		log.Fatalf("Failed to connect to database. Error: %v", err)
	}

	application := api.NewApplication(dbAdapter)
}
