package models

import (
	"database/sql"
	"time"
)

type CustomerEntity struct {
	CustomerID   int64          `db:"customer_id"`
	CustomerUUID string         `db:"customer_uuid"`
	Email        string         `db:"email"`
	FirstName    string         `db:"first_name"`
	MiddleName   sql.NullString `db:"middle_name"`
	LastName     string         `db:"last_name"`
	CreatedAt    time.Time      `db:"created_at"`
}
