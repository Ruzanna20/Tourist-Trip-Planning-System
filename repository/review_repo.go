package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
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
		slog.Error("Failed to insert review",
			"user_id", review.UserID,
			"entity_type", review.EntityType,
			"entity_id", review.EntityID,
			"error", err,
		)
		return 0, fmt.Errorf("failed to insert review for user %d: %w", review.UserID, err)
	}

	slog.Debug("Review inserted successfully", "review_id", reviewID, "user_id", review.UserID)
	return reviewID, nil
}

func (r *ReviewRepository) GetByUserID(userID int) ([]models.Review, error) {
	query := `
        SELECT 
            rv.review_id, rv.rating, rv.comment, rv.created_at,
            CASE 
                WHEN rv.hotel_id IS NOT NULL THEN h.name
                WHEN rv.restaurant_id IS NOT NULL THEN res.name
                WHEN rv.attraction_id IS NOT NULL THEN a.name
                ELSE 'Unknown'
            END as entity_name
        FROM reviews rv
        LEFT JOIN hotels h ON rv.hotel_id = h.hotel_id
        LEFT JOIN restaurants res ON rv.restaurant_id = res.restaurant_id
        LEFT JOIN attractions a ON rv.attraction_id = a.attraction_id
        WHERE rv.user_id = $1`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		slog.Error("Failed to fetch reviews for user", "user_id", userID, "error", err)
		return nil, fmt.Errorf("failed to fetch reviews for user %d: %w", userID, err)
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var rev models.Review
		if err := rows.Scan(
			&rev.ReviewID,
			&rev.Rating,
			&rev.Comment,
			&rev.CreatedAt,
			&rev.EntityName,
		); err != nil {
			slog.Warn("Error scanning review row", "user_id", userID, "error", err)
			continue
		}
		reviews = append(reviews, rev)
	}
	return reviews, nil
}

func (r *ReviewRepository) Delete(reviewID, userID int) error {
	query := `DELETE FROM reviews WHERE review_id = $1 AND user_id = $2`

	res, err := r.db.Exec(query, reviewID, userID)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("review not found or unauthorized")
	}
	return nil
}
