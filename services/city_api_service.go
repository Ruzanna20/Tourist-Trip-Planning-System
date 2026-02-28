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
)

const cityAPIURL = "https://overpass-api.de/api/interpreter"
const CitiesLimit = "5"
const MinPopoluation = "[population~\"^[1-9][0-9]{4,}$\"]"

type CityAPIService struct {
	client *http.Client
	cache  *cache.RedisCache
}

func NewCityAPIService(cache *cache.RedisCache) *CityAPIService {
	return &CityAPIService{
		client: &http.Client{Timeout: 45 * time.Second},
		cache:  cache,
	}
}

type CityAPIResponse struct {
	Elements []struct {
		Lat  float64           `json:"lat"`
		Lon  float64           `json:"lon"`
		Tags map[string]string `json:"tags"`
	} `json:"elements"`
}

func (s *CityAPIService) FetchCitiesByCountry(countryCode string) ([]*models.City, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("cities:%s", countryCode)
	l := slog.With("country_code", countryCode)

	var cachedCities []*models.City
	err := s.cache.Get(ctx, cacheKey, &cachedCities)
	if err == nil {
		l.Info("Serving cities from Redis cache", "country_code", countryCode)
		return cachedCities, nil
	}

	l.Info("Fetching cities for country from Overpass API")

	//OverPass QL
	query := fmt.Sprintf(`
		[out:json][timeout:30];
		area["ISO3166-1:alpha2"="%s"]->.country;
		(
	  	node[place=city]%s(area.country);
  		way[place=city]%s(area.country);
 	 	relation[place=city]%s(area.country);
		);
		out center %s;`, countryCode, MinPopoluation, MinPopoluation, MinPopoluation, CitiesLimit)

	data := url.Values{}
	data.Set("data", query)

	startTime := time.Now()

	resp, err := s.client.Post(
		cityAPIURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		l.Error("Overpass City API request failed", "error", err)
		return nil, fmt.Errorf("failed to make overpass API request for %s: %w", countryCode, err)
	}
	defer resp.Body.Close()

	l.Debug("Overpass City response received", "duration", time.Since(startTime))

	if resp.StatusCode != http.StatusOK {
		l.Warn("Overpass City API returned non-OK status", "status", resp.StatusCode)
		return nil, fmt.Errorf("API request failed: %d", resp.StatusCode)
	}

	var apicities CityAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apicities); err != nil {
		l.Error("Failed to decode City JSON response", "error", err)
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	var cities []*models.City
	for _, element := range apicities.Elements {
		name := element.Tags["name:en"]

		populationStr := element.Tags["population"]

		if name == "" || element.Lat == 0 || element.Lon == 0 || populationStr == "" {
			l.Debug("Skipping incomplete city data", "city_name", name)
			continue
		}

		population, _ := strconv.Atoi(populationStr)

		newCity := &models.City{
			Name:        name,
			Latitude:    element.Lat,
			Longitude:   element.Lon,
			Description: fmt.Sprintf("Top city in %s (Pop: %d)", countryCode, population),
		}
		cities = append(cities, newCity)
	}

	if len(cities) > 0 {
		err := s.cache.Set(ctx, cacheKey, cities, 24*time.Hour)
		if err != nil {
			l.Error("Failed to save cities to Redis", "error", err)
		}
	}

	l.Info("Successfully fetched cities", "count", len(cities))
	return cities, nil
}
