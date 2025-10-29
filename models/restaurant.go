package models

import "time"

type Restaurant struct {
	RestaurantID int       `json:"restaurant_id" db:"restaurant_id"`
	CityID       int       `json:"city_id" db:"city_id"`
	Name         string    `json:"name" db:"name"`
	CuisineType  string    `json:"cuisine_type" db:"cuisine_type"`
	Address      string    `json:"address" db:"address"`
	Rating       float64   `json:"rating" db:"rating"`
	PriceRange   string    `json:"price_range" db:"price_range"`
	Phone        string    `json:"phone" db:"phone"`
	Website      string    `json:"website" db:"website"`
	ImageURL     string    `json:"image_url" db:"image_url"`
	OpeningHours string    `json:"opening_hours" db:"opening_hours"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
