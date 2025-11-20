package repository

import (
	"database/sql"
	"fmt"
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
