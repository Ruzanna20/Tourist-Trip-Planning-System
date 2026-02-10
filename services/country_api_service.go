package services

import (
	"encoding/json"
	"fmt"
	"log/slog"
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

	slog.Info("Fetching all countries from REST Countries API", "url", apiURL)

	resp, err := s.client.Get(apiURL)
	if err != nil {
		slog.Error("Failed to reach REST Countries API", "error", err)
		return nil, fmt.Errorf("failed to fetch data from REST Countries API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		slog.Warn("REST Countries API returned non-OK status", "status", resp.StatusCode)
		return nil, fmt.Errorf("API request failed: %d", resp.StatusCode)
	}

	var apiCountries []CountryAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiCountries); err != nil {
		slog.Error("Failed to decode country API response", "error", err)
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

	slog.Info("Successfully fetched countries", "count", len(countries))
	return countries, nil
}
