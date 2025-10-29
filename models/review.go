package models

import "time"

type Review struct {
	ReviewID   int       `json:"review_id" db:"review_id"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   int       `json:"entity_id" db:"entity_id"`
	AuthorName string    `json:"author_name" db:"author_name"`
	Rating     int       `json:"rating" db:"rating"`
	Comment    string    `json:"comment" db:"comment"`
	ReviewDate time.Time `json:"review_date" db:"review_date"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}
