package repository

import (
	"database/sql"
	"travel-planning/models"
)

type HotelRepository struct {
	db *sql.DB
}

func NewHotelRepository(db *sql.DB) *HotelRepository {
	return &HotelRepository{db: db}
}

func (r *HotelRepository) Insert(hotel *models.Hotel) (int, error) {
	query := `INSERT INTO hotels (
		city_id,name,address,stars,rating,price_per_night,
		currency,amenities,phone,email,website,image_url,description)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING hotel_id`

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
		hotel.Amenities,
		hotel.Phone,
		hotel.Email,
		hotel.Website,
		hotel.ImageURL,
		hotel.Description,
	).Scan(&hotelID)

	if err != nil {
		return 0, err
	}

	return hotelID, nil
}
