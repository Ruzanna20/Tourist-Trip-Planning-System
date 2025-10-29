package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"travel-planning/models"
)

type AttractionRepository struct {
	db *sql.DB
}

func NewAttractionRepository(db *sql.DB) *AttractionRepository {
	return &AttractionRepository{
		db: db,
	}
}

func (r *AttractionRepository) GetByName(name string,cityID int) (*models.Attraction,error) {
	attraction := models.Attraction{}
	query := `SELECT 
			  attraction_id,city_id,name,category,latitude,longitude,rating,entry_fee,currency,opening_hours,description,image_url,website,created_at,updated_at
			  FROM attractions
			  WHERE name = $1 AND city_id = $2`

	var (
			ratingSql sql.NullFloat64
			entryFeeSql sql.NullFloat64
			currencySql sql.NullString
			openingHoursSql sql.NullString
			descriptionSql sql.NullString
			imageUrlSql sql.NullString
			websiteSql sql.NullString
			createdAtSql time.Time 
			updatedAtSql time.Time
		)

	err := r.db.QueryRow(query, name, cityID).Scan(
		&attraction.AttractionID,
		&attraction.CityID,
		&attraction.Name,
		&attraction.Category,
		&attraction.Latitude,
		&attraction.Longitude,
		&ratingSql,
		&entryFeeSql,
		&currencySql,
		&openingHoursSql,
		&descriptionSql,
		&imageUrlSql,
		&websiteSql,
		&createdAtSql,
		&updatedAtSql,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error fetching attraction %s: %w", name, err)
	}

	attraction.Rating = ratingSql.Float64
	attraction.EntryFee = entryFeeSql.Float64
	attraction.Currency = currencySql.String
	attraction.OpeningHours = openingHoursSql.String
	attraction.Description = descriptionSql.String
	attraction.ImageURL = imageUrlSql.String
	attraction.Website = websiteSql.String
	attraction.CreatedAt = createdAtSql
	attraction.UpdatedAt = updatedAtSql

	return &attraction, nil
}

func (r *AttractionRepository) Insert(attraction *models.Attraction) (int, error) {
	existing, err := r.GetByName(attraction.Name, attraction.CityID)
	if err != nil {
		return 0, err
	}
	if existing != nil {
		return existing.AttractionID, nil 
	}
	
	query := `INSERT INTO attractions 
		(city_id, name, category, latitude, longitude, rating, entry_fee, currency, 
		opening_hours, description, image_url, website, created_at, updated_at) 
		VALUES 
		($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING attraction_id;`

	var lastInsertID int
    currentTime := time.Now()

	err = r.db.QueryRow(
		query,
		attraction.CityID,
		attraction.Name,
		attraction.Category,
		attraction.Latitude,
		attraction.Longitude,
		attraction.Rating,
		attraction.EntryFee,
		attraction.Currency,
		attraction.OpeningHours,
		attraction.Description,
		attraction.ImageURL,
		attraction.Website,
		currentTime,
		currentTime,
	).Scan(&lastInsertID)

	if err != nil {
		log.Printf("ERROR inserting attraction %s (CityID %d): %v", attraction.Name, attraction.CityID, err)
		return 0, fmt.Errorf("failed to insert attraction %s: %w", attraction.Name, err)
	}
	return lastInsertID, nil
} 