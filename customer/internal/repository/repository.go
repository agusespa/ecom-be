package repository

import (
	"database/sql"
	"net/http"

	"github.com/agusespa/ecom-be/customer/internal/httperrors"
	"github.com/agusespa/ecom-be/customer/internal/models"
	"github.com/go-sql-driver/mysql"
)

type CustomerRepository interface {
	ReadCustomerByID(id int64) (models.CustomerEntity, error)
	ReadCustomerUUID(id int64) (string, error)
	DeleteCustomer(id int64) error
	CreateCustomer(body models.CustomerRequest) (int64, error)
}

type MySqlRepository struct {
	DB *sql.DB
}

func NewMySqlRepository(db *sql.DB) *MySqlRepository {
	return &MySqlRepository{DB: db}
}

func (repo *MySqlRepository) ReadCustomerUUID(id int64) (string, error) {
	var uuid string

	query := `
		SELECT 
			c.customer_uuid
		FROM customers c
		WHERE c.customer_id = ?
	`

	row := repo.DB.QueryRow(query, id)
	err := row.Scan(&uuid)

	if err != nil {
		if err == sql.ErrNoRows {
			err = httperrors.NewError(err, http.StatusNotFound)
			return "", err
		}
		err = httperrors.NewError(err, http.StatusInternalServerError)
		return "", err
	}

	return uuid, nil
}

func (repo *MySqlRepository) ReadCustomerByID(id int64) (models.CustomerEntity, error) {
	var customer models.CustomerEntity

	query := `
		SELECT 
			c.customer_id, 
			c.customer_uuid, 
			c.first_name, 
			c.middle_name, 
			c.last_name, 
			c.email, 
			c.created_at
		FROM customers c
		WHERE c.customer_id = ?
	`

	row := repo.DB.QueryRow(query, id)
	err := row.Scan(
		&customer.CustomerID,
		&customer.CustomerUUID,
		&customer.FirstName,
		&customer.MiddleName,
		&customer.LastName,
		&customer.Email,
		&customer.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			err = httperrors.NewError(err, http.StatusNotFound)
			return customer, err
		}
		err = httperrors.NewError(err, http.StatusInternalServerError)
		return customer, err
	}

	return customer, nil
}

func (repo *MySqlRepository) CreateCustomer(body models.CustomerRequest) (int64, error) {
	var middleName *string
	if body.MiddleName == "" {
		middleName = nil
	} else {
		middleName = &body.MiddleName
	}

	query := `
		INSERT INTO customers (customer_uuid, first_name, middle_name, last_name, email)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := repo.DB.Exec(query, body.CustomerUUID, body.FirstName, middleName, body.LastName, body.Email)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			err = httperrors.NewError(err, http.StatusConflict)
			return 0, err
		}
		err = httperrors.NewError(err, http.StatusInternalServerError)
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		err = httperrors.NewError(err, http.StatusInternalServerError)
		return 0, err
	}

	return id, nil
}

func (repo *MySqlRepository) DeleteCustomer(id int64) error {
	_, err := repo.DB.Exec("DELETE FROM customers WHERE customer_id = ?", id)
	if err != nil {
		err = httperrors.NewError(err, http.StatusInternalServerError)
		return err
	}
	return nil
}
