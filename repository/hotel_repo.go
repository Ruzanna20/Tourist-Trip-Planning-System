package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	"travel-planning/models"
)

type HotelRepository struct {
	db *sql.DB
}

func NewHotelRepository(db *sql.DB) *HotelRepository {
	return &HotelRepository{db: db}
}

func (r *HotelRepository) Upsert(hotel *models.Hotel) (int, error) {
	query := `INSERT INTO hotels (
        city_id, name, address, stars, rating, price_per_night,
       	website, description, 
        created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
        ON CONFLICT (name, city_id) DO UPDATE 
        SET 
            address = EXCLUDED.address,
            stars = EXCLUDED.stars,
            rating = EXCLUDED.rating,
            price_per_night = EXCLUDED.price_per_night,
            description = EXCLUDED.description,
            website = COALESCE(EXCLUDED.website,hotels.website),
            updated_at = NOW() AT TIME ZONE 'Asia/Yerevan'
        RETURNING hotel_id`

	if hotel.CreatedAt.IsZero() {
		hotel.CreatedAt = time.Now()
	}

	if hotel.UpdatedAt.IsZero() {
		hotel.UpdatedAt = time.Now()
	}

	var hotelID int
	err := r.db.QueryRow(
		query,
		hotel.CityID,
		hotel.Name,
		hotel.Address,
		hotel.Stars,
		hotel.Rating,
		hotel.PricePerNight,
		hotel.Website,
		hotel.Description,
		hotel.CreatedAt,
		hotel.UpdatedAt,
	).Scan(&hotelID)

	if err != nil {
		slog.Error("Failed to upsert hotel",
			"hotel_name", hotel.Name,
			"city_id", hotel.CityID,
			"error", err,
		)
		return 0, fmt.Errorf("failed to upsert hotel %s: %w", hotel.Name, err)
	}

	slog.Debug("Hotel upserted successfully", "hotel_id", hotelID, "name", hotel.Name)
	return hotelID, nil
}

func (r *HotelRepository) GetAllHotels() ([]models.Hotel, error) {
	query := `SELECT 
                hotel_id, city_id, name, address, stars, rating, price_per_night, 
                website, description, 
                created_at, updated_at
              FROM hotels;`

	rows, err := r.db.Query(query)
	if err != nil {
		slog.Error("Failed to fetch all hotels", "error", err)
		return nil, fmt.Errorf("failed to fetch all hotels: %w", err)
	}
	defer rows.Close()

	var hotels []models.Hotel
	for rows.Next() {
		var h models.Hotel
		var starsSql sql.NullInt32
		var ratingSql, priceSql sql.NullFloat64
		var websiteSql, descriptionSql sql.NullString

		if err := rows.Scan(
			&h.HotelID, &h.CityID, &h.Name, &h.Address,
			&starsSql, &ratingSql, &priceSql,
			&websiteSql, &descriptionSql, &h.CreatedAt, &h.UpdatedAt,
		); err != nil {
			slog.Warn("Error scanning hotel row", "error", err)
			continue
		}

		h.Stars = int(starsSql.Int32)
		h.Rating = ratingSql.Float64
		h.PricePerNight = priceSql.Float64
		h.Website = websiteSql.String
		h.Description = descriptionSql.String
		hotels = append(hotels, h)
	}

	return hotels, nil
}

func (r *HotelRepository) GetBestHotelByTier(cityID int, budgetMax float64, tier string) (*models.Hotel, error) {
	slog.Info("Searching for best hotel", "city_id", cityID, "budget_max", budgetMax, "tier", tier)

	hotel := &models.Hotel{}
	var orderBy string
	var filter string

	switch tier {
	case "Economy":
		orderBy = "price_per_night ASC"
	case "Balanced":
		orderBy = "rating DESC, price_per_night DESC"
		filter = "AND price_per_night <= $2 * 0.6"
	case "Luxury":
		orderBy = "rating DESC, price_per_night DESC"
	default:
		orderBy = "rating DESC"
	}

	query := fmt.Sprintf(`
    SELECT 
        hotel_id, city_id, name, address, stars, rating, price_per_night, 
        website, description
    FROM hotels 
    WHERE city_id = $1 AND price_per_night <= $2  %s
    ORDER BY %s
    LIMIT 1`, filter, orderBy)

	err := r.db.QueryRow(query, cityID, budgetMax).Scan(
		&hotel.HotelID,
		&hotel.CityID,
		&hotel.Name,
		&hotel.Address,
		&hotel.Stars,
		&hotel.Rating,
		&hotel.PricePerNight,
		&hotel.Website,
		&hotel.Description,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.Debug("No hotels found matching criteria", "city_id", cityID, "tier", tier)
			return nil, nil
		}
		slog.Error("Database error searching for hotel", "error", err, "city_id", cityID)
		return nil, fmt.Errorf("failed to find hotel: %w", err)
	}

	return hotel, nil
}

func (r *HotelRepository) GetVisitedHotels(userID int) ([]models.Hotel, error) {
	query := `
        SELECT DISTINCT 
            h.hotel_id AS hotel_id, 
            h.city_id AS city_id, 
            h.name AS name, 
            h.address AS address, 
            h.stars AS stars, 
            h.rating AS rating, 
            h.price_per_night AS price_per_night
        FROM hotels h
        JOIN itinerary_activities ia ON h.hotel_id = ia.hotel_id
        JOIN trip_itinerary ti ON ia.itinerary_id = ti.itinerary_id
        JOIN trips t ON ti.trip_id = t.trip_id
        WHERE t.user_id = $1 AND t.status ILIKE 'completed'`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		slog.Error("Database error in GetVisitedHotels", "error", err, "user_id", userID)
		return nil, err
	}
	defer rows.Close()

	var hotels []models.Hotel
	for rows.Next() {
		var h models.Hotel
		if err := rows.Scan(
			&h.HotelID,
			&h.CityID,
			&h.Name,
			&h.Address,
			&h.Stars,
			&h.Rating,
			&h.PricePerNight,
		); err != nil {
			slog.Warn("Error scanning visited hotel row", "error", err)
			return nil, err
		}
		hotels = append(hotels, h)
	}
	return hotels, nil
}
