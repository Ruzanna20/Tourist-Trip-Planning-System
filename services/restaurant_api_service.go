package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"travel-planning/models"
)

type RestaurantAPIService struct {
	amadeus       *AmadeusService
	searcRadiusKm int
}

func NewRestaurantAPIService(amadeus *AmadeusService) *RestaurantAPIService {
	return &RestaurantAPIService{
		amadeus:       amadeus,
		searcRadiusKm: 10,
	}
}

type RestaurantSearchResponse struct {
	Data []struct {
		Name     string   `json:"name"`
		Category []string `json:"category"`
		GeoCode  struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
		} `json:"geoCode"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}

func (s *RestaurantAPIService) SearchRestaurantsByGeo(lat, lon float64, radiusKm int) (*RestaurantSearchResponse, error) {
	endpoint := "/v1/reference-data/locations/pois"

	params := url.Values{}
	params.Add("latitude", strconv.FormatFloat(lat, 'f', 6, 64))
	params.Add("longitude", strconv.FormatFloat(lon, 'f', 6, 64))
	params.Add("categories", "RESTAURANT")
	params.Add("radius", fmt.Sprintf("%d", radiusKm))
	params.Add("radiusUnit", "KM")
	params.Add("page[limit]", "50")

	resp, err := s.amadeus.ExecuteGetRequest(endpoint, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed: %d", resp.StatusCode)
	}

	var apirestaurants RestaurantSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&apirestaurants); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	if len(apirestaurants.Errors) > 0 {
		log.Printf("Amadeus API returned errors: %v", apirestaurants.Errors)
		return nil, fmt.Errorf("amadeus API reported data errors")
	}

	return &apirestaurants, nil
}

func (s *RestaurantAPIService) FetchRestaurantsByCity(cityID int, lat, lon float64) ([]*models.Restaurant, error) {
	resp, err := s.SearchRestaurantsByGeo(lat, lon, s.searcRadiusKm)
	if err != nil {
		return nil, fmt.Errorf("amadues search failed for city %v:%w", cityID, err)
	}

	var restaurants []*models.Restaurant
	for _, restaurant := range resp.Data {
		newRestaurant := &models.Restaurant{
			CityID:       cityID,
			Name:         restaurant.Name,
			CuisineType:  strings.Join(restaurant.Category, ","),
			Address:      "",
			Rating:       1,
			PriceRange:   "",
			Phone:        "",
			Website:      "",
			ImageURL:     "",
			OpeningHours: "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		restaurants = append(restaurants, newRestaurant)
	}

	return restaurants, nil
}
