package services

import (
	"fmt"
	"log/slog"
	"strings"
	"travel-planning/models"
	"travel-planning/repository"
)

type ReviewService struct {
	ReviewRepo *repository.ReviewRepository
}

func NewReviewService(reviewRepo *repository.ReviewRepository) *ReviewService {
	return &ReviewService{
		ReviewRepo: reviewRepo,
	}
}

func (s *ReviewService) CreateReview(userID int, req models.CreateReviewRequest) (int, error) {
	entityType := strings.ToLower(req.EntityType)
	l := slog.With("user_id", userID, "entity_type", entityType, "entity_id", req.EntityID)

	if req.Rating < 1 || req.Rating > 5 {
		l.Warn("Review creation failed: invalid rating", "rating", req.Rating)
		return 0, fmt.Errorf("rating must be >= 1 and  <= 5")
	}

	if entityType != "hotel" && entityType != "attraction" && entityType != "restaurant" {
		l.Warn("Review creation failed: invalid entity type")
		return 0, fmt.Errorf("invalid entity type")
	}

	if req.EntityID <= 0 {
		l.Warn("Review creation failed: invalid entity ID")
		return 0, fmt.Errorf("invalid entity ID")
	}

	newReview := &models.Review{
		UserID:     userID,
		Rating:     req.Rating,
		Comment:    req.Comment,
		EntityType: entityType,
		EntityID:   req.EntityID,
	}

	l.Debug("Attempting to insert review into database")

	reviewID, err := s.ReviewRepo.Insert(newReview)
	if err != nil {
		l.Error("Database error: failed to save review", "error", err)
		return 0, fmt.Errorf("failed to save review: %w", err)
	}

	l.Info("Review created successfully", "review_id", reviewID)
	return reviewID, nil
}

func (s *ReviewService) GetUserReviews(userID int) ([]models.Review, error) {
	reviews, err := s.ReviewRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch reviews: %w", err)
	}
	return reviews, nil
}

func (s *ReviewService) DeleteReview(reviewID, userID int) error {
	return s.ReviewRepo.Delete(reviewID, userID)
}
