package models

import "time"

type UserPreferences struct {
	PreferenceID        int       `json:"preference_id" db:"preference_id"`
	UserID              int       `json:"user_id" db:"user_id"`
	HomeCityID          int       `json:"home_city_id" db:"home_city_id"`
	BudgetMin           float64   `json:"budget_min" db:"budget_min"`
	BudgetMax           float64   `json:"budget_max" db:"budget_max"`
	TravelStyle         string    `json:"travel_style" db:"travel_style"`
	PreferredCategories string    `json:"preferred_categories" db:"preferred_categories"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}
