package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/agusespa/ecom-be/customer/internal/helpers"
	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dbUser, dbAddr, dbPassword, err := helpers.GetDatabaseVars()
	if err != nil {
		log.Fatal(fmt.Errorf("failed to get database variables: %s", err.Error()))
	}
	cfg := mysql.Config{
		User:   dbUser,
		Passwd: dbPassword,
		Net:    "tcp",
		Addr:   dbAddr,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(fmt.Errorf("failed to establish database connection: %s", err.Error()))
	}

	defer db.Close()

	dbName := "ecom_customer"

	if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName); err != nil {
		log.Fatal(err)
	}

	if _, err := db.Exec("USE " + dbName); err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			customer_id INT AUTO_INCREMENT PRIMARY KEY,
			customer_uuid VARCHAR(36) NOT NULL UNIQUE,
			first_name VARCHAR(20) NOT NULL,
			middle_name VARCHAR(20),
			last_name VARCHAR(30) NOT NULL,
			email VARCHAR(100) NOT NULL UNIQUE,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Customers table ensured")

	fmt.Println("\nDatabase setup completed successfully")
}
