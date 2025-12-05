package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"travel-planning/models"
)

func (h *AppHandlers) CreateTripHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID == 0 {
		log.Printf("CRITICAL: JWT UserID is invalid: %v", err)
		http.Error(w, "Authentication error: Invalid User ID", http.StatusUnauthorized)
		return
	}

	var req models.TripPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.StartDate == "" || req.EndDate == "" {
		http.Error(w, "Trip name, start date, and end date are required.", http.StatusBadRequest)
		return
	}

	tripID, err := h.TripPlanningService.PlanTrip(userID, req)
	if err != nil {
		log.Printf("Trip Planning Failed: %v", err)
		http.Error(w, fmt.Sprintf("Failed to process trip plan: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Trip created and initial planing started",
		"trip_id": tripID,
		"user_id": userID,
	})
}
