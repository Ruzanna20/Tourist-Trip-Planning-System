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
	Status            string    `json:"status" db:"status"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

type TripPlanRequest struct {
	Name         string  `json:"name"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	ToCityID     int     `json:"destination_city_id"`
	BudgetAmount float64 `json:"total_price"`
}

type TripOption struct {
	Tier             string  `json:"tier"`
	OutBoundFlight   *Flight `json:"outbound_flight"`
	InBoundFlight    *Flight `json:"inbound_flight"`
	Hotel            *Hotel  `json:"hotel"`
	LogisticsBudget  float64 `json:"logistics_budget"`
	ActivitiesBudget float64 `json:"activites_budget"`
	MoreMoney        float64 `json:"more_money"`
	TotalPriceOfTrip float64 `json:"total_price_of_money"`
}
