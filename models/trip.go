package models

import "time"

type Trip struct {
	TripID            int       `json:"trip_id" db:"trip_id"`
	UserID            int       `json:"user_id" db:"user_id"`
	DestinationCityID int       `json:"destination_city_id" db:"destination_city_id"`
	Title             string    `json:"title" db:"title"`
	StartDate         time.Time `json:"start_date" db:"start_date"`
	EndDate           time.Time `json:"end_date" db:"end_date"`
	TotalPrice        float64   `json:"total_price" db:"total_price"`
	Currency          string    `json:"currency" db:"currency"`
	Status            string    `json:"status" db:"status"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type TripPlanRequest struct {
	Name              string  `json:"name"`
	StartDate         string  `json:"start_date"`
	EndDate           string  `json:"end_date"`
	BudgetAmount      float64 `json:"budget_amount"`
	Currency          string  `json:"currency"`
	DestinationCityID int     `json:"destination_city_id"`
}
