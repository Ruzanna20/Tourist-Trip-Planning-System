package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"travel-planning/models"
	"travel-planning/services"

	"github.com/gorilla/mux"
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

// GetUserReviewsHandler godoc
// @Summary Get all reviews submitted by the authenticated user
// @Security BearerAuth
// @Tags Reviews
// @Produce json
// @Success 200 {array} models.Review
// @Router /api/reviews [get]
func (h *ReviewHandlers) GetUserReviewsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	l := slog.With("user_id", userID, "path", r.URL.Path)

	if err != nil || userID <= 0 {
		l.Error("Unauthorized access attempt: invalid User ID")
		http.Error(w, "Authentication error: Invalid User ID", http.StatusUnauthorized)
		return
	}

	reviews, err := h.ReviewService.GetUserReviews(userID)
	if err != nil {
		l.Error("Failed to fetch user reviews", "error", err)
		http.Error(w, "Error fetching reviews", http.StatusInternalServerError)
		return
	}

	l.Debug("Fetched reviews for user", "count", len(reviews))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

// DeleteReviewHandler godoc
// @Summary Delete a review
// @Security BearerAuth
// @Tags Reviews
// @Param id path int true "Review ID"
// @Success 204 "No Content"
// @Router /api/reviews/{id} [delete]
func (h *ReviewHandlers) DeleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	reviewID, err := strconv.Atoi(vars["id"])
	userIDStr := r.Header.Get("X-User-ID")
	userID, _ := strconv.Atoi(userIDStr)

	if err != nil {
		http.Error(w, "Invalid review ID", http.StatusBadRequest)
		return
	}

	err = h.ReviewService.DeleteReview(reviewID, userID)
	if err != nil {
		slog.Error("Failed to delete review", "review_id", reviewID, "error", err)
		http.Error(w, "Failed to delete review: "+err.Error(), http.StatusInternalServerError)
		return
	}

	slog.Info("Review deleted", "review_id", reviewID, "user_id", userID)
	w.WriteHeader(http.StatusNoContent)
}
