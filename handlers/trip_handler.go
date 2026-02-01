package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"travel-planning/models"
	"travel-planning/services"

	"github.com/gorilla/mux"
)

type TripHandlers struct {
	TripPlanningService *services.TripPlanningService
}

func NewTripHandlers(tripPlanningService *services.TripPlanningService) *TripHandlers {
	return &TripHandlers{
		TripPlanningService: tripPlanningService,
	}
}

func (h *TripHandlers) CreateTripHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
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
		"trip_id": tripID,
		"user_id": userID,
	})
}

func (h *TripHandlers) GenerateTripOptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	tripID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Trip ID", http.StatusBadRequest)
		return
	}

	trip, err := h.TripPlanningService.TripRepo.GetTripByID(tripID)
	if err != nil {
		http.Error(w, "Trip not found", http.StatusNotFound)
		return
	}

	options, err := h.TripPlanningService.GenerateOptions(trip)
	if err != nil {
		log.Printf("ERROR in GenerateOptions: %v", err)
		http.Error(w, "Failed to generate plan", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(options); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *TripHandlers) GetTripItineraryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "Authentication error: Invalid User ID", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	tripIDStr := vars["id"]
	tripID, err := strconv.Atoi(tripIDStr)
	if err != nil || tripID <= 0 {
		http.Error(w, "Invalid Trip ID format", http.StatusBadRequest)
		return
	}

	itineraryDays, err := h.TripPlanningService.ItineraryRepo.GetItineraryDaysByTripID(tripID)
	if err != nil {
		log.Printf("DB Error fetching itinerary for trip %d: %v", tripID, err)
		http.Error(w, "Error fetching itinerary days.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(itineraryDays) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No itinerary found for this trip"})
		return
	}
	json.NewEncoder(w).Encode(itineraryDays)
}

func (h *TripHandlers) GetActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		http.Error(w, "Authentication error: Invalid User ID", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	itineraryIDStr := vars["id"]
	itineraryID, err := strconv.Atoi(itineraryIDStr)
	if err != nil || itineraryID <= 0 {
		http.Error(w, "Invalid Itinerary ID format", http.StatusBadRequest)
		return
	}

	activities, err := h.TripPlanningService.ItineraryActivitiesRepo.GetActivitiesByItineraryID(itineraryID)
	if err != nil {
		log.Printf("DB Error fetching activities for itinerary %d: %v", itineraryID, err)
		http.Error(w, "Error fetching itinerary days.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(activities) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No activities found for this day"})
		return
	}

	json.NewEncoder(w).Encode(activities)
}

func (h *TripHandlers) SelectTripOption(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tripID, _ := strconv.Atoi(vars["id"])

	var req struct {
		Tier string `json:"tier"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.TripPlanningService.FinalizeTripPlan(tripID, req.Tier)
	if err != nil {
		http.Error(w, "Failed to finalize trip: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
