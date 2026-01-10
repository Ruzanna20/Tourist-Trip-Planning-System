package models

import "time"

type Hotel struct {
	HotelID       int       `json:"hotel_id" db:"hotel_id"`
	CityID        int       `json:"city_id" db:"city_id"`
	Name          string    `json:"name" db:"name"`
	Address       string    `json:"address" db:"address"`
	Stars         int       `json:"stars" db:"stars"`
	Rating        float64   `json:"rating" db:"rating"`
	PricePerNight float64   `json:"price_per_night" db:"price_per_night"`
	Currency      string    `json:"currency" db:"currency"`
	Phone         string    `json:"phone" db:"phone"`
	Website       string    `json:"website" db:"website"`
	ImageURL      string    `json:"image_url" db:"image_url"`
	Description   string    `json:"description" db:"description"`
	CreatedAt     time.Time `json:"-" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}
