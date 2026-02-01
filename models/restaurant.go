package models

import "time"

type Restaurant struct {
	RestaurantID int       `json:"restaurant_id" db:"restaurant_id"`
	CityID       int       `json:"city_id" db:"city_id"`
	Name         string    `json:"name" db:"name"`
	Cuisine      string    `json:"cuisine" db:"cuisine"`
	Latitude     float64   `json:"latitude" db:"latitude"`
	Longitude    float64   `json:"longitude" db:"longitude"`
	Rating       float64   `json:"rating" db:"rating"`
	PriceRange   string    `json:"price_range" db:"price_range"`
	Website      string    `json:"website" db:"website"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
