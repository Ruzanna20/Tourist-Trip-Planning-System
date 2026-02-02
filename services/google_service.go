package services

import (
	"encoding/json"
	"fmt"
	"log/slog"
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
	l := slog.With("hotel_name", name, "city", cityName)

	if s.apiKey == "" {
		l.Error("Google API key missing in environment variables")
		return nil, fmt.Errorf("google api key is not in .env")
	}

	l.Debug("Searching place ID via Google Text Search")

	query := url.QueryEscape(name + " in " + cityName)
	endpoint := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/place/textsearch/json?query=%s&key=%s",
		query, s.apiKey,
	)

	resp, err := s.client.Get(endpoint)
	if err != nil {
		l.Error("Google Text Search HTTP request failed", "error", err)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		l.Warn("Google API returned non-OK status", "status", resp.StatusCode)
		return nil, fmt.Errorf("api request failed with status %d", resp.StatusCode)
	}

	var apiResponse GoogleAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		l.Error("Failed to decode Google Text Search response", "error", err)
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	if apiResponse.Status != "OK" {
		if apiResponse.Status == "ZERO_RESULTS" {
			l.Debug("Google found no results for this hotel")
		} else {
			l.Warn("Google API returned status error", "api_status", apiResponse.Status)
		}
		return nil, nil
	}

	place := apiResponse.Results[0]
	l.Debug("Place ID found, fetching details", "place_id", place.PlaceID)

	detailsURL := fmt.Sprintf(
		"https://maps.googleapis.com/maps/api/place/details/json?place_id=%s&fields=international_phone_number,website,editorial_summary&key=%s",
		place.PlaceID, s.apiKey,
	)

	respDetails, err := s.client.Get(detailsURL)
	if err != nil {
		l.Error("Google Place Details request failed", "place_id", place.PlaceID, "error", err)
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
		Status string `json:"status"`
	}

	if err := json.NewDecoder(respDetails.Body).Decode(&detRes); err != nil {
		l.Error("Failed to decode Google Details response", "place_id", place.PlaceID, "error", err)
		return nil, fmt.Errorf("failed to decode details: %w", err)
	}

	if detRes.Status != "OK" {
		l.Warn("Google Place Details status error", "api_status", detRes.Status)
		return nil, nil
	}

	photoURL := ""
	if len(place.Photos) > 0 {
		photoURL = fmt.Sprintf("https://maps.googleapis.com/maps/api/place/photo?maxwidth=800&photoreference=%s&key=%s",
			place.Photos[0].Reference, s.apiKey)
	}

	l.Info("Successfully enriched hotel data via Google", "has_photo", photoURL != "")

	return &GoogleEnrichmentData{
		Rating:      place.Rating,
		Phone:       detRes.Result.Phone,
		Website:     detRes.Result.Website,
		Photo:       photoURL,
		Description: detRes.Result.EditorialSummary.Overview,
	}, nil
}
