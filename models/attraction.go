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
	Website      string    `json:"website" db:"website"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}
