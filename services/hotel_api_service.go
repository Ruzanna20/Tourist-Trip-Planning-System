package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand/v2"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"travel-planning/internal/cache"
	"travel-planning/models"
)

const hotelAPIURL = "https://overpass-api.de/api/interpreter"
const hotelLimit = 10
const searchRadiusKm = 20

type HotelAPIService struct {
	client *http.Client
	cache  *cache.RedisCache
}

func NewHotelAPIService(cache *cache.RedisCache) *HotelAPIService {
	return &HotelAPIService{
		client: &http.Client{Timeout: 45 * time.Second},
		cache:  cache,
	}
}

type OverpassHotelResponse struct {
	Elements []struct {
		Lat  float64           `json:"lat"`
		Lon  float64           `json:"lon"`
		Tags map[string]string `json:"tags"`
	} `json:"elements"`
}

func (s *HotelAPIService) FetchHotelsByCity(cityID int, lat, lon float64) ([]*models.Hotel, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("Hotels:%d:%.3f:%.3f", cityID, lat, lon)
	l := slog.With("city_id", cityID, "lat", lat, "lon", lon)

	var cachedHotels []*models.Hotel
	err := s.cache.Get(ctx, cacheKey, &cachedHotels)
	if err == nil {
		l.Info("Hotels retrieved from cache")
		return cachedHotels, nil
	}
	l.Info("Fetching hotels from Overpass API")

	query := fmt.Sprintf(`
		[out:json][timeout:90];
		(
		  node["tourism"="hotel"](around:20000,%[1]f,%[2]f);
		  way["tourism"="hotel"](around:20000,%[1]f,%[2]f);
		);
		out center %d;`, lat, lon, hotelLimit)

	data := url.Values{}
	data.Set("data", query)

	resp, err := s.client.Post(
		hotelAPIURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		l.Error("Overpass request failed", "error", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("overpass API error: %d", resp.StatusCode)
	}

	var osmResp OverpassHotelResponse
	if err := json.NewDecoder(resp.Body).Decode(&osmResp); err != nil {
		return nil, err
	}

	var hotels []*models.Hotel
	dif := 5.0 - 1.0
	for _, el := range osmResp.Elements {
		name := el.Tags["name"]
		if name == "" {
			continue
		}

		stars, _ := strconv.Atoi(el.Tags["stars"])
		if stars == 0 {
			stars = 3 + rand.IntN(3)
		}

		rating := 1.0 + rand.Float64()*dif
		price := 50.0 + rand.Float64()*150.0

		cityName := el.Tags["addr:city"]
		street := el.Tags["addr:street"]
		address := fmt.Sprintf("%s, %s %s", name, street, cityName)
		address = strings.Trim(address, ", ")
		if street == "" && cityName == "" {
			address = name + ", Local Area"
		}

		website := el.Tags["contact:website"]
		if website == "" {
			website = el.Tags["website"]
		}
		if website == "" {
			cleanName := strings.ToLower(strings.ReplaceAll(name, " ", ""))
			website = fmt.Sprintf("https://www.%s.com", cleanName)
		}

		var amenities []string

		if el.Tags["internet_access"] != "" || el.Tags["wifi"] != "" {
			amenities = append(amenities, "Free Wi-Fi")
		}
		if el.Tags["swimming_pool"] == "yes" {
			amenities = append(amenities, "Swimming Pool")
		}
		if el.Tags["parking"] == "yes" {
			amenities = append(amenities, "Private Parking")
		}
		if el.Tags["air_conditioning"] == "yes" {
			amenities = append(amenities, "Air Conditioning")
		}
		if el.Tags["wheelchair"] == "yes" {
			amenities = append(amenities, "Wheelchair Accessible")
		}

		description := el.Tags["description"]
		if description == "" {
			amenitiesText := ""
			if len(amenities) > 0 {
				amenitiesText = " Features: " + strings.Join(amenities, ", ") + "."
			}
			description = fmt.Sprintf("A wonderful stay at %s. %s ", name, amenitiesText)
		}

		newHotel := &models.Hotel{
			CityID:        cityID,
			Name:          name,
			Address:       address,
			Stars:         stars,
			Rating:        rating,
			PricePerNight: price,
			Phone:         el.Tags["contact:phone"],
			Website:       website,
			Description:   description,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}

		hotels = append(hotels, newHotel)
	}

	if len(hotels) > 0 {
		err = s.cache.Set(ctx, cacheKey, hotels, 24*time.Hour)
		if err != nil {
			l.Error("Failed to save hotels to cache", "error", err)
		}
	}

	l.Info("Successfully processed hotels", "added_to_db", len(hotels))
	return hotels, nil
}
