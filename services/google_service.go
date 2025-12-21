package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"
)

type GoogleService struct {
	client *http.Client
	apiKey string
}

func NewGoogleService() *GoogleService {
	return &GoogleService{
		client: &http.Client{Timeout: 30 * time.Second},
		apiKey: os.Getenv("GOOGLE_API_KEY"),
	}
}

type GoogleEnrichmentData struct {
	Rating      float64
	Website     string
	Phone       string
	Photo       string
	Description string
}

type GoogleAPIResponse struct {
	Results []struct {
		Rating  float64 `json:"rating"`
		PlaceID string  `json:"place_id"`
		Photos  []struct {
			Reference string `json:"photo_reference"`
		} `json:"photos"`
	} `json:"results"`
	Status string `json:"status"`
}

func (s *GoogleService) EnrichHotelData(name string, cityName string) (*GoogleEnrichmentData, error) {
	if s.apiKey == "" {
		return nil, fmt.Errorf("google api key is not .env")
	}

	query := url.QueryEscape(name + " in " + cityName)
	endpoint := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/place/textsearch/json?query=%s&key=%s",
		query, s.apiKey,
	)

	resp, err := s.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("api request failed with status %d", resp.StatusCode)
	}

	var apiResponse GoogleAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	if len(apiResponse.Results) == 0 || apiResponse.Status != "OK" {
		return nil, nil
	}

	place := apiResponse.Results[0]
	detailsURL := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/place/details/json?place_id=%s&fields=international_phone_number,website,editorial_summary&key=%s",
		place.PlaceID, s.apiKey,
	)

	respDetails, err := s.client.Get(detailsURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get details: %w", err)
	}
	defer respDetails.Body.Close()

	var detRes struct {
		Result struct {
			Phone            string `json:"international_phone_number"`
			Website          string `json:"website"`
			EditorialSummary struct {
				Overview string `json:"overview"`
			} `json:"editorial_summary"`
		} `json:"result"`
	}

	if err := json.NewDecoder(respDetails.Body).Decode(&detRes); err != nil {
		return nil, fmt.Errorf("failed to decode details: %w", err)
	}

	photoURL := ""
	if len(place.Photos) > 0 {
		photoURL = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=800&photoreference=%s&key=%s",
			place.Photos[0].Reference, s.apiKey)
	}

	return &GoogleEnrichmentData{
		Rating:      place.Rating,
		Phone:       detRes.Result.Phone,
		Website:     detRes.Result.Website,
		Photo:       photoURL,
		Description: detRes.Result.EditorialSummary.Overview,
	}, nil

}
