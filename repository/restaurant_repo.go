package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
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
        city_id, name, cuisine, latitude, longitude, rating, price_range, 
        website, created_at, updated_at
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    ON CONFLICT (city_id,name) DO UPDATE 
    SET 
        cuisine = EXCLUDED.cuisine,
        latitude = EXCLUDED.latitude,
        longitude = EXCLUDED.longitude,
        rating = EXCLUDED.rating,
        price_range = EXCLUDED.price_range,
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
		restaurant.Cuisine,
		restaurant.Latitude,
		restaurant.Longitude,
		restaurant.Rating,
		restaurant.PriceRange,
		restaurant.Website,
		restaurant.CreatedAt,
		restaurant.UpdatedAt,
	).Scan(&restaurantID)

	if err != nil {
		slog.Error("Failed to upsert restaurant",
			"name", restaurant.Name,
			"city_id", restaurant.CityID,
			"error", err,
		)
		return 0, fmt.Errorf("failed to upsert restaurant %s: %w", restaurant.Name, err)
	}

	slog.Debug("Restaurant upserted successfully", "id", restaurantID, "name", restaurant.Name)
	return restaurantID, nil
}

func (r *RestaurantRepository) GetAllRestaurants() ([]models.Restaurant, error) {
	query := `SELECT 
                restaurant_id, city_id, name, cuisine, latitude, longitude, rating, price_range, 
                website, created_at, updated_at
            FROM restaurants;`

	rows, err := r.db.Query(query)
	if err != nil {
		slog.Error("Failed to fetch all restaurants", "error", err)
		return nil, fmt.Errorf("failed to fetch all restaurants: %w", err)
	}
	defer rows.Close()

	var restaurants []models.Restaurant
	for rows.Next() {
		var r models.Restaurant
		var ratingSql, latitudeSql, longitudeSql sql.NullFloat64
		var cuisineSql, priceRangeSql, websiteSql sql.NullString

		if err := rows.Scan(
			&r.RestaurantID, &r.CityID, &r.Name,
			&cuisineSql, &latitudeSql, &longitudeSql, &ratingSql, &priceRangeSql,
			&websiteSql,
			&r.CreatedAt, &r.UpdatedAt,
		); err != nil {
			slog.Warn("Error scanning restaurant row", "error", err)
			continue
		}

		r.Cuisine = cuisineSql.String
		r.Latitude = latitudeSql.Float64
		r.Longitude = longitudeSql.Float64
		r.Rating = ratingSql.Float64
		r.PriceRange = priceRangeSql.String
		r.Website = websiteSql.String
		restaurants = append(restaurants, r)
	}

	return restaurants, nil
}

func (r *RestaurantRepository) GetBestRestaurantByTier(cityID int, tier string) ([]models.Restaurant, error) {
	slog.Info("Fetching best restaurants by tier", "city_id", cityID, "tier", tier)

	var priceFilter string
	switch tier {
	case "Economy":
		priceFilter = "AND price_range = '$'"
	case "Balanced":
		priceFilter = "AND price_range IN ('$', '$$')"
	case "Luxury":
		priceFilter = "AND price_range IN ('$$', '$$$')"
	default:
		priceFilter = ""
	}

	query := fmt.Sprintf(`
        SELECT restaurant_id, city_id, name, cuisine, latitude, longitude, rating, price_range, website
        FROM restaurants 
        WHERE city_id = $1 %s
        ORDER BY rating DESC
        LIMIT 15`, priceFilter)

	rows, err := r.db.Query(query, cityID)
	if err != nil {
		slog.Error("Database error in GetBestRestaurantByTier", "city_id", cityID, "error", err)
		return nil, err
	}
	defer rows.Close()

	var results []models.Restaurant
	for rows.Next() {
		var res models.Restaurant
		if err := rows.Scan(
			&res.RestaurantID,
			&res.CityID,
			&res.Name,
			&res.Cuisine,
			&res.Latitude,
			&res.Longitude,
			&res.Rating,
			&res.PriceRange,
			&res.Website); err != nil {
			slog.Warn("Skipping restaurant row due to scan error", "error", err)
			continue
		}
		results = append(results, res)
	}

	slog.Debug("Restaurants fetched successfully", "count", len(results), "city_id", cityID)
	return results, nil
}
