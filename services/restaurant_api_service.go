package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
	"travel-planning/models"
)

const restaurantAPIUrl = "https://overpass-api.de/api/interpreter"

type RestaurantAPIService struct {
	client         *http.Client
	searchRadiusKm int
}

func NewRestaurantAPIService() *RestaurantAPIService {
	return &RestaurantAPIService{
		client:         &http.Client{Timeout: 60 * time.Second},
		searchRadiusKm: 10,
	}
}

type RestaurantAPIResponse struct {
	Elements []struct {
		Lat  float64           `json:"lat"`
		Lon  float64           `json:"lon"`
		Tags map[string]string `json:"tags"`
	} `json:"elements"`
}

func (s *RestaurantAPIService) FetchRestaurantsByCity(cityID int, lat, lon float64) ([]*models.Restaurant, error) {
	query := fmt.Sprintf(`
		[out:json][timeout:60];
		(
		  node(around:%d, %.6f, %.6f)["amenity"="restaurant"];
		  node(around:%d, %.6f, %.6f)["amenity"="cafe"];
		  node(around:%d, %.6f, %.6f)["amenity"="bar"];
		);
		out center 50; 
	`, s.searchRadiusKm*1000, lat, lon, s.searchRadiusKm*1000, lat, lon, s.searchRadiusKm*1000, lat, lon)

	data := url.Values{}
	data.Set("data", query)

	resp, err := s.client.Post(
		restaurantAPIUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to make overpass API request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var apirestaurants RestaurantAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apirestaurants); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	var restaurants []*models.Restaurant
	for _, element := range apirestaurants.Elements {
		name := element.Tags["name:en"]
		if name == "" {
			continue
		}
		cuisine := element.Tags["cuisine"]
		phone := element.Tags["phone"]
		website := element.Tags["website"]

		newRestaurant := &models.Restaurant{
			CityID:       cityID,
			Name:         name,
			CuisineType:  strings.TrimSpace(cuisine),
			Address:      fmt.Sprintf("%.6f, %.6f", element.Lat, element.Lon),
			Rating:       1,
			PriceRange:   element.Tags["cuisine:price"],
			Phone:        phone,
			Website:      website,
			ImageURL:     "",
			OpeningHours: "",
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}
		restaurants = append(restaurants, newRestaurant)
	}
	return restaurants, nil
}
