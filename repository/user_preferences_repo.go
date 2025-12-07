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

func (r *UserPreferencesRepository) GetByUserID(userID int) (*models.UserPreferences, error) {
	preferences := models.UserPreferences{}

	query := `SELECT 
                preference_id, user_id, budget_min, budget_max, currency, 
                travel_style, preferred_categories, created_at, updated_at
              FROM user_preferences 
              WHERE user_id = $1`

	var currencySql, travelStyleSql, categoriesSql sql.NullString
	err := r.db.QueryRow(query, userID).Scan(
		&preferences.PreferenceID,
		&preferences.UserID,
		&preferences.BudgetMin,
		&preferences.BudgetMax,
		&currencySql,
		&travelStyleSql,
		&categoriesSql,
		&preferences.CreatedAt,
		&preferences.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to fetch user preferences for user %d:%w", userID, err)
	}

	preferences.Currency = currencySql.String
	preferences.TravelStyle = travelStyleSql.String
	preferences.PreferredCategories = categoriesSql.String

	return &preferences, nil
}
