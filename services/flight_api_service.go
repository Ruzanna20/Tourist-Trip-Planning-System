package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"travel-planning/internal/cache"
	"travel-planning/models"
	"travel-planning/repository"
)

type FlightAPIService struct {
	amadeus  *AmadeusService
	cityRepo *repository.CityRepository
	cache    *cache.RedisCache
}

func NewFlightAPIService(
	amadeus *AmadeusService,
	cityRepo *repository.CityRepository,
	cache *cache.RedisCache) *FlightAPIService {
	return &FlightAPIService{
		amadeus:  amadeus,
		cityRepo: cityRepo,
		cache:    cache,
	}
}

type AirportAPIResponse struct {
	Data []struct {
		IataCode string `json:"iataCode"`
	} `json:"data"`
	Errors []interface{} `json:"errors"`
}

func (s *FlightAPIService) FindNearestAirportIataCode(lat, lon float64) (string, error) {
	l := slog.With("lat", lat, "lon", lon)
	l.Debug("Searching nearest airport IATA code")

	if lat == 0 || lon == 0 {
		l.Warn("Invalid coordinates provided for airport search")
		return "", fmt.Errorf("invalid zero coordinates")
	}

	endpoint := "/v1/reference-data/locations/airports"
	params := url.Values{}
	params.Add("latitude", strconv.FormatFloat(lat, 'f', 6, 64))
	params.Add("longitude", strconv.FormatFloat(lon, 'f', 6, 64))
	params.Add("radius", "300")

	resp, err := s.amadeus.ExecuteGetRequest(endpoint, params)
	if err != nil {
		l.Error("Amadeus airport search request failed", "error", err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		l.Warn("Airport API returned non-OK status", "status", resp.StatusCode)
		return "", fmt.Errorf("airport API request failed: %d", resp.StatusCode)
	}

	var apiairports AirportAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiairports); err != nil {
		l.Error("Failed to decode airport response", "error", err)
		return "", fmt.Errorf("failed to decode airport search response: %w", err)
	}

	if len(apiairports.Data) > 0 {
		iata := apiairports.Data[0].IataCode
		l.Debug("Found nearest airport", "iata", iata)
		return iata, nil
	}

	l.Warn("No airports found within 300km radius")
	return "", nil
}

func (s *FlightAPIService) CityLocationsToIATA() (map[int]string, error) {
	slog.Info("Starting City to IATA mapping process")

	cityLocations, err := s.cityRepo.GetAllCityLocations()
	if err != nil {
		slog.Error("Failed to fetch city locations from DB", "error", err)
		return nil, fmt.Errorf("failed to get city locations: %w", err)
	}

	iataMap := make(map[int]string)
	for _, city := range cityLocations {
		cl := slog.With("city_name", city.Name, "city_id", city.ID)

		if city.IataCode != "" {
			iataMap[city.ID] = city.IataCode
			continue
		}

		cl.Debug("IATA missing, searching via API")
		iataCode, err := s.FindNearestAirportIataCode(city.Latitude, city.Longitude)
		if err != nil {
			cl.Error("Failed to resolve IATA code", "error", err)
			time.Sleep(4 * time.Second)
			continue
		}

		if iataCode != "" {
			iataMap[city.ID] = iataCode
			if err := s.cityRepo.UpsertCityIata(city.ID, iataCode); err != nil {
				cl.Error("CRITICAL: Failed to save IATA to database", "iata", iataCode, "error", err)
			} else {
				cl.Info("IATA code saved to DB", "iata", iataCode)
			}
		}

		time.Sleep(4 * time.Second)
	}

	slog.Info("City to IATA mapping completed", "mapped_count", len(iataMap))
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
	ctx := context.Background()
	cacheKey := fmt.Sprintf("flights:%s,%s", fromCityIata, toCityIata)
	l := slog.With("from", fromCityIata, "to", toCityIata)

	var cachedFlights *models.Flight
	err := s.cache.Get(ctx, cacheKey, &cachedFlights)
	if err == nil {
		l.Info("Flight offer retrieved from cache")
		return cachedFlights, nil
	}

	l.Info("Searching for best flight offer")

	departureDate := time.Now().Add(30 * 24 * time.Hour).Format("2006-01-02")
	endpoint := "/v2/shopping/flight-offers"

	params := url.Values{}
	params.Add("originLocationCode", fromCityIata)
	params.Add("destinationLocationCode", toCityIata)
	params.Add("departureDate", departureDate)
	params.Add("adults", "1")
	params.Add("max", "5")
	params.Add("currencyCode", "USD")

	resp, err := s.amadeus.ExecuteGetRequest(endpoint, params)
	if err != nil {
		l.Error("Amadeus flight search HTTP error", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		l.Warn("Flight API returned error status", "status", resp.StatusCode)
		return nil, fmt.Errorf("flight search failed with status: %d", resp.StatusCode)
	}

	var flightRes FlightSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&flightRes); err != nil {
		l.Error("Failed to decode flight offers", "error", err)
		return nil, fmt.Errorf("failed to decode flight response: %w", err)
	}

	if len(flightRes.Data) == 0 {
		l.Warn("No flight offers found for this route")
		return nil, nil
	}

	bestFlight := flightRes.Data[0]
	totalPrice, _ := strconv.ParseFloat(bestFlight.Price.Total, 64)

	l.Info("Best flight offer found", "price", totalPrice, "airline", bestFlight.Itineraries[0].Segments[0].CarrierCode)

	durationStr := bestFlight.Itineraries[0].Segments[0].Duration
	durationStr = strings.TrimPrefix(durationStr, "PT")

	duration, _ := time.ParseDuration(strings.ToLower(durationStr))
	durationMinutes := int(duration.Minutes())

	flightResult := &models.Flight{
		Airline:         bestFlight.Itineraries[0].Segments[0].CarrierCode,
		DurationMinutes: durationMinutes,
		Price:           totalPrice,
		Website:         "Amadeus.com",
	}

	err = s.cache.Set(ctx, cacheKey, flightResult, 24*time.Hour)
	if err != nil {
		l.Error("Failed to save flight to cache", "error", err)
	}

	return flightResult, nil

}
