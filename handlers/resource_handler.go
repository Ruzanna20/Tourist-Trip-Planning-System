package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"travel-planning/services"
)

type ResourceHandlers struct {
	ResourceService *services.ResourceService
}

func NewResourceHandlers(resourceHandlers *services.ResourceService) *ResourceHandlers {
	return &ResourceHandlers{
		ResourceService: resourceHandlers,
	}
}

// GetAllCountriesHandler godoc
// @Summary Get all countries
// @Tags Resources
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Country
// @Router /api/countries [get]
func (h *ResourceHandlers) GetAllCountriesHandler(w http.ResponseWriter, r *http.Request) {
	l := slog.With("endpoint", "GetAllCountries", "method", r.Method)

	countries, err := h.ResourceService.GetAllCountries()
	if err != nil {
		l.Error("Service error", "error", err)
		http.Error(w, "Error fetching countries", http.StatusInternalServerError)
		return
	}

	l.Debug("Countries fetched successfully", "count", len(countries))
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(countries); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// GetAllCitiesHandler godoc
// @Summary Get all cities
// @Tags Resources
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.City
// @Router /api/cities [get]
func (h *ResourceHandlers) GetAllCitiesHandler(w http.ResponseWriter, r *http.Request) {
	l := slog.With("endpoint", "GetAllCities", "method", r.Method)

	cities, err := h.ResourceService.GetAllCities()
	if err != nil {
		l.Error("Service error", "error", err)
		http.Error(w, "Error fetching cities", http.StatusInternalServerError)
		return
	}

	l.Debug("Cities fetched successfully", "count", len(cities))
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(cities); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

// GetAllAttractionssHandler godoc
// @Summary Get all attractions
// @Tags Resources
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Attraction
// @Router /api/attractions [get]
func (h *ResourceHandlers) GetAllAttractionssHandler(w http.ResponseWriter, r *http.Request) {
	l := slog.With("endpoint", "GetAllAttractions", "method", r.Method)

	attractions, err := h.ResourceService.GetAllAttractions()
	if err != nil {
		l.Error("Service error", "error", err)
		http.Error(w, "Error fetching attractions", http.StatusInternalServerError)
		return
	}

	l.Debug("Attractions fetched successfully", "count", len(attractions))
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(attractions); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

// GetAllHotelsHandler godoc
// @Summary Get all hotels
// @Tags Resources
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Hotel
// @Router /api/hotels [get]
func (h *ResourceHandlers) GetAllHotelsHandler(w http.ResponseWriter, r *http.Request) {
	l := slog.With("endpoint", "GetAllHotels", "method", r.Method)

	hotels, err := h.ResourceService.GetAllHotels()
	if err != nil {
		l.Error("Service error", "error", err)
		http.Error(w, "Error fetching hotels", http.StatusInternalServerError)
		return
	}

	l.Debug("Hotels fetched successfully", "count", len(hotels))
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(hotels); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

// GetAllRestaurantssHandler godoc
// @Summary Get all restaurants
// @Tags Resources
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Restaurant
// @Router /api/restaurants [get]
func (h *ResourceHandlers) GetAllRestaurantssHandler(w http.ResponseWriter, r *http.Request) {
	l := slog.With("endpoint", "GetAllRestaurants", "method", r.Method)

	restaurants, err := h.ResourceService.GetAllRestaurants()
	if err != nil {
		l.Error("Service error", "error", err)
		http.Error(w, "Error fetching restaurants", http.StatusInternalServerError)
		return
	}

	l.Debug("Restaurants fetched successfully", "count", len(restaurants))
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(restaurants); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

// GetAllFlightsHandler godoc
// @Summary Get all flights
// @Tags Resources
// @Security BearerAuth
// @Produce json
// @Success 200 {array} models.Flight
// @Router /api/flights [get]
func (h *ResourceHandlers) GetAllFlightsHandler(w http.ResponseWriter, r *http.Request) {
	l := slog.With("endpoint", "GetAllFlights", "method", r.Method)

	flights, err := h.ResourceService.GetAllFlights()
	if err != nil {
		l.Error("Service error", "error", err)
		http.Error(w, "Error fetching flights", http.StatusInternalServerError)
		return
	}

	l.Debug("Flights fetched successfully", "count", len(flights))
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(flights); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

// GetVisitedEntitiesHandler godoc
// @Summary Get entities (hotels, attractions, restaurants) visited by the user
// @Description Returns a list of entities that appear in the user's completed trips
// @Tags Resources
// @Security BearerAuth
// @Param type query string true "Entity type (hotel, attraction, restaurant)"
// @Produce json
// @Success 200 {array} interface{}
// @Router /api/users/me/visited [get]
func (h *ResourceHandlers) GetVisitedEntitiesHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	slog.Info("Checking Visited for User", "userID", userIDStr)
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID <= 0 {
		slog.Error("Unauthorized: Invalid User ID in GetVisitedEntities", "id_raw", userIDStr)
		http.Error(w, "Unauthorized: Invalid User ID", http.StatusUnauthorized)
		return
	}

	entityType := r.URL.Query().Get("type")
	l := slog.With("endpoint", "GetVisitedEntities", "user_id", userID, "type", entityType)

	if entityType == "" {
		http.Error(w, "Missing type parameter", http.StatusBadRequest)
		return
	}

	l.Info("Fetching visited entities from service")
	entities, err := h.ResourceService.GetVisitedEntities(userID, entityType)
	if err != nil {
		l.Error("Service error", "error", err)
		http.Error(w, "Error fetching visited entities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entities)
}
