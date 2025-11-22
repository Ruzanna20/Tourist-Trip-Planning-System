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

type AirportSearchResponse struct {
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
	params.Add("page[limit]", "1")

	resp, err := s.amadeus.ExecuteGetRequest(endpoint, params)
	if err != nil {
		return "", fmt.Errorf("amadeus airport search failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("airport search failed with status: %d", resp.StatusCode)
	}

	var airportRes AirportSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&airportRes); err != nil {
		return "", fmt.Errorf("failed to decode airport search response: %w", err)
	}

	if len(airportRes.Data) > 0 {
		return airportRes.Data[0].IataCode, nil
	}

	return "", nil
}

func (s *TransportationAPIService) MapCityLocationsToIATA() (map[int]string, error) {
	cityLocations, err := s.cityRepo.GetAllCityLocations()
	if err != nil {
		return nil, fmt.Errorf("failed to get city locations from repo: %w", err)
	}

	iataMap := make(map[int]string)
	sleepInterval := 3500 * time.Millisecond

	for _, city := range cityLocations {
		iataCode, err := s.FindNearestAirportIataCode(city.Latitude, city.Longitude)
		if err != nil {
			log.Printf("Warning: Failed to get IATA for %s (ID %d): %v", city.Name, city.ID, err)
			time.Sleep(sleepInterval)
			continue
		}

		if iataCode != "" {
			iataMap[city.ID] = iataCode
		}

		time.Sleep(sleepInterval)

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

func (s *TransportationAPIService) FindBestFlightOffer(originIata, destinationIata string) (*models.Transportation, error) {
	departureDate := time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02")
	endpoint := "/v2/shopping/flight-offers"

	params := url.Values{}
	params.Add("originLocationCode", originIata)
	params.Add("destinationLocationCode", destinationIata)
	params.Add("departureDate", departureDate)
	params.Add("adults", "1")
	params.Add("max", "1")

	resp, err := s.amadeus.ExecuteGetRequest(endpoint, params)
	if err != nil {
		return nil, fmt.Errorf("flight search failed: %w", err)
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
