package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	"travel-planning/models"
)

type TripItineraryRepository struct {
	db *sql.DB
}

func NewTripItineraryRepository(db *sql.DB) *TripItineraryRepository {
	return &TripItineraryRepository{db: db}
}

func (r *TripItineraryRepository) Insert(tx *sql.Tx, itinerary *models.TripItinerary) (int, error) {
	query := `INSERT INTO trip_itinerary (trip_id, day_number, notes, date, created_at)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING itinerary_id;`

	var itineraryID int
	currTime := time.Now()

	err := tx.QueryRow(
		query,
		itinerary.TripID,
		itinerary.DayNumber,
		itinerary.Notes,
		itinerary.Date,
		currTime,
	).Scan(&itineraryID)

	if err != nil {
		slog.Error("Failed to insert trip itinerary day",
			"trip_id", itinerary.TripID,
			"day_number", itinerary.DayNumber,
			"error", err,
		)
		return 0, fmt.Errorf("failed to insert trip itinerary: %w", err)
	}

	slog.Debug("Trip itinerary day inserted",
		"itinerary_id", itineraryID,
		"trip_id", itinerary.TripID,
		"day", itinerary.DayNumber,
	)
	return itineraryID, nil
}

func (r *TripItineraryRepository) GetItineraryDaysByTripID(tripID int) ([]*models.TripItinerary, error) {
	slog.Info("Fetching itinerary days", "trip_id", tripID)

	query := `SELECT itinerary_id, trip_id, day_number, notes, date, created_at
              FROM trip_itinerary 
              WHERE trip_id = $1 
              ORDER BY day_number ASC`

	rows, err := r.db.Query(query, tripID)
	if err != nil {
		slog.Error("Error querying trip itinerary days", "trip_id", tripID, "error", err)
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
			slog.Warn("Error scanning trip itinerary row", "error", err)
			continue
		}
		days = append(days, day)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Rows iteration error in GetItineraryDaysByTripID", "error", err)
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	slog.Debug("Itinerary days fetched successfully", "trip_id", tripID, "count", len(days))
	return days, nil
}
