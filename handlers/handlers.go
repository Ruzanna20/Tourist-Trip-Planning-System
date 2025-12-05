package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"travel-planning/repository"
	"travel-planning/services"
)

type AppHandlers struct {
	HotelRepo      *repository.HotelRepository
	CityRepo       *repository.CityRepository
	AttractionRepo *repository.AttractionRepository
	CountryRepo    *repository.CountryRepository
	RestaurantRepo *repository.RestaurantRepository
	FlightRepo     *repository.FlightRepository

	TripRepo *repository.TripRepository

	TripPlanningService *services.TripPlanningService
}

func NewAppHandlers(HotelRepo *repository.HotelRepository,
	CityRepo *repository.CityRepository,
	AttractionRepo *repository.AttractionRepository,
	CountryRepo *repository.CountryRepository,
	RestaurantRepo *repository.RestaurantRepository,
	FlightRepo *repository.FlightRepository,
	TripRepo *repository.TripRepository,
	TripPlanningService *services.TripPlanningService,
) *AppHandlers {
	return &AppHandlers{
		HotelRepo:           HotelRepo,
		CityRepo:            CityRepo,
		AttractionRepo:      AttractionRepo,
		CountryRepo:         CountryRepo,
		RestaurantRepo:      RestaurantRepo,
		FlightRepo:          FlightRepo,
		TripRepo:            TripRepo,
		TripPlanningService: TripPlanningService,
	}
}

func (h *AppHandlers) GetAllCountriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	countries, err := h.CountryRepo.GetAll()
	if err != nil {
		log.Printf("DB error fetching data:%v", err)
		http.Error(w, "DB error fetching city data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(countries); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *AppHandlers) GetAllCitiesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cities, err := h.CityRepo.GetAllCities()
	if err != nil {
		log.Printf("DB error fetching data:%v", err)
		http.Error(w, "DB error fetching city data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(cities); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

func (h *AppHandlers) GetAllAttractionssHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	attractions, err := h.AttractionRepo.GetAllAttractions()
	if err != nil {
		log.Printf("DB error fetching data:%v", err)
		http.Error(w, "DB error fetching city data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(attractions); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

func (h *AppHandlers) GetAllHotelsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hotels, err := h.HotelRepo.GetAllHotels()
	if err != nil {
		log.Printf("DB error fetching data:%v", err)
		http.Error(w, "DB error fetching city data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(hotels); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

func (h *AppHandlers) GetAllRestaurantssHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	restaurants, err := h.RestaurantRepo.GetAllRestaurants()
	if err != nil {
		log.Printf("DB error fetching data:%v", err)
		http.Error(w, "DB error fetching city data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(restaurants); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

func (h *AppHandlers) GetAllFlightsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	flights, err := h.FlightRepo.GetAllFlights()
	if err != nil {
		log.Printf("DB error fetching data:%v", err)
		http.Error(w, "DB error fetching flight data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(flights); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *AppHandlers) GetTripsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || userID == 0 {
		http.Error(w, "User authentication required", http.StatusUnauthorized)
		return
	}

	trips, err := h.TripRepo.GetAllTripsByUserID(userID)
	if err != nil {
		log.Printf("DB error fetching trips for user %d: %v", userID, err)
		http.Error(w, "Database error fetching trips", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(trips); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
