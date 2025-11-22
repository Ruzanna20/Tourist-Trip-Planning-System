package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"travel-planning/models"
)

type CountryAPIService struct {
	client *http.Client
}

func NewCountryAPIService() *CountryAPIService {
	return &CountryAPIService{
		client: &http.Client{Timeout: 15 * time.Second},
	}
}

type CountryAPIResponse struct {
	Name struct {
		Common string `json:"common" `
	} `json:"name"`
	Code string `json:"cca2"`
}

func (s *CountryAPIService) FetchAllCountries() ([]models.Country, error) {
	const apiURL = "https://restcountries.com/v3.1/all?fields=name,cca2"
	resp, err := s.client.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from REST Countries API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed: %d", resp.StatusCode)
	}

	var apiCountries []CountryAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiCountries); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	var countries []models.Country
	now := time.Now()
	for _, ac := range apiCountries {
		countries = append(countries, models.Country{
			Name:      ac.Name.Common,
			Code:      ac.Code,
			CreatedAt: now,
		})
	}

	return countries, nil
}
