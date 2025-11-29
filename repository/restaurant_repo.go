package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"travel-planning/models"
)

type RestaurantRepository struct {
	db *sql.DB
}

func NewRestaurantRepository(db *sql.DB) *RestaurantRepository {
	return &RestaurantRepository{
		db: db,
	}
}

func (r *RestaurantRepository) Upsert(restaurant *models.Restaurant) (int, error) {
	query := `INSERT INTO restaurants (
        city_id, name, cuisine_type, address, rating, price_range, phone, 
        website, image_url, opening_hours, created_at, updated_at
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
    ON CONFLICT (name, city_id) DO UPDATE 
    SET 
        cuisine_type = EXCLUDED.cuisine_type,
        address = EXCLUDED.address,
        rating = EXCLUDED.rating,
        price_range = EXCLUDED.price_range,
        opening_hours = EXCLUDED.opening_hours,
        image_url = EXCLUDED.image_url,
        phone = COALESCE(EXCLUDED.phone, restaurants.phone),
        website = COALESCE(EXCLUDED.website, restaurants.website),
        updated_at = NOW() AT TIME ZONE 'Asia/Yerevan' 
    RETURNING restaurant_id;`

	if restaurant.CreatedAt.IsZero() {
		restaurant.CreatedAt = time.Now()
	}

	if restaurant.UpdatedAt.IsZero() {
		restaurant.UpdatedAt = time.Now()
	}

	var restaurantID int
	err := r.db.QueryRow(
		query,
		restaurant.CityID,
		restaurant.Name,
		restaurant.CuisineType,
		restaurant.Address,
		restaurant.Rating,
		restaurant.PriceRange,
		restaurant.Phone,
		restaurant.Website,
		restaurant.ImageURL,
		restaurant.OpeningHours,
		restaurant.CreatedAt,
		restaurant.UpdatedAt,
	).Scan(&restaurantID)

	if err != nil {
		return 0, fmt.Errorf("failed to upsert hotel %s: %w", restaurant.Name, err)
	}
	return restaurantID, nil
}

func (r *RestaurantRepository) GetAllRestaurants() ([]models.Restaurant, error) {
	query := `SELECT 
                restaurant_id, city_id, name, cuisine_type, address, rating, price_range, 
                phone, website, image_url, opening_hours, created_at, updated_at
              FROM restaurants;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all restaurants: %w", err)
	}
	defer rows.Close()

	var restaurants []models.Restaurant
	for rows.Next() {
		var r models.Restaurant
		var ratingSql sql.NullFloat64
		var cuisineSql, addressSql, priceRangeSql, phoneSql, websiteSql, imageUrlSql, openingHoursSql sql.NullString

		if err := rows.Scan(
			&r.RestaurantID, &r.CityID, &r.Name,
			&cuisineSql, &addressSql, &ratingSql, &priceRangeSql,
			&phoneSql, &websiteSql, &imageUrlSql, &openingHoursSql,
			&r.CreatedAt, &r.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning restaurant row: %v", err)
			continue
		}

		r.CuisineType = cuisineSql.String
		r.Address = addressSql.String
		r.Rating = ratingSql.Float64
		r.PriceRange = priceRangeSql.String
		r.Phone = phoneSql.String
		r.Website = websiteSql.String
		r.ImageURL = imageUrlSql.String
		r.OpeningHours = openingHoursSql.String
		restaurants = append(restaurants, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return restaurants, nil
}
