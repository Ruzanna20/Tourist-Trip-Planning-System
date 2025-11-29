package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"travel-planning/models"
)

type CityRepository struct {
	db *sql.DB
}

func NewCityRepository(db *sql.DB) *CityRepository {
	return &CityRepository{
		db: db,
	}
}

func (r *CityRepository) Upsert(city *models.City) (int, error) {
	query := `INSERT INTO cities  (country_id, name, latitude, longitude, description, created_at, updated_at) 
			VALUES ($1, $2, $3, $4, $5, $6,$7) 
			ON CONFLICT (latitude, longitude, country_id) DO UPDATE
			SET 
				name = EXCLUDED.name,
				description = EXCLUDED.description,
				updated_at = NOW() AT TIME ZONE 'Asia/Yerevan'
			RETURNING city_id;`

	if city.Description == "" {
		city.Description = "No description provided."
	}

	if city.CreatedAt.IsZero() {
		city.CreatedAt = time.Now()
	}

	if city.UpdatedAt.IsZero() {
		city.UpdatedAt = time.Now()
	}

	var cityID int
	err := r.db.QueryRow(
		query,
		city.CountryID,
		city.Name,
		city.Latitude,
		city.Longitude,
		city.Description,
		city.CreatedAt,
		city.UpdatedAt,
	).Scan(&cityID)

	if err != nil {
		log.Printf("ERROR inserting city %s (CountryID %d): %v", city.Name, city.CountryID, err)
		return 0, fmt.Errorf("failed to insert city %s: %w", city.Name, err)
	}

	return cityID, nil
}

type CityLocation struct {
	ID        int
	Name      string
	Latitude  float64
	Longitude float64
}

func (r *CityRepository) GetAllCityLocations() ([]CityLocation, error) {
	query := `SELECT city_id, name, latitude, longitude FROM cities;`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all city locations: %w", err)
	}
	defer rows.Close()

	var locations []CityLocation
	for rows.Next() {
		var loc CityLocation
		if err := rows.Scan(&loc.ID, &loc.Name, &loc.Latitude, &loc.Longitude); err != nil {
			log.Printf("Error scanning city location row: %v", err)
			continue
		}
		locations = append(locations, loc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning city rows: %w", err)
	}

	return locations, nil
}

func (r *CityRepository) GetAllCities() ([]models.City, error) {
	query := `SELECT city_id,country_id,name,latitude,longitude,description,created_at,updated_at
			  FROM cities;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all city locations: %w", err)
	}
	defer rows.Close()

	var cities []models.City
	for rows.Next() {
		var city models.City
		if err := rows.Scan(
			&city.CityID,
			&city.CountryID,
			&city.Name,
			&city.Latitude,
			&city.Longitude,
			&city.Description,
			&city.CreatedAt,
			&city.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning city location row: %v", err)
			continue
		}
		cities = append(cities, city)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error after scanning city rows: %w", err)
	}
	return cities, nil
}
