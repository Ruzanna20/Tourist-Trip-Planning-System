package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"travel-planning/models"
)

type FlightRepository struct {
	db *sql.DB
}

func NewFlightRepository(db *sql.DB) *FlightRepository {
	return &FlightRepository{
		db: db,
	}
}

func (r *FlightRepository) Upsert(flight *models.Flight) (int, error) {
	query := `INSERT INTO flights (
        from_city_id, to_city_id, airline, duration_minutes, price, 
    	website, created_at, updated_at
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    ON CONFLICT (from_city_id, to_city_id, airline) DO UPDATE 
    SET 
        duration_minutes = EXCLUDED.duration_minutes,
        price = EXCLUDED.price,
        website = COALESCE(EXCLUDED.website, flights.website),
        updated_at = NOW() AT TIME ZONE 'Asia/Yerevan' 
    RETURNING flight_id;`

	if flight.CreatedAt.IsZero() {
		flight.CreatedAt = time.Now()
	}

	if flight.UpdatedAt.IsZero() {
		flight.UpdatedAt = time.Now()
	}

	var flightID int
	err := r.db.QueryRow(
		query,
		flight.FromCityID,
		flight.ToCityID,
		flight.Airline,
		flight.DurationMinutes,
		flight.Price,
		flight.Website,
		flight.CreatedAt,
		flight.UpdatedAt,
	).Scan(&flightID)

	if err != nil {
		return 0, fmt.Errorf("ERROR upserting flight between %d and %d: %w", flight.FromCityID, flight.ToCityID, err)
	}
	return flightID, nil
}

func (r *FlightRepository) GetAllFlights() ([]models.Flight, error) {
	query := `SELECT flight_id, from_city_id, to_city_id, airline, duration_minutes, price, website, created_at, updated_at
              FROM flights;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all hotels: %w", err)
	}
	defer rows.Close()

	var flights []models.Flight
	for rows.Next() {
		var f models.Flight
		var websiteSql sql.NullString

		if err := rows.Scan(
			&f.FlightID,
			&f.FromCityID,
			&f.ToCityID,
			&f.Airline,
			&f.DurationMinutes,
			&f.Price,
			&websiteSql,
			&f.CreatedAt,
			&f.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning flight row: %v", err)
			continue
		}
		f.Website = websiteSql.String
		flights = append(flights, f)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows interation error: %w", err)
	}
	return flights, nil
}

func (r *FlightRepository) GetBestFlightByTier(fromCityID, toCityID int, budgetMax float64, tier string) (*models.Flight, error) {
	var orderBy string

	switch tier {
	case "Economy":
		orderBy = "price ASC"
	case "Balanced":
		orderBy = "price ASC"
	case "Luxury":
		orderBy = "price DESC"
	default:
		orderBy = "price ASC"
	}

	query := fmt.Sprintf(`SELECT
	flight_id, from_city_id, to_city_id, airline, duration_minutes, price, website
	FROM flights
	WHERE from_city_id = $1 AND to_city_id = $2 AND price <= $3
	ORDER BY %s
	LIMIT 1`, orderBy)

	flight := &models.Flight{}
	err := r.db.QueryRow(
		query, fromCityID, toCityID, budgetMax).Scan(
		&flight.FlightID,
		&flight.FromCityID,
		&flight.ToCityID,
		&flight.Airline,
		&flight.DurationMinutes,
		&flight.Price,
		&flight.Website,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find flight: %w", err)
	}
	return flight, nil
}
