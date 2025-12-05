package repository

import (
	"database/sql"
	"fmt"
	"time"
	"travel-planning/models"
)

type TripItineraryRepository struct {
	db *sql.DB
}

func NewTripItineraryRepository(db *sql.DB) *TripItineraryRepository {
	return &TripItineraryRepository{db: db}
}

func (r *TripItineraryRepository) Insert(itinerary *models.TripItinerary) (int, error) {
	query := `INSERT INTO trip_itinerary (
        trip_id, day_number, notes, date, created_at
    )
    VALUES ($1, $2, $3, $4, $5)
    RETURNING itinerary_id;`

	var itineraryID int
	currTime := time.Now()

	err := r.db.QueryRow(
		query,
		itinerary.TripID,
		itinerary.DayNumber,
		itinerary.Notes,
		itinerary.Date,
		currTime,
	).Scan(&itineraryID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert trip itinerary: %w", err)
	}

	return itineraryID, nil
}
