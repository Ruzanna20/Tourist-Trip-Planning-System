package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	"travel-planning/models"
)

type CountryRepository struct {
	db *sql.DB
}

func NewCountryRepository(db *sql.DB) *CountryRepository {
	return &CountryRepository{
		db: db,
	}
}

func (r *CountryRepository) Upsert(country *models.Country) (int, error) {
	query := `INSERT INTO countries (name,code,created_at,updated_at)
			  VALUES ($1,$2,$3,$4)
			  ON CONFLICT (code) DO UPDATE
			  SET name = EXCLUDED.name,updated_at = NOW() AT TIME ZONE 'Asia/Yerevan'
			  RETURNING country_id`

	if country.CreatedAt.IsZero() {
		country.CreatedAt = time.Now()
	}

	if country.UpdatedAt.IsZero() {
		country.UpdatedAt = time.Now()
	}

	var countryID int
	err := r.db.QueryRow(
		query,
		country.Name,
		country.Code,
		country.CreatedAt,
		country.UpdatedAt,
	).Scan(&countryID)

	if err != nil {
		slog.Error("Failed to upsert country",
			"country_code", country.Code,
			"country_name", country.Name,
			"error", err,
		)
		return 0, fmt.Errorf("failed to insert country with code %s: %w", country.Code, err)
	}

	slog.Debug("Country upserted successfully", "code", country.Code, "id", countryID)
	return countryID, nil
}

func (r *CountryRepository) GetByCode(code string) (*models.Country, error) {
	query := `SELECT country_id,name,code,created_at,updated_at
			  FROM countries
			  WHERE code = $1`

	country := &models.Country{}
	err := r.db.QueryRow(query, code).Scan(
		&country.CountryID,
		&country.Name,
		&country.Code,
		&country.CreatedAt,
		&country.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			slog.Debug("Country not found", "code", code)
			return nil, nil
		}
		slog.Error("Database error in GetByCode", "code", code, "error", err)
		return nil, fmt.Errorf("failed to get country by code %s: %w", code, err)
	}

	return country, nil
}

func (r *CountryRepository) GetAll() ([]models.Country, error) {
	rows, err := r.db.Query("SELECT country_id,name,code,created_at,updated_at FROM countries")
	if err != nil {
		slog.Error("Failed to fetch all countries", "error", err)
		return nil, fmt.Errorf("failed to execute select all countries: %w", err)
	}

	defer rows.Close()

	var countries []models.Country
	for rows.Next() {
		var c models.Country
		if err := rows.Scan(&c.CountryID, &c.Name, &c.Code, &c.CreatedAt, &c.UpdatedAt); err != nil {
			slog.Warn("Error scanning country row", "error", err)
			return nil, fmt.Errorf("error scanning country row: %w", err)
		}
		countries = append(countries, c)
	}

	if err := rows.Err(); err != nil {
		slog.Error("Rows iteration error in GetAll countries", "error", err)
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	slog.Debug("All countries fetched", "count", len(countries))
	return countries, nil
}
