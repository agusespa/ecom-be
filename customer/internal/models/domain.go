package models

import (
	"database/sql"
	"time"
)

type Customer struct {
	CustomerID   int64     `json:"customer_id"`
	CustomerUUID string    `json:"uuid"`
	Email        string    `json:"email"`
	FirstName    string    `json:"first_name"`
	MiddleName   string    `json:"middle_name"`
	LastName     string    `json:"last_name"`
	CreatedAt    time.Time `json:"created_at"`
}

type CustomerRequest struct {
	CustomerUUID string `json:"uuid"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	FirstName    string `json:"first_name"`
	MiddleName   string `json:"middle_name"`
	LastName     string `json:"last_name"`
}

type RegistrationResponse struct {
	CustomerID int64 `json:"customerID"`
}

func NewCustomer(id int64, uuid, email, firstName string, middleNameNullStr sql.NullString, lastName string, createdAt time.Time) Customer {
	return Customer{
		CustomerID:   id,
		CustomerUUID: uuid,
		Email:        email,
		FirstName:    firstName,
		MiddleName:   middleNameNullStr.String,
		LastName:     lastName,
		CreatedAt:    createdAt,
	}
}
