package models

import "time"

type Transportation struct {
	TransportID     int       `json:"transport_id" db:"transport_id"`
	FromCityID      int       `json:"from_city_id" db:"from_city_id"`
	ToCityID        int       `json:"to_city_id" db:"to_city_id"`
	Type            string    `json:"type" db:"type"`
	Carrier         string    `json:"carrier" db:"carrier"`
	DurationMinutes int       `json:"duration_minutes" db:"duration_minutes"`
	Price           float64   `json:"price" db:"price"`
	Currency        string    `json:"currency" db:"currency"`
	Website         string    `json:"website" db:"website"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
