package models

import "time"

const (
	PrefMuseum     = "museum"
	PrefViewpoint  = "viewpoint"
	PrefGallery    = "gallery"
	PrefAttraction = "attraction"
	PrefMonument   = "monument"
	PrefHistoric   = "historic"
)

type UserPreferences struct {
	PreferenceID        int       `json:"preference_id" db:"preference_id"`
	UserID              int       `json:"user_id" db:"user_id"`
	HomeCityID          int       `json:"home_city_id" db:"home_city_id"`
	PreferredCategories string    `json:"preferred_categories" db:"preferred_categories"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}
