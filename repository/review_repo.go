package repository

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"travel-planning/models"
)

type ReviewRepository struct {
	db *sql.DB
}

func NewReviewRepository(db *sql.DB) *ReviewRepository {
	return &ReviewRepository{
		db: db,
	}
}

func (r *ReviewRepository) Insert(review *models.Review) (int, error) {
	query := `INSERT INTO reviews (
	    user_id, entity_type, entity_id, rating, comment, review_date, created_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING review_id;`

	var reviewID int
	currTime := time.Now()
	reviewDate := review.ReviewDate
	if reviewDate.IsZero() {
		reviewDate = currTime
	}

	err := r.db.QueryRow(
		query,
		review.UserID,
		review.EntityType,
		review.EntityID,
		review.Rating,
		review.Comment,
		reviewDate,
		currTime,
	).Scan(&reviewID)

	if err != nil {
		log.Printf("DB Error inserting review for user %d: %v", review.UserID, err)
		return 0, fmt.Errorf("failed to insert review: %w", err)
	}
	return reviewID, nil

}
