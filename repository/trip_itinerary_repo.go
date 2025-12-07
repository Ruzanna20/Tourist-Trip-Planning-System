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
	query := `INSERT INTO trip_itinerary (trip_id, day_number, notes, date, created_at)
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

func (r *TripItineraryRepository) GetItineraryDaysByTripID(tripID int) ([]*models.TripItinerary, error) {

	query := `SELECT itinerary_id, trip_id, day_number, notes, date, created_at
	          FROM trip_itinerary 
	          WHERE trip_id = $1 
	          ORDER BY day_number ASC`

	rows, err := r.db.Query(query, tripID)
	if err != nil {
		return nil, fmt.Errorf("error querying trip itinerary days: %w", err)
	}
	defer rows.Close()

	var days []*models.TripItinerary

	for rows.Next() {
		day := &models.TripItinerary{}

		err := rows.Scan(
			&day.ItineraryID,
			&day.TripID,
			&day.DayNumber,
			&day.Notes,
			&day.Date,
			&day.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning trip itinerary row: %w", err)
		}
		days = append(days, day)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return days, nil
}
