package models

import "time"

type Country struct {
	CountryID int    `json:"country_id" db:"country_id"`
	Name      string `json:"name" db:"name"`
	Code      string `json:"code" db:"code"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
