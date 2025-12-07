package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"travel-planning/models"
)

func (h *AppHandlers) CreateReviewHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID == 0 {
		http.Error(w, "Authentication error: Invalid User ID", http.StatusUnauthorized)
		return
	}

	var req models.CreateReviewRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}
	EntityType := strings.ToLower(req.EntityType)
	if req.Rating < 1 || req.Rating > 5 || req.EntityID == 0 || (EntityType != "hotel" && EntityType != "attraction" && EntityType != "restaurant") {
		http.Error(w, "Invalid input. Check rating, entity type (hotel/attraction/restaurant), and entity ID.", http.StatusBadRequest)
		return
	}

	newReview := &models.Review{
		UserID:     userID,
		Rating:     req.Rating,
		Comment:    req.Comment,
		EntityType: EntityType,
		EntityID:   req.EntityID,
	}

	reviewID, err := h.ReviewRepo.Insert(newReview)
	if err != nil {
		log.Printf("DB Error inserting review: %v", err)
		http.Error(w, "Failed to create review.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":   fmt.Sprintf("Review created for %s successfully.", EntityType),
		"review_id": reviewID,
	})

}
