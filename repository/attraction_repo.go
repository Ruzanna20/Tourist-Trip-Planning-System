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

func (r *AttractionRepository) Upsert(attraction *models.Attraction) (int, error) {
	query := `INSERT INTO attractions (city_id, name, category, latitude, longitude, rating, entry_fee, currency, 
           opening_hours, description, image_url, website, created_at, updated_at) 
          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
          ON CONFLICT (name, city_id) DO UPDATE  
          SET 
            category = EXCLUDED.category,
            latitude = EXCLUDED.latitude, 
            longitude = EXCLUDED.longitude,
            rating = EXCLUDED.rating,
            entry_fee = EXCLUDED.entry_fee,
            description = EXCLUDED.description,
            currency = $8,
            opening_hours = $9,
            image_url = $11,
            website = $12,
            updated_at = NOW() AT TIME ZONE 'Asia/Yerevan' 
          RETURNING attraction_id;`

	if attraction.CreatedAt.IsZero() {
		attraction.CreatedAt = time.Now()
	}

	if attraction.UpdatedAt.IsZero() {
		attraction.UpdatedAt = time.Now()
	}

	var attractionID int
	err := r.db.QueryRow(
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
		attraction.CreatedAt,
		attraction.UpdatedAt,
	).Scan(&attractionID)

	if err != nil {
		log.Printf("ERROR inserting attraction %s (CityID %d): %v", attraction.Name, attraction.CityID, err)
		return 0, fmt.Errorf("failed to insert attraction %s: %w", attraction.Name, err)
	}
	return attractionID, nil
}

func (r *AttractionRepository) GetAllAttractions() ([]models.Attraction, error) {
	query := `SELECT 
                attraction_id, city_id, name, category, latitude, longitude, 
                rating, entry_fee, currency, opening_hours, description, 
                image_url, website, created_at, updated_at
              FROM attractions;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all attractionS:%w", err)
	}
	defer rows.Close()

	var attractions []models.Attraction
	for rows.Next() {
		var a models.Attraction
		var ratingSql, entryFeeSql sql.NullFloat64
		var currencySql, openingHoursSql, descriptionSql, imageUrlSql, websiteSql sql.NullString
		if err := rows.Scan(
			&a.AttractionID, &a.CityID, &a.Name, &a.Category, &a.Latitude, &a.Longitude,
			&ratingSql, &entryFeeSql, &currencySql, &openingHoursSql, &descriptionSql,
			&imageUrlSql, &websiteSql, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning attraction row: %v", err)
			continue
		}

		a.Rating = ratingSql.Float64
		a.EntryFee = entryFeeSql.Float64
		a.Currency = currencySql.String
		a.OpeningHours = openingHoursSql.String
		a.Description = descriptionSql.String
		a.ImageURL = imageUrlSql.String
		a.Website = websiteSql.String

		attractions = append(attractions, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return attractions, nil
}
