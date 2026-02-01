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
	query := `INSERT INTO attractions (city_id, name, category, latitude, longitude, rating, entry_fee, website, created_at, updated_at) 
          VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
          ON CONFLICT (name, city_id) DO UPDATE  
          SET 
            category = EXCLUDED.category,
            latitude = EXCLUDED.latitude, 
            longitude = EXCLUDED.longitude,
            rating = EXCLUDED.rating,
            entry_fee = EXCLUDED.entry_fee,
            website = $8,
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
                rating, entry_fee, website, created_at, updated_at
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
		var websiteSql sql.NullString
		if err := rows.Scan(
			&a.AttractionID, &a.CityID, &a.Name, &a.Category, &a.Latitude, &a.Longitude,
			&ratingSql, &entryFeeSql, &websiteSql, &a.CreatedAt, &a.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning attraction row: %v", err)
			continue
		}

		a.Rating = ratingSql.Float64
		a.EntryFee = entryFeeSql.Float64
		a.Website = websiteSql.String

		attractions = append(attractions, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return attractions, nil
}

func (s *AttractionRepository) GetBestAttractionsByTier(cityID int, budgetLimit float64, tier string) ([]models.Attraction, error) {
	var orderBy string
	switch tier {
	case "Economy":
		orderBy = "entry_fee ASC, rating DESC"
	case "Luxury":
		orderBy = "rating DESC, entry_fee DESC"
	default:
		orderBy = "rating DESC, entry_fee ASC"
	}

	query := fmt.Sprintf(`
		SELECT attraction_id, city_id, name, category, latitude, longitude, rating, entry_fee, website
		FROM attractions 
		WHERE city_id = $1 AND entry_fee <= $2
		ORDER BY %s
		LIMIT 10`, orderBy)

	rows, err := s.db.Query(query, cityID, budgetLimit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.Attraction
	for rows.Next() {
		var a models.Attraction
		if err := rows.Scan(
			&a.AttractionID,
			&a.CityID,
			&a.Name,
			&a.Category,
			&a.Latitude,
			&a.Longitude,
			&a.Rating,
			&a.EntryFee,
			&a.Website); err != nil {
			continue
		}
		results = append(results, a)
	}
	return results, nil
}
