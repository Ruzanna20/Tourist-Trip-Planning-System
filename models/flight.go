package models

import "time"

type Flight struct {
	FlightID        int       `json:"flight_id" db:"flight_id"`
	FromCityID      int       `json:"from_city_id" db:"from_city_id"`
	ToCityID        int       `json:"to_city_id" db:"to_city_id"`
	Airline         string    `json:"airline" db:"airline"`
	DurationMinutes int       `json:"duration_minutes" db:"duration_minutes"`
	Price           float64   `json:"price" db:"price"`
	Website         string    `json:"website" db:"website"`
	CreatedAt       time.Time `json:"-" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
