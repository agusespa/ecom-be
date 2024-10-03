package repository

import (
	"database/sql"
	"net/http"

	"github.com/agusespa/ecom-be/customer/internal/errors"
	"github.com/agusespa/ecom-be/customer/internal/models"
	"github.com/go-sql-driver/mysql"
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

func (repo *CustomerRepository) CreateUser(uuid string, body models.UserRequest, passwordHash []byte) (int64, error) {
	var middleName *string
	if body.MiddleName == "" {
		middleName = nil
	} else {
		middleName = &body.MiddleName
	}
	result, err := repo.DB.Exec("INSERT INTO users (user_uuid, first_name, middle_name, last_name, email, password_hash) VALUES (?, ?, ?, ?, ?, ?)", uuid, body.FirstName, middleName, body.LastName, body.Email, passwordHash)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			err := httperrors.NewError(err, http.StatusConflict)
			return 0, err
		}
		err := httperrors.NewError(err, http.StatusInternalServerError)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		err := httperrors.NewError(err, http.StatusInternalServerError)
		return 0, err
	}

	return id, nil
}
