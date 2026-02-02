package services

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"travel-planning/models"
)

type HotelAPIService struct {
	amadeus        *AmadeusService
	googleService  *GoogleService
	searchRadiusKm int
}

func NewHotelAPIService(amadeus *AmadeusService, googleService *GoogleService) *HotelAPIService {
	return &HotelAPIService{
		amadeus:        amadeus,
		googleService:  googleService,
		searchRadiusKm: 10,
	}
}

type HotelSearchResponse struct {
	Data []struct {
		HotelID   string  `json:"hotelId"`
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
		Address   struct {
			Lines    []string `json:"lines"`
			CityName string   `json:"cityName"`
		} `json:"address"`
		Price struct {
			Currency string `json:"currency"`
			Total    string `json:"total"`
		} `json:"price"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}

func (s *HotelAPIService) SearchHotelsByGeo(lat, lon float64, radiusKm int) (*HotelSearchResponse, error) {
	l := slog.With("lat", lat, "lon", lon, "radius", radiusKm)
	l.Debug("Searching hotels by geocode via Amadeus")

	endpoint := "/v1/reference-data/locations/hotels/by-geocode"

	params := url.Values{}
	params.Add("latitude", strconv.FormatFloat(lat, 'f', 6, 64))
	params.Add("longitude", strconv.FormatFloat(lon, 'f', 6, 64))
	params.Add("radius", fmt.Sprintf("%d", radiusKm))
	params.Add("radiusUnit", "KM")

	resp, err := s.amadeus.ExecuteGetRequest(endpoint, params)
	if err != nil {
		l.Error("Amadeus request failed", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		l.Warn("Amadeus returned non-OK status", "status", resp.StatusCode)
		return nil, fmt.Errorf("API request failed: %d", resp.StatusCode)
	}

	var apihotels HotelSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&apihotels); err != nil {
		l.Error("Failed to decode Amadeus response", "error", err)
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	if len(apihotels.Errors) > 0 {
		l.Error("Amadeus API reported data errors", "errors", apihotels.Errors)
		return nil, fmt.Errorf("amadeus API reported data errors")
	}

	l.Info("Successfully fetched hotels from Amadeus", "count", len(apihotels.Data))
	return &apihotels, nil
}

func (s *HotelAPIService) FetchHotelsByCity(cityID int, lat, lon float64) ([]*models.Hotel, error) {
	l := slog.With("city_id", cityID)
	l.Info("Starting hotel fetch and enrichment process")

	resp, err := s.SearchHotelsByGeo(lat, lon, s.searchRadiusKm)
	if err != nil {
		return nil, err
	}

	if resp == nil || len(resp.Data) == 0 {
		l.Warn("No hotels found in this area")
		return nil, nil
	}

	var hotels []*models.Hotel
	for _, hotel := range resp.Data {
		hl := l.With("hotel_name", hotel.Name)
		hl.Debug("Enriching hotel data via Google")

		enrichedData, err := s.googleService.EnrichHotelData(hotel.Name, hotel.Address.CityName)
		if err != nil {
			hl.Debug("Google enrichment failed or skipped", "error", err)
			continue
		}

		if enrichedData == nil || enrichedData.Rating == 0 || enrichedData.Phone == "" || enrichedData.Description == "" {
			hl.Debug("Skipping hotel: incomplete enriched data",
				"has_rating", enrichedData != nil && enrichedData.Rating != 0,
				"has_phone", enrichedData != nil && enrichedData.Phone != "",
				"has_desc", enrichedData != nil && enrichedData.Description != "")
			continue
		}

		price := 65.0 + rand.Float64()*100.0
		stars := int(enrichedData.Rating)

		newHotel := &models.Hotel{
			CityID:        cityID,
			Name:          hotel.Name,
			Address:       strings.Join(hotel.Address.Lines, ", ") + ", " + hotel.Address.CityName,
			Stars:         stars,
			Rating:        enrichedData.Rating,
			PricePerNight: price,
			Phone:         enrichedData.Phone,
			Website:       enrichedData.Website,
			ImageURL:      enrichedData.Photo,
			Description:   enrichedData.Description,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		hotels = append(hotels, newHotel)
		hl.Debug("Hotel successfully enriched and added")
	}

	l.Info("Hotel fetch and enrichment completed", "total_added", len(hotels))
	return hotels, nil
}
