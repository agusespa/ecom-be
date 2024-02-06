package repository

import (
	"database/sql"
	"net/http"

	"github.com/agusespa/ecom-be/customer/internal/errors"
	"github.com/agusespa/ecom-be/customer/internal/models"
)

type CustomerRepository struct {
	DB *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{DB: db}
}

func (repo *CustomerRepository) QueryCustomerByUUID(uuid string) (models.CustomerEntity, error) {
	row := repo.DB.QueryRow("SELECT * FROM customer WHERE uuid = ?", uuid)
	var customer models.CustomerEntity

	if err := row.Scan(
		&customer.ID,
		&customer.UUID,
		&customer.Email,
		&customer.FirstName,
		&customer.MiddleName,
		&customer.LastName,
		&customer.CreatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			error := errors.NewError(err, http.StatusNotFound)
			return customer, error
		}
		error := errors.NewError(err, http.StatusInternalServerError)
		return customer, error
	}
	return customer, nil
}
