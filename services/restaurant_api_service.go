package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
	"travel-planning/internal/cache"
	"travel-planning/models"
)

const restaurantAPIUrl = "https://overpass-api.de/api/interpreter"

type RestaurantAPIService struct {
	client *http.Client
	cache  *cache.RedisCache
}

func NewRestaurantAPIService(cache *cache.RedisCache) *RestaurantAPIService {
	return &RestaurantAPIService{
		client: &http.Client{Timeout: 60 * time.Second},
		cache:  cache,
	}
}

type RestaurantAPIResponse struct {
	Elements []struct {
		Lat  float64           `json:"lat"`
		Lon  float64           `json:"lon"`
		Tags map[string]string `json:"tags"`
	} `json:"elements"`
}

func (s *RestaurantAPIService) FetchRestaurantsByCity(cityID int, lat, lon float64) ([]*models.Restaurant, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("Hotels:%d:%.3f:%.3f", cityID, lat, lon)
	l := slog.With("city_id", cityID, "lat", lat, "lon", lon)

	var cachedRestaurants []*models.Restaurant
	err := s.cache.Get(ctx, cacheKey, &cachedRestaurants)
	if err == nil {
		l.Info("Hotels retrieved from cache")
		return cachedRestaurants, nil
	}

	l.Info("Fetching restaurants from Overpass API")

	searchRadiusM := searchRadiusKm * 1000

	query := fmt.Sprintf(`
		[out:json][timeout:90];
		(
		  node(around:%d, %.6f, %.6f)["amenity"="restaurant"];
		  way["amenity"="restaurant"](around:%d, %.6f, %.6f);
		  node(around:%d, %.6f, %.6f)["amenity"="fast_food"];
		  way["amenity"="fast_food"](around:%d, %.6f, %.6f);
		  node(around:%d, %.6f, %.6f)["amenity"="cafe"];
		  way["amenity"="cafe"](around:%d, %.6f, %.6f);
		  node(around:%d, %.6f, %.6f)["amenity"="bar"];
		  way["amenity"="bar"](around:%d, %.6f, %.6f);
		);
		out center 50; 
	`, searchRadiusM, lat, lon,
		searchRadiusM, lat, lon,
		searchRadiusM, lat, lon,
		searchRadiusM, lat, lon,
		searchRadiusM, lat, lon,
		searchRadiusM, lat, lon,
		searchRadiusM, lat, lon,
		searchRadiusM, lat, lon)

	data := url.Values{}
	data.Set("data", query)
	startTime := time.Now()

	resp, err := s.client.Post(
		restaurantAPIUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		l.Error("Overpass Restaurant API request failed", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	l.Debug("Overpass response received", "duration", time.Since(startTime))

	if resp.StatusCode != http.StatusOK {
		l.Warn("Overpass Restaurant API returned non-OK status", "status", resp.StatusCode)
		return nil, fmt.Errorf("API request failed: %d", resp.StatusCode)
	}

	var apirestaurants RestaurantAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apirestaurants); err != nil {
		l.Error("Failed to decode Restaurant JSON", "error", err)
		return nil, err
	}

	var restaurants []*models.Restaurant
	dif := 5.0 - 1.0
	for _, element := range apirestaurants.Elements {
		name := element.Tags["name:en"]
		if name == "" {
			name = element.Tags["name"]
		}

		cuisine := element.Tags["cuisine"]

		website := element.Tags["website"]
		if website == "" {
			website = element.Tags["contact:website"]
		}

		if name == "" || website == "" {
			continue
		}

		price := element.Tags["cuisine:price"]
		if price == "" {
			prices := []string{"$", "$$", "$$$"}
			price = prices[rand.Intn(len(prices))]
		}

		rating := 1.0 + rand.Float64()*dif

		newRestaurant := &models.Restaurant{
			CityID:     cityID,
			Name:       name,
			Cuisine:    strings.TrimSpace(cuisine),
			Latitude:   element.Lat,
			Longitude:  element.Lon,
			Rating:     rating,
			PriceRange: price,
			Website:    website,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}
		restaurants = append(restaurants, newRestaurant)
	}

	if len(restaurants) > 0 {
		err = s.cache.Set(ctx, cacheKey, restaurants, 24*time.Hour)
		if err != nil {
			l.Error("Failed to save restaurants to cache", "error", err)
		}
	}

	l.Info("Successfully processed restaurants", "total_found", len(apirestaurants.Elements), "added_after_filter", len(restaurants))
	return restaurants, nil
}
