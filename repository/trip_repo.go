package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	"travel-planning/models"
)

type TripRepository struct {
	db *sql.DB
}

func NewTripRepository(db *sql.DB) *TripRepository {
	return &TripRepository{
		db: db,
	}
}

func (r *TripRepository) GetConn() *sql.DB {
	return r.db
}

func (r *TripRepository) Insert(tx *sql.Tx, trip *models.Trip) (int, error) {
	query := `INSERT INTO trips (user_id, destination_city_id, title, start_date, end_date, total_price, status, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    RETURNING trip_id;`

	var tripID int
	currTime := time.Now()
	err := tx.QueryRow(
		query,
		trip.UserID,
		trip.DestinationCityID,
		trip.Title,
		trip.StartDate,
		trip.EndDate,
		trip.TotalPrice,
		trip.Status,
		currTime,
		currTime,
	).Scan(&tripID)

	if err != nil {
		slog.Error("Failed to insert new trip",
			"user_id", trip.UserID,
			"title", trip.Title,
			"error", err,
		)
		return 0, fmt.Errorf("failed to insert new trip: %w", err)
	}

	slog.Debug("Trip inserted successfully", "trip_id", tripID, "user_id", trip.UserID)
	return tripID, nil
}

func (r *TripRepository) GetAllTripsByUserID(userID int) ([]models.Trip, error) {
	slog.Info("Fetching all trips for user", "user_id", userID)

	query := `SELECT 
                trip_id, user_id, destination_city_id, title, start_date, end_date, 
                total_price, status, created_at, updated_at
              FROM trips 
              WHERE user_id = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		slog.Error("Failed to fetch trips for user", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to fetch trips for user %d: %w", userID, err)
	}
	defer rows.Close()

	var trips []models.Trip

	for rows.Next() {
		var t models.Trip
		var totalPriceSql sql.NullFloat64
		var statusSql sql.NullString

		if err := rows.Scan(
			&t.TripID, &t.UserID, &t.DestinationCityID, &t.Title,
			&t.StartDate, &t.EndDate, &totalPriceSql,
			&statusSql, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			slog.Warn("Error scanning trip row", "user_id", userID, "error", err)
			continue
		}
		t.TotalPrice = totalPriceSql.Float64
		t.Status = statusSql.String

		trips = append(trips, t)
	}

	return trips, nil
}

func (r *TripRepository) GetTripByID(tripID int) (*models.Trip, error) {
	slog.Debug("Fetching trip by ID", "trip_id", tripID)

	query := `SELECT 
                trip_id, user_id, destination_city_id, title, start_date, end_date, 
                total_price, status, created_at, updated_at
              FROM trips 
              WHERE trip_id = $1`

	var t models.Trip
	var totalPriceSql sql.NullFloat64
	var statusSql sql.NullString

	err := r.db.QueryRow(query, tripID).Scan(
		&t.TripID,
		&t.UserID,
		&t.DestinationCityID,
		&t.Title,
		&t.StartDate,
		&t.EndDate,
		&totalPriceSql,
		&statusSql,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.Warn("Trip not found", "trip_id", tripID)
			return nil, fmt.Errorf("trip with id %d not found", tripID)
		}
		slog.Error("Database error in GetTripByID", "trip_id", tripID, "error", err)
		return nil, err
	}

	t.TotalPrice = totalPriceSql.Float64
	t.Status = statusSql.String

	return &t, nil
}

func (r *TripRepository) DeleteByIDAndUserID(tripID, userID int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRow(
		`SELECT EXISTS(SELECT 1 FROM trips WHERE trip_id = $1 AND user_id = $2)`,
		tripID, userID,
	).Scan(&exists)

	if err != nil {
		return fmt.Errorf("failed to verify trip ownership: %w", err)
	}
	if !exists {
		return fmt.Errorf("trip not found or access denied")
	}

	_, err = tx.Exec(`
        DELETE FROM itinerary_activities 
        WHERE itinerary_id IN (SELECT itinerary_id FROM trip_itinerary WHERE trip_id = $1)`,
		tripID)
	if err != nil {
		slog.Error("Failed to delete trip activities", "trip_id", tripID, "error", err)
		return fmt.Errorf("failed to delete activities: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM trip_itinerary WHERE trip_id = $1`, tripID)
	if err != nil {
		slog.Error("Failed to delete trip itinerary", "trip_id", tripID, "error", err)
		return fmt.Errorf("failed to delete itinerary: %w", err)
	}

	_, err = tx.Exec(`DELETE FROM trips WHERE trip_id = $1 AND user_id = $2`, tripID, userID)
	if err != nil {
		slog.Error("Failed to delete trip", "trip_id", tripID, "error", err)
		return fmt.Errorf("failed to delete trip: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	slog.Info("Trip and all related data deleted successfully", "trip_id", tripID, "user_id", userID)
	return nil
}
