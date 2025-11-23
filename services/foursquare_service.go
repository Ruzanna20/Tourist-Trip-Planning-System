package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

type FoursquareService struct {
	client *http.Client
	apiKey string
}

func NewFoursquareService() *FoursquareService {
	return &FoursquareService{
		client: &http.Client{Timeout: 30 * time.Second},
		apiKey: os.Getenv("FOURSQUARE_API_KEY"),
	}
}

type FoursquareEnrichmentData struct {
	Rating    float64
	Website   string
	Phone     string
	Amenities string
}

type FoursquareAPIResponse struct {
	Results []struct {
		FsqId    string `json:"fsq_id"`
		Name     string `json:"name"`
		Website  string `json:"website"`
		Location struct {
			FormattedAddress string `json:"formatted_address"`
		} `json:"location"`
		Rating float64 `json:"rating"`
		Price  int     `json:"price"`
	} `json:"results"`
}

func (s *FoursquareService) EnrichHotelData(name string, lat, lon float64) (*FoursquareEnrichmentData, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("foursquare api key is not .env")
	}

	endpoint := "https://api.foursquare.com/v3/places/search"

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("query", name)
	q.Add("ll", fmt.Sprintf("%.6f,%.6f", lat, lon))
	q.Add("limit", "1")
	q.Add("categories", "19014")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", s.apiKey)

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api request failed with status %d", resp.StatusCode)
	}

	var apiResponse FoursquareAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}
	if len(apiResponse.Results) == 0 {
		return nil, nil
	}

	result := apiResponse.Results[0]
	return &FoursquareEnrichmentData{
		Rating:    result.Rating,
		Website:   result.Website,
		Phone:     "N/A",
		Amenities: "N/A",
	}, nil
}
