package models

import (
	"database/sql"
	"time"
)

type ItineraryActivity struct {
	ActivityID   int64         `json:"activity_id" db:"activity_id"`
	ItineraryID  int64         `json:"itinerary_id" db:"itinerary_id"`
	ActivityType string        `json:"activity_type" db:"activity_type"`
	HotelID      sql.NullInt64 `json:"hotel_id" db:"hotel_id"`
	AttractionID sql.NullInt64 `json:"attraction_id" db:"attraction_id"`
	RestaurantID sql.NullInt64 `json:"restaurant_id" db:"restaurant_id"`
	FlightID     sql.NullInt64 `json:"flight_id" db:"flight_id"`
	OrderNumber  int           `json:"order_number" db:"order_number"`
	StartTime    time.Time     `json:"start_time" db:"start_time"`
	EndTime      time.Time     `json:"end_time" db:"end_time"`
	Notes        string        `json:"notes" db:"notes"`

	EntityName   string    `json:"entity_name"`
	EntityDetail string    `json:"entity_detail"`
	EntityExtra  string    `json:"entity_extra"`
	EntityRating float64   `json:"entity_rating"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
