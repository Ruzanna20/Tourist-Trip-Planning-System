package repository

import (
	"database/sql"
	"fmt"
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

func (r *CountryRepository) Insert(country *models.Country) (int, error) {
	query := `INSERT INTO countries (name,code,created_at)
			  VALUES ($1,$2,$3)
			  ON CONFLICT (code) DO UPDATE
			  SET name = EXCLUDED.name,created_at = EXCLUDED.created_at
			  RETURNING country_id`

	var countryID int
	if country.CreatedAt.IsZero() {
		country.CreatedAt = time.Now()
	}

	err := r.db.QueryRow(
		query,
		country.Name,
		country.Code,
		country.CreatedAt,
	).Scan(&countryID)

	if err != nil {
		return 0, fmt.Errorf("failed tp insert country with code %s: %w", country.Code, err)
	}
	return countryID, nil
}

func (r *CountryRepository) GetByCode(code string) (*models.Country,error) {
	query := `SELECT country_id,name,code,created_at
			  FROM countries
			  WHERE code = $1`

	country := &models.Country{}
	err := r.db.QueryRow(query,code).Scan(
		&country.CountryID,
		&country.Name,
		&country.Code,
		&country.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil,nil
		}
		return nil,fmt.Errorf("failed to get country by code %s: %w",code,err)
	}

	return country,nil
}

func (r *CountryRepository) GetAll() ([]models.Country,error) {
	rows,err := r.db.Query("SELECT country_id,name,code,created_at FROM countries")
	if err != nil {
		return nil,fmt.Errorf("failed to execute select all countries: %w",err)
	}

	defer rows.Close()

	var countries []models.Country
	for rows.Next() {
		var c models.Country
		if err := rows.Scan(&c.CountryID, &c.Name, &c.Code, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("error scanning country row: %w", err)
		}
		countries = append(countries, c)
	}
	if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows iteration error: %w", err)
		}

	return countries, nil
}

func (r *CountryRepository) GetCountryCodeToIDMap() (map[string]int, error) {
    countries, err := r.GetAll() 
    if err != nil {
        return nil, err
    }

    countryMap := make(map[string]int)
    for _, country := range countries {
        countryMap[country.Code] = country.CountryID 
    }
    return countryMap, nil
}