package services

import (
	"fmt"
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
	if req.Rating < 1 || req.Rating > 5 {
		return 0, fmt.Errorf("rating must be >= 1 and  <= 5")
	}

	if entityType != "hotel" && entityType != "attraction" && entityType != "restaurant" {
		return 0, fmt.Errorf("invalid entity type")
	}

	if req.EntityID <= 0 {
		return 0, fmt.Errorf("invalid entity ID")
	}

	newReview := &models.Review{
		UserID:     userID,
		Rating:     req.Rating,
		Comment:    req.Comment,
		EntityType: entityType,
		EntityID:   req.EntityID,
	}

	reviewID, err := s.ReviewRepo.Insert(newReview)
	if err != nil {
		return 0, fmt.Errorf("failed to save review: %w", err)
	}

	return reviewID, nil
}
