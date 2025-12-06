package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"travel-planning/models"
)

func (h *AppHandlers) SetPreferencesHandler(w http.ResponseWriter, r *http.Request) {
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

	var req models.UserPreferences
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}

	if req.BudgetMin <= 0 || req.BudgetMax <= req.BudgetMin {
		http.Error(w, "Invalid budget ranges", http.StatusBadRequest)
		return
	}

	req.UserID = userID
	prefID, err := h.UserPreferencesRepo.Upsert(&req)
	if err != nil {
		log.Printf("DB Error: %v", err)
		http.Error(w, "Failed to save user preferences", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "User preferences saved successfully",
		"preference_id": prefID,
	})
}
