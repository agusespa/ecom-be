package models

import "time"

type CustomerEntity struct {
	ID         int64     `db:"customer_id"`
	UUID       string    `db:"uuid"`
	Email      string    `db:"email"`
	FirstName  string    `db:"first_name"`
	MiddleName string    `db:"middle_name"`
	LastName   string    `db:"last_name"`
	CreatedAt  time.Time `db:"created_at"`
}
