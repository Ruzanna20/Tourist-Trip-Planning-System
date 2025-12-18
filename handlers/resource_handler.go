package handlers

import (
	"encoding/json"
	"log"
	"net/http"
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

func (h *ResourceHandlers) GetAllCountriesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	countries, err := h.ResourceService.GetAllCountries()
	if err != nil {
		log.Printf("Service error:%v", err)
		http.Error(w, "Error fetching countries", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(countries); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (h *ResourceHandlers) GetAllCitiesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cities, err := h.ResourceService.GetAllCities()
	if err != nil {
		log.Printf("Service error:%v", err)
		http.Error(w, "Error fetching cities", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(cities); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

func (h *ResourceHandlers) GetAllAttractionssHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	attractions, err := h.ResourceService.GetAllAttractions()
	if err != nil {
		log.Printf("Service error:%v", err)
		http.Error(w, "Error fetching attractions", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(attractions); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

func (h *ResourceHandlers) GetAllHotelsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	hotels, err := h.ResourceService.GetAllHotels()
	if err != nil {
		log.Printf("Service error:%v", err)
		http.Error(w, "Error fetching hotels", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(hotels); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

func (h *ResourceHandlers) GetAllRestaurantssHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	restaurants, err := h.ResourceService.GetAllRestaurants()
	if err != nil {
		log.Printf("Service error:%v", err)
		http.Error(w, "Error fetching restaurants", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(restaurants); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}

}

func (h *ResourceHandlers) GetAllFlightsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	flights, err := h.ResourceService.GetAllFlights()
	if err != nil {
		log.Printf("Service error:%v", err)
		http.Error(w, "Error fetching flights", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(flights); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
