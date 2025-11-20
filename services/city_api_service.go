package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
	"travel-planning/models"
)

type CityAPIService struct {
	client    *http.Client
	geoAPIKey string
	geoHost   string
}

func NewCityAPIService() *CityAPIService {
	return &CityAPIService{
		client:    &http.Client{Timeout: 15 * time.Second},
		geoAPIKey: os.Getenv("GEODB_API_KEY"),
		geoHost:   os.Getenv("GEODB_HOST"),
	}
}

type GeoDBCityResponse struct {
	Data []struct {
		CityID      int     `json:"id"`
		CityName    string  `json:"name"`
		CountryCode string  `json:"countryCode"`
		Latitude    float64 `json:"latitude"`
		Longitude   float64 `json:"longitude"`
		Population  int     `json:"population"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}

func (s *CityAPIService) FetchCitiesByCountry(countryCode string) ([]models.City, error) {
	if s.geoAPIKey == "" || s.geoHost == "" {
		return nil, fmt.Errorf("geodb api  key or host is not in .env")
	}

	endpoint := "https://" + s.geoHost + "/v1/geo/cities"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("countryCode", countryCode)
	q.Add("limit", "5")
	q.Add("sort", "-population")
	q.Add("minPopulation", "100000")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("x-rapidapi-key", s.geoAPIKey)
	req.Header.Add("x-rapidapi-host", s.geoHost)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed geo api request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("geo api request failed: %d", resp.StatusCode)
	}

	var geoResp GeoDBCityResponse
	if err := json.NewDecoder(resp.Body).Decode(&geoResp); err != nil {
		return nil, fmt.Errorf("failed to decode geo api response: %w", err)
	}

	var cities []models.City
	for _, city := range geoResp.Data {
		newCity := models.City{
			Name:        city.CityName,
			Latitude:    city.Latitude,
			Longitude:   city.Longitude,
			Description: fmt.Sprintf("City with %d population", city.Population),
		}
		cities = append(cities, newCity)
	}
	return cities, nil
}
