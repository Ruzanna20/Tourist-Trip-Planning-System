package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
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

// CreateTripHandler godoc
// @Summary Creating a new itinerary
// @Security BearerAuth
// @Tags Trips
// @Accept json
// @Produce json
// @Param trip body models.TripPlanRequest true "Travel information"
// @Success 201 {object} map[string]interface{} "trip_id"
// @Router /api/trips/create [post]
func (h *TripHandlers) CreateTripHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	l := slog.With("user_id", userID, "path", r.URL.Path)

	if err != nil || userID <= 0 {
		l.Error("Unauthorized access: invalid or missing X-User-ID")
		http.Error(w, "Authentication error", http.StatusUnauthorized)
		return
	}

	var req models.TripPlanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		l.Warn("Invalid request body format", "error", err)
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}

	if req.Name == "" || req.StartDate == "" || req.EndDate == "" {
		l.Warn("Missing required trip fields", "name", req.Name)
		http.Error(w, "Trip name, start date, and end date are required.", http.StatusBadRequest)
		return
	}

	l.Info("Starting trip planning", "trip_name", req.Name)
	tripID, err := h.TripPlanningService.PlanTrip(userID, req)
	if err != nil {
		l.Error("Trip Planning Failed", "error", err)
		http.Error(w, fmt.Sprintf("Failed to process trip plan: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"trip_id": tripID})
}

// GenerateTripOptions godoc
// @Summary Generation of travel options (with budget types)
// @Security BearerAuth
// @Tags Trips
// @Param id path int true "Trip ID"
// @Produce json
// @Success 200 {array} models.TripOption
// @Router /api/trips/{id}/generate-options [post]
func (h *TripHandlers) GenerateTripOptions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tripID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid Trip ID", http.StatusBadRequest)
		return
	}
	l := slog.With("trip_id", tripID)
	l.Info("Generating trip options (Budget tiers)")

	options, err := h.TripPlanningService.GenerateOptions(tripID)
	if err != nil {
		l.Error("Failed to generate trip options", "error", err)
		http.Error(w, "Failed to generate plan", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(options)
}

// GetTripItineraryHandler godoc
// @Summary Get trip itinerary
// @Description Fetch all days of a specific trip's itinerary
// @Security BearerAuth
// @Tags Trips
// @Param id path int true "Trip ID"
// @Produce json
// @Success 200 {array} models.TripItinerary
// @Router /api/trips/{id}/itinerary [get]
func (h *TripHandlers) GetTripItineraryHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	vars := mux.Vars(r)
	tripIDStr := vars["id"]
	tripID, err := strconv.Atoi(tripIDStr)

	l := slog.With("user_id", userID, "path", r.URL.Path)

	if err != nil || userID <= 0 {
		l.Error("Unauthorized access: invalid or missing X-User-ID")
		http.Error(w, "Authentication error", http.StatusUnauthorized)
		return
	}

	if err != nil || tripID <= 0 {
		l.Warn("Invalid Trip ID format", "trip_id_raw", vars["id"])
		http.Error(w, "Invalid Trip ID format", http.StatusBadRequest)
		return
	}

	l.Debug("Fetching itinerary days from DB", "trip_id", tripID)
	itineraryDays, err := h.TripPlanningService.GetItineraryDays(tripID)
	if err != nil {
		l.Error("DB Error fetching itinerary days", "trip_id", tripID, "error", err)
		http.Error(w, "Error fetching itinerary days.", http.StatusInternalServerError)
		return
	}

	if len(itineraryDays) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "No itinerary found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(itineraryDays)
}

// GetActivitiesHandler godoc
// @Summary Get itinerary activities
// @Description Fetch all activities for a specific itinerary day
// @Security BearerAuth
// @Tags Trips
// @Param id path int true "Itinerary ID"
// @Produce json
// @Success 200 {array} models.ItineraryActivity
// @Router /api/itineraries/{id}/activities [get]
func (h *TripHandlers) GetActivitiesHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, _ := strconv.Atoi(userIDStr)
	vars := mux.Vars(r)
	itineraryIDStr := vars["id"]
	itineraryID, err := strconv.Atoi(itineraryIDStr)

	l := slog.With("user_id", userID, "path", r.URL.Path)

	if err != nil || itineraryID <= 0 {
		l.Warn("Invalid Itinerary ID format", "itinerary_id_raw", vars["id"])
		http.Error(w, "Invalid Itinerary ID format", http.StatusBadRequest)
		return
	}

	l.Debug("Fetching activities for itinerary day", "itinerary_id", itineraryID)
	activities, err := h.TripPlanningService.GetActivitiesByDay(itineraryID)
	if err != nil {
		l.Error("DB Error fetching activities", "itinerary_id", itineraryID, "error", err)
		http.Error(w, "Error fetching itinerary days.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activities)
}

// SelectTripOption godoc
// @Summary Finalize trip selection
// @Description Confirm the chosen travel tier and logistics IDs to finalize the plan
// @Security BearerAuth
// @Tags Trips
// @Param id path int true "Trip ID"
// @Param selection body object true "Selected tier and entity IDs"
// @Success 200 {string} string "Trip finalized successfully"
// @Router /api/trips/{id}/select-option [post]
func (h *TripHandlers) SelectTripOption(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, _ := strconv.Atoi(userIDStr)
	vars := mux.Vars(r)
	tripID, _ := strconv.Atoi(vars["id"])
	l := slog.With("user_id", userID, "trip_id", tripID)

	var req struct {
		Tier             string `json:"tier"`
		HotelID          int    `json:"hotel_id"`
		OutboundFlightID int    `json:"outbound_flight_id"`
		InboundFlightID  int    `json:"inbound_flight_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		l.Warn("Invalid selection body", "error", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	l.Info("Finalizing trip selection", "tier", req.Tier)
	err := h.TripPlanningService.FinalizeTripPlan(
		tripID,
		req.Tier,
		req.HotelID,
		req.OutboundFlightID,
		req.InboundFlightID)
	if err != nil {
		l.Error("Failed to finalize trip plan", "tier", req.Tier, "error", err)
		http.Error(w, "Failed to finalize trip: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// DeleteTripHandler godoc
// @Summary Delete a trip by ID (owner only)
// @Security BearerAuth
// @Tags Trips
// @Param id path int true "Trip ID"
// @Success 204
// @Router /api/trips/{id} [delete]
func (h *TripHandlers) DeleteTripHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	vars := mux.Vars(r)
	tripID, errT := strconv.Atoi(vars["id"])
	l := slog.With("user_id", userID, "trip_id", tripID)

	if err != nil || userID <= 0 {
		l.Error("Unauthorized: invalid X-User-ID")
		http.Error(w, "Authentication error", http.StatusUnauthorized)
		return
	}
	if errT != nil || tripID <= 0 {
		l.Warn("Invalid trip ID")
		http.Error(w, "Invalid trip ID", http.StatusBadRequest)
		return
	}

	if err := h.TripPlanningService.DeleteUserTrip(tripID, userID); err != nil {
		l.Error("Failed to delete trip", "error", err)
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetUserTripsHandler godoc
// @Summary Get all trips for the authenticated user
// @Security BearerAuth
// @Tags Trips
// @Produce json
// @Success 200 {array} models.Trip
// @Router /api/trips [get]
func (h *TripHandlers) GetUserTripsHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	l := slog.With("user_id", userID, "path", r.URL.Path)

	if err != nil || userID <= 0 {
		l.Error("Unauthorized access: invalid or missing X-User-ID")
		http.Error(w, "Authentication error", http.StatusUnauthorized)
		return
	}

	trips, err := h.TripPlanningService.GetUserTrips(userID)
	if err != nil {
		l.Error("Failed to fetch user trips", "error", err)
		http.Error(w, "Error fetching trips", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trips)
}
