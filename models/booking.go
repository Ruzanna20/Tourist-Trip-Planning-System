package models

import "time"

type Booking struct {
	BookingID   int       `json:"booking_id" db:"booking_id"`
	TripID      int       `json:"trip_id" db:"trip_id"`
	UserID      int       `json:"user_id" db:"user_id"`
	BookingDate time.Time `json:"booking_date" db:"booking_date"`
	Status      string    `json:"status" db:"status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
