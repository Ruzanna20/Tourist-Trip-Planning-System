package handlers

import (
	"encoding/json"
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

func (h *ReviewHandlers) CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "Authentication error: Invalid User ID", http.StatusUnauthorized)
		return
	}

	var req models.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}

	reviewID, err := h.ReviewService.CreateReview(userID, req)
	if err != nil {
		http.Error(w, "Failed to create review.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"review_id": reviewID,
	})

}
