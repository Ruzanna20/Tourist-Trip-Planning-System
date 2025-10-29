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

func (r *CityRepository) Insert(city *models.City) (int, error) {
	existing, err := r.GetByName(city.Name, city.CountryID)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return existing.CityID, nil
	}
	query := `INSERT INTO cities (country_id,name,latitude,longitude,description,created_at)
			  VALUES ($1,$2,$3,$4,$5,$6)
			  RETURNING city_id;`

	var lastInsertID int

	if city.Description == "" {
		city.Description = "No description provided."
	}

	err = r.db.QueryRow(
		query,
		city.CountryID,
		city.Name,
		city.Latitude,
		city.Longitude,
		city.Description,
		time.Now(),
	).Scan(&lastInsertID)

	if err != nil {
		log.Printf("ERROR inserting city %s (CountryID %d): %v", city.Name, city.CountryID, err)
		return 0, fmt.Errorf("failed to insert city %s: %w", city.Name, err)
	}

	return lastInsertID, nil
}

func (r *CityRepository) GetByName(name string, countryID int) (*models.City, error) {
	city := models.City{}

	query := `SELECT city_id,country_id,name,latitude,longitude
			  FROM cities 
			  WHERE name = $1 AND country_id = $2;`

	err := r.db.QueryRow(query, name, countryID).Scan(
		&city.CityID,
		&city.CountryID,
		&city.Name,
		&city.Latitude,
		&city.Longitude,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching city %s: %w", name, err)
	}

	return &city, nil
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