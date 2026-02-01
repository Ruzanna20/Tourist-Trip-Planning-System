package repository

import (
	"database/sql"
	"fmt"
	"log"
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

func (r *TripRepository) Insert(trip *models.Trip) (int, error) {
	query := `INSERT INTO trips (user_id, destination_city_id, title, start_date, end_date, total_price, currency, status, created_at, updated_at)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    RETURNING trip_id;`

	var tripID int
	currTime := time.Now()
	err := r.db.QueryRow(
		query,
		trip.UserID,
		trip.DestinationCityID,
		trip.Title,
		trip.StartDate,
		trip.EndDate,
		trip.TotalPrice,
		trip.Currency,
		trip.Status,
		currTime,
		currTime,
	).Scan(&tripID)

	if err != nil {
		log.Printf("DB error inserting new trip for user %d: %v", trip.UserID, err)
		return 0, fmt.Errorf("failed to insert new trip: %w", err)
	}

	return tripID, nil
}

func (r *TripRepository) GetAllTripsByUserID(userID int) ([]models.Trip, error) {
	query := `SELECT 
                trip_id, user_id, destination_city_id, title, start_date, end_date, 
                total_price, currency, status, created_at, updated_at
              FROM trips 
              WHERE user_id = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch trips for user %d: %w", userID, err)
	}
	defer rows.Close()

	var trips []models.Trip

	for rows.Next() {
		var t models.Trip
		var totalPriceSql sql.NullFloat64
		var currencySql, statusSql sql.NullString

		if err := rows.Scan(
			&t.TripID, &t.UserID, &t.DestinationCityID, &t.Title,
			&t.StartDate, &t.EndDate, &totalPriceSql, &currencySql,
			&statusSql, &t.CreatedAt, &t.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning trip row for user %d: %v", userID, err)
			continue
		}
		t.TotalPrice = totalPriceSql.Float64
		t.Currency = currencySql.String
		t.Status = statusSql.String

		trips = append(trips, t)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return trips, nil
}

func (r *TripRepository) GetTripByID(tripID int) (*models.Trip, error) {
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
			return nil, fmt.Errorf("trip with id %d not found", tripID)
		}
		return nil, err
	}

	t.TotalPrice = totalPriceSql.Float64
	t.Status = statusSql.String

	return &t, nil
}
