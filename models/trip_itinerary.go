package models

import "time"

type TripItinerary struct {
	ItineraryID   int       `json:"Itinerary_id" db:"Itinerary_id"`
	TripID        int       `json:"trip_id" db:"trip_id"`
	DayNumber     int       `json:"day_number" db:"day_number"`
	HotelID       int       `json:"hotel_id" db:"hotel_id"`
	AttractionIDs string    `json:"attaction_ids" db:"attraction_ids"`
	RestaurantIDs string    `json:"restaurant_ids" db:"restaurant_ids"`
	Notes         string    `json:"notes" db:"notes"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}
