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

type FlightAPIService struct {
	amadeus  *AmadeusService
	cityRepo *repository.CityRepository
}

func NewFlightAPIService(amadeus *AmadeusService, cityRepo *repository.CityRepository) *FlightAPIService {
	return &FlightAPIService{
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

func (s *FlightAPIService) FindNearestAirportIataCode(lat, lon float64) (string, error) {
	endpoint := "/v1/reference-data/locations/airports"

	if lat == 0 || lon == 0 {
		return "", fmt.Errorf("invalid zero coordinates")
	}

	params := url.Values{}
	params.Add("latitude", strconv.FormatFloat(lat, 'f', 6, 64))
	params.Add("longitude", strconv.FormatFloat(lon, 'f', 6, 64))
	params.Add("radius", "300")

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

func (s *FlightAPIService) CityLocationsToIATA() (map[int]string, error) {
	cityLocations, err := s.cityRepo.GetAllCityLocations()
	if err != nil {
		return nil, fmt.Errorf("failed to get city locations: %w", err)
	}

	iataMap := make(map[int]string)
	for _, city := range cityLocations {
		if city.IataCode != "" {
			iataMap[city.ID] = city.IataCode
			continue
		}

		iataCode, err := s.FindNearestAirportIataCode(city.Latitude, city.Longitude)
		if err != nil {
			log.Printf("Failed to get IATA for %s (ID %d): %v", city.Name, city.ID, err)
			time.Sleep(4 * time.Second)
			continue
		}

		if iataCode != "" {
			iataMap[city.ID] = iataCode
			if err := s.cityRepo.UpsertCityIata(city.ID, iataCode); err != nil {
				log.Printf("CRITICAL: Failed to save IATA for city %d: %v", city.ID, err)
			}
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

func (s *FlightAPIService) FindBestFlightOffer(fromCityIata, toCityIata string) (*models.Flight, error) {
	departureDate := time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02")
	endpoint := "/v2/shopping/flight-offers"

	params := url.Values{}
	params.Add("originLocationCode", toCityIata)
	params.Add("destinationLocationCode", fromCityIata)
	params.Add("departureDate", departureDate)
	params.Add("adults", "1")
	params.Add("max", "5")
	params.Add("currencyCode", "USD")

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
	durationStr = strings.TrimPrefix(durationStr, "PT")
	duration, _ := time.ParseDuration(strings.ToLower(durationStr[1:]))
	durationMinutes := int(duration.Minutes())

	return &models.Flight{
		Airline:         carrier,
		DurationMinutes: durationMinutes,
		Price:           totalPrice,
		Website:         "Amadeus.com",
	}, nil
}
