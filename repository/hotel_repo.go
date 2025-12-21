package repository

import (
	"database/sql"
	"fmt"
	"log"
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
        currency, phone, website, image_url, description, 
        created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
        ON CONFLICT (name, city_id) DO UPDATE 
        SET 
            address = EXCLUDED.address,
            stars = EXCLUDED.stars,
            rating = EXCLUDED.rating,
            price_per_night = EXCLUDED.price_per_night,
            currency = EXCLUDED.currency,
            description = EXCLUDED.description,
            image_url = EXCLUDED.image_url,
            phone = COALESCE(EXCLUDED.phone, hotels.phone), 
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
		hotel.Currency,
		hotel.Phone,
		hotel.Website,
		hotel.ImageURL,
		hotel.Description,
		hotel.CreatedAt,
		hotel.UpdatedAt,
	).Scan(&hotelID)

	if err != nil {
		return 0, fmt.Errorf("failed to upsert hotel %s: %w", hotel.Name, err)
	}
	return hotelID, nil
}

func (r *HotelRepository) GetAllHotels() ([]models.Hotel, error) {
	query := `SELECT 
                hotel_id, city_id, name, address, stars, rating, price_per_night, 
                currency, phone, website, image_url, description, 
                created_at, updated_at
              FROM hotels;`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all hotels: %w", err)
	}
	defer rows.Close()

	var hotels []models.Hotel
	for rows.Next() {
		var h models.Hotel
		var starsSql sql.NullInt32
		var ratingSql, priceSql sql.NullFloat64
		var currencySql, phoneSql, websiteSql, imageUrlSql, descriptionSql sql.NullString

		if err := rows.Scan(
			&h.HotelID, &h.CityID, &h.Name, &h.Address,
			&starsSql, &ratingSql, &priceSql, &currencySql, &phoneSql,
			&websiteSql, &imageUrlSql, &descriptionSql, &h.CreatedAt, &h.UpdatedAt,
		); err != nil {
			log.Printf("Error scanning hotel row: %v", err)
			continue
		}

		h.Stars = int(starsSql.Int32)
		h.Rating = ratingSql.Float64
		h.PricePerNight = priceSql.Float64
		h.Currency = currencySql.String
		h.Phone = phoneSql.String
		h.Website = websiteSql.String
		h.ImageURL = imageUrlSql.String
		h.Description = descriptionSql.String
		hotels = append(hotels, h)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return hotels, nil
}

func (r *HotelRepository) GetBestHotelByTier(cityID int, budgetMax float64, tier string) (*models.Hotel, error) {
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

	query := fmt.Sprintf(`SELECT 
	hotel_id, city_id, name, address, stars, rating, price_per_night, 
                currency, phone, website, image_url, description
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
		&hotel.Currency,
		&hotel.Phone,
		&hotel.Website,
		&hotel.ImageURL,
		&hotel.Description,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find hotel: %w", err)
	}
	return hotel, nil
}
