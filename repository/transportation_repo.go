package repository

import (
	"database/sql"
	"fmt"
	"time"
	"travel-planning/models"
)

type TransportationRepository struct {
	db *sql.DB
}

func NewTransportationRepository(db *sql.DB) *TransportationRepository {
	return &TransportationRepository{
		db: db,
	}
}

func (r *TransportationRepository) Upsert(transportation *models.Transportation) (int, error) {
	query := `INSERT INTO transportation (
        from_city_id, to_city_id, type, carrier, duration_minutes, price, 
        currency, website, created_at, updated_at
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
    ON CONFLICT (from_city_id, to_city_id, type) DO UPDATE 
    SET 
        carrier = EXCLUDED.carrier,
        duration_minutes = EXCLUDED.duration_minutes,
        price = EXCLUDED.price,
        currency = EXCLUDED.currency,
        website = COALESCE(EXCLUDED.website, transportation.website),
        updated_at = NOW() AT TIME ZONE 'Asia/Yerevan' 
    RETURNING transport_id;`

	if transportation.CreatedAt.IsZero() {
		transportation.CreatedAt = time.Now()
	}

	if transportation.UpdatedAt.IsZero() {
		transportation.UpdatedAt = time.Now()
	}

	var transportationID int
	err := r.db.QueryRow(
		query,
		transportation.FromCityID,
		transportation.ToCityID,
		transportation.Type,
		transportation.Carrier,
		transportation.DurationMinutes,
		transportation.Price,
		transportation.Currency,
		transportation.Website,
		transportation.CreatedAt,
		transportation.UpdatedAt,
	).Scan(&transportationID)

	if err != nil {
		return 0, fmt.Errorf("ERROR upserting transportation between %d and %d: %w", transportation.FromCityID, transportation.ToCityID, err)
	}
	return transportationID, nil
}
