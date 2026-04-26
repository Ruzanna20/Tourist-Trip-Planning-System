package models

import "time"

type Notification struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	TripID    int       `json:"trip_id" db:"trip_id"`
	TripTitle string    `json:"trip_title" db:"-"`
	Message   string    `json:"message" db:"message"`
	Type      string    `json:"type" db:"type"`
	IsRead    bool      `json:"is_read" db:"is_read"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
