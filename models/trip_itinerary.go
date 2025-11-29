package models

import "time"

type TripItinerary struct {
	ItineraryID int       `json:"Itinerary_id" db:"Itinerary_id"`
	TripID      int       `json:"trip_id" db:"trip_id"`
	DayNumber   int       `json:"day_number" db:"day_number"`
	Notes       string    `json:"notes" db:"notes"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Date        time.Time `json:"date" db:"date"`
}
