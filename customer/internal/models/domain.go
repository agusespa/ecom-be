package models

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Customer struct {
	ID         int64     `json:"customer_id"`
	UUID       string    `json:"uuid"`
	Email      string    `json:"email"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
	CreatedAt  time.Time `json:"created_at"`
}

func NewCustomer(id int64, uuid, email, firstName, middleName, lastName string, createdAt time.Time) Customer {
	return Customer{
		ID:         id,
		UUID:       uuid,
		Email:      email,
		FirstName:  firstName,
		MiddleName: middleName,
		LastName:   lastName,
		CreatedAt:  createdAt,
	}
}

type CustomClaims struct {
	User TokenUser `json:"user"`
	Type string    `json:"type"`
	jwt.StandardClaims
}
type TokenUser struct {
	UserID   int64  `json:"userID"`
	UserUUID string `json:"userUUID"`
}
