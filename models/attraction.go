package models

import "time"

type Attraction struct {
	AttractionID int       `json:"attraction_id" db:"attraction_id"`
	CityID       int       `json:"city_id" db:"city_id"`
	Name         string    `json:"name" db:"name"`
	Category     string    `json:"category" db:"category"`
	Latitude     float64   `json:"latitude" db:"latitude"`
	Longitude    float64   `json:"longitude" db:"longitude"`
	Rating       float64   `json:"rating" db:"rating"`
	EntryFee     float64   `json:"entry_fee" db:"entry_fee"`
	Currency     string    `json:"currency" db:"currency"`
	OpeningHours string    `json:"opening_hours" db:"opening_hours"`
	Description  string    `json:"description" db:"description"`
	ImageURL     string    `json:"image_url" db:"image_url"`
	Website      string    `json:"website" db:"website"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
