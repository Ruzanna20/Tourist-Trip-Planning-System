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
	"travel-planning/repository"
)

type TransportationAPIService struct {
	amadeus  *AmadeusService
	cityRepo *repository.CityRepository
}

func NewTransportationAPIService(amadeus *AmadeusService, cityRepo *repository.CityRepository) *TransportationAPIService {
	return &TransportationAPIService{
		amadeus:  amadeus,
		cityRepo: cityRepo,
	}
}

type AirportAPIResponse struct {
	Data []struct {
		IataCode string `json:"iataCode"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}

func (s *TransportationAPIService) FindNearestAirportIataCode(lat, lon float64) (string, error) {
	endpoint := "/v1/reference-data/locations/airports"

	if lat == 0 || lon == 0 {
		return "", fmt.Errorf("invalid zero coordinates")
	}

	params := url.Values{}
	params.Add("latitude", strconv.FormatFloat(lat, 'f', 6, 64))
	params.Add("longitude", strconv.FormatFloat(lon, 'f', 6, 64))
	params.Add("radius", "50")

	resp, err := s.amadeus.ExecuteGetRequest(endpoint, params)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("airport API request failed: %d", resp.StatusCode)
	}

	var apiairports AirportAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiairports); err != nil {
		return "", fmt.Errorf("failed to decode airport search response: %w", err)
	}

	if len(apiairports.Data) > 0 {
		return apiairports.Data[0].IataCode, nil
	}

	return "", nil
}

func (s *TransportationAPIService) CityLocationsToIATA() (map[int]string, error) {
	cityLocations, err := s.cityRepo.GetAllCityLocations()
	if err != nil {
		return nil, fmt.Errorf("failed to get city locations: %w", err)
	}

	iataMap := make(map[int]string)
	for _, city := range cityLocations {
		iataCode, err := s.FindNearestAirportIataCode(city.Latitude, city.Longitude)
		if err != nil {
			log.Printf("Failed to get IATA for %s (ID %d): %v", city.Name, city.ID, err)
			time.Sleep(4 * time.Second)
			continue
		}

		if iataCode != "" {
			iataMap[city.ID] = iataCode
		}

		time.Sleep(4 * time.Second)

	}

	log.Printf("Successfully mapped %d cities to IATA codes.", len(iataMap))
	return iataMap, nil
}

type FlightSearchResponse struct {
	Data []struct {
		Price struct {
			Total string `json:"grandTotal"`
		} `json:"price"`
		Itineraries []struct {
			Segments []struct {
				CarrierCode string `json:"carrierCode"`
				Duration    string `json:"duration"`
			} `json:"segments"`
		} `json:"itineraries"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}

func (s *TransportationAPIService) FindBestFlightOffer(fromCityIata, toCityIata string) (*models.Transportation, error) {
	departureDate := time.Now().Add(30 * 24 * time.Hour).Format("2025-11-22")
	endpoint := "/v2/shopping/flight-offers"

	params := url.Values{}
	params.Add("fromCIty_Code", fromCityIata)
	params.Add("toCity_Code", toCityIata)
	params.Add("departureDate", departureDate)
	params.Add("max", "5")

	resp, err := s.amadeus.ExecuteGetRequest(endpoint, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("flight search failed with status: %d", resp.StatusCode)
	}

	var flightRes FlightSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&flightRes); err != nil {
		return nil, fmt.Errorf("failed to decode flight response: %w", err)
	}

	if len(flightRes.Data) == 0 || len(flightRes.Data[0].Itineraries) == 0 {
		return nil, nil
	}

	bestFlight := flightRes.Data[0]
	totalPrice, _ := strconv.ParseFloat(bestFlight.Price.Total, 64)
	carrier := bestFlight.Itineraries[0].Segments[0].CarrierCode
	durationStr := bestFlight.Itineraries[0].Segments[0].Duration
	duration, _ := time.ParseDuration(strings.ToLower(durationStr[1:]))
	durationMinutes := int(duration.Minutes())

	return &models.Transportation{
		Type:            "FLIGHT",
		Carrier:         carrier,
		DurationMinutes: durationMinutes,
		Price:           totalPrice,
		Currency:        "USD",
		Website:         "Amadeus.com",
	}, nil
}
