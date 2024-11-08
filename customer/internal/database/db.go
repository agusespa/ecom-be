package database

import (
	"database/sql"

	"github.com/agusespa/ecom-be/customer/internal/models"
	"github.com/go-sql-driver/mysql"
)

func ConnectDB(config models.Database) (*sql.DB, error) {
	cfg := mysql.Config{
		User:      config.User,
		Passwd:    config.Password,
		Net:       "tcp",
		Addr:      config.Address,
		DBName:    "ecom_customer",
		ParseTime: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
