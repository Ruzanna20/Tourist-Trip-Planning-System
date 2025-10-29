package models

import "time"

type City struct {
	CityID      int       `json:"city_id" db:"city_id"`
	CountryID   int       `json:"country_id" db:"country_id"`
	Name        string    `json:"name" db:"name"`
	Latitude    float64   `json:"latitude" db:"latitude"`
	Longitude   float64   `json:"longitude" db:"longitude"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}
