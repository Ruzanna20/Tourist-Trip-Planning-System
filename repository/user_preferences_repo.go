package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"travel-planning/models"
)

type UserPreferencesRepository struct {
	db *sql.DB
}

func NewUserPreferencesRepository(db *sql.DB) *UserPreferencesRepository {
	return &UserPreferencesRepository{
		db: db,
	}
}

func (r *UserPreferencesRepository) Upsert(preferences *models.UserPreferences) (int, error) {
	query := `INSERT INTO user_preferences (
        user_id, budget_min, budget_max, currency, travel_style, preferred_categories,created_at, updated_at
    )
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
    ON CONFLICT (user_id) DO UPDATE 
    SET 
        budget_min = EXCLUDED.budget_min,
        budget_max = EXCLUDED.budget_max,
        currency = EXCLUDED.currency,
        travel_style = EXCLUDED.travel_style,
        preferred_categories = EXCLUDED.preferred_categories,
        updated_at = NOW() AT TIME ZONE 'Asia/Yerevan' 
    RETURNING preference_id;`

	var preferenceID int
	currTime := time.Now()

	err := r.db.QueryRow(
		query,
		preferences.UserID,
		preferences.BudgetMin,
		preferences.BudgetMax,
		preferences.Currency,
		preferences.TravelStyle,
		preferences.PreferredCategories,
		currTime,
		currTime,
	).Scan(&preferenceID)

	if err != nil {
		log.Printf("DB Error upserting preferences for user %d: %v", preferences.UserID, err)
		return 0, fmt.Errorf("failed to upsert user preferences: %w", err)
	}
	return preferenceID, nil
}
