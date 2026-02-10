package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"travel-planning/models"
	"travel-planning/services"
)

type ReviewHandlers struct {
	ReviewService *services.ReviewService
}

func NewReviewHandlers(reviewService *services.ReviewService) *ReviewHandlers {
	return &ReviewHandlers{
		ReviewService: reviewService,
	}
}

// CreateReviewHandler godoc
// @Summary Create a new review
// @Description Submit a review for a hotel, attraction, or restaurant
// @Security BearerAuth
// @Tags Reviews
// @Accept json
// @Produce json
// @Param review body models.CreateReviewRequest true "Review details"
// @Success 201 {object} map[string]interface{} "review_id"
// @Router /api/reviews [post]
func (h *ReviewHandlers) CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	l := slog.With("user_id", userID, "path", r.URL.Path, "method", r.Method)

	if r.Method != http.MethodPost {
		l.Warn("Method not allowed")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err != nil && userID <= 0 {
		l.Error("Unauthorized access attempt: invalid User ID")
		http.Error(w, "Authentication error: Invalid User ID", http.StatusUnauthorized)
		return
	}

	var req models.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		l.Warn("Invalid request body format", "error", err)
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}

	l.Info("Attempting to create review", "entity_type", req.EntityType, "entity_id", req.EntityID)

	reviewID, err := h.ReviewService.CreateReview(userID, req)
	if err != nil {
		l.Error("Failed to create review", "error", err)
		http.Error(w, "Failed to create review.", http.StatusInternalServerError)
		return
	}

	l.Info("Review created successfully", "review_id", reviewID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"review_id": reviewID,
	})

}
