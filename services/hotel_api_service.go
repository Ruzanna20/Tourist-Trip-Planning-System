package services

import (
	"encoding/json"
	"fmt"
	"log"
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
	endpoint := "/v1/reference-data/locations/hotels/by-geocode"

	params := url.Values{}
	params.Add("latitude", strconv.FormatFloat(lat, 'f', 6, 64))
	params.Add("longitude", strconv.FormatFloat(lon, 'f', 6, 64))
	params.Add("radius", fmt.Sprintf("%d", radiusKm))
	params.Add("radiusUnit", "KM")

	resp, err := s.amadeus.ExecuteGetRequest(endpoint, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed: %d", resp.StatusCode)
	}

	var apihotels HotelSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&apihotels); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	if len(apihotels.Errors) > 0 {
		log.Printf("Amadeus API returned errors: %v", apihotels.Errors)
		return nil, fmt.Errorf("amadeus API reported data errors")
	}

	return &apihotels, nil
}

func (s *HotelAPIService) FetchHotelsByCity(cityID int, lat, lon float64) ([]*models.Hotel, error) {
	resp, err := s.SearchHotelsByGeo(lat, lon, s.searchRadiusKm)
	if err != nil {
		return nil, err
	}

	if resp == nil || len(resp.Data) == 0 {
		return nil, nil
	}

	var hotels []*models.Hotel
	for _, hotel := range resp.Data {
		enrichedData, err := s.googleService.EnrichHotelData(hotel.Name, hotel.Address.CityName)
		if err != nil || enrichedData == nil || enrichedData.Rating == 0 || enrichedData.Phone == "" {
			continue
		}

		price := 65.0 + rand.Float64()*100.0
		stars := int(enrichedData.Rating)
		description := ""

		if enrichedData.Description != "" {
			description = enrichedData.Description
		} else {
			continue
		}

		newHotel := &models.Hotel{
			HotelID:       0,
			CityID:        cityID,
			Name:          hotel.Name,
			Address:       strings.Join(hotel.Address.Lines, ", ") + ", " + hotel.Address.CityName,
			Stars:         stars,
			Rating:        enrichedData.Rating,
			PricePerNight: price,
			Phone:         enrichedData.Phone,
			Website:       enrichedData.Website,
			ImageURL:      enrichedData.Photo,
			Description:   description,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		hotels = append(hotels, newHotel)
	}
	return hotels, nil
}
