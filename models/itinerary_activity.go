package models

import (
	"database/sql"
	"time"
)

type ItineraryActivity struct {
	ActivityID   int64         `json:"activity_id"`
	ItineraryID  int64         `json:"itinerary_id"`
	ActivityType string        `json:"activity_type"`
	HotelID      sql.NullInt64 `json:"hotel_id"`
	AttractionID sql.NullInt64 `json:"attraction_id"`
	RestaurantID sql.NullInt64 `json:"restaurant_id"`
	FlightID     sql.NullInt64 `json:"flight_id"`
	OrderNumber  int           `json:"order_number"`
	StartTime    time.Time     `json:"start_time"`
	EndTime      time.Time     `json:"end_time"`
	Notes        string        `json:"notes"`
	CreatedAt    time.Time     `json:"created_at"`

	EntityName    string  `json:"entity_name"`
	EntityDetail  string  `json:"entity_detail"`
	EntityExtra   string  `json:"entity_extra"`
	EntityRating  float64 `json:"entity_rating"`
	EntityPrice   float64 `json:"entity_price"`
	EntityWebsite string  `json:"entity_website"`
	HotelStars    int     `json:"hotel_stars"`
}
