package models

import "time"

type Payment struct {
	PaymentID     int       `json:"payment_id" db:"payment_id"`
	BookingID     int       `json:"booking_id" db:"booking_id"`
	Amount        float64   `json:"amount" db:"amount"`
	Currency      string    `json:"currency" db:"currency"`
	PaymentDate   time.Time `json:"payment_date" db:"payment_date"`
	PaymentMethod string    `json:"payment_method" db:"payment_method"`
	Status        string    `json:"status" db:"status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
