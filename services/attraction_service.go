package services

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
	"travel-planning/models"
)

const AttractionAPIUrl = "https://overpass-api.de/api/interpreter"

type AttractionAPIService struct {
	client *http.Client
}

func NewAttractionAPIService() *AttractionAPIService {
	return &AttractionAPIService{
		client: &http.Client{Timeout: 90 * time.Second},
	}
}

type AttractionAPIResponse struct {
	Elements []struct {
		Lat  float64           `json:"lat"`
		Lon  float64           `json:"lon"`
		Tags map[string]string `json:"tags"`
	} `json:"elements"`
}

func (s *AttractionAPIService) FetchAttractionByCity(cityID int, lat, lon float64) ([]models.Attraction, error) {
	l := slog.With("city_id", cityID, "lat", lat, "lon", lon)
	l.Info("Fetching attractions from Overpass API")

	query := fmt.Sprintf(`
    [out:json][timeout:90];
    (
      node["tourism"~"museum|viewpoint|gallery|attraction|monument|historic"](around:10000, %f, %f);
      way["tourism"~"museum|viewpoint|gallery|attraction|monument|historic"](around:10000, %f, %f);
      relation["tourism"~"museum|viewpoint|gallery|attraction|monument|historic"](around:10000, %f, %f);
    );
    out center 50;`, lat, lon, lat, lon, lat, lon)

	data := url.Values{}
	data.Set("data", query)

	startTime := time.Now()

	resp, err := s.client.Post(
		AttractionAPIUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		l.Error("Overpass API request failed", "error", err)
		return nil, fmt.Errorf("failed to make overpass API request: %w", err)
	}
	defer resp.Body.Close()

	l.Debug("Overpass API response received", "duration", time.Since(startTime))

	if resp.StatusCode != http.StatusOK {
		l.Warn("Overpass API returned non-OK status", "status", resp.StatusCode)
		return nil, fmt.Errorf("API request failed: %d", resp.StatusCode)
	}

	var apiattractions AttractionAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiattractions); err != nil {
		l.Error("Failed to decode Overpass JSON", "error", err)
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	var attractions []models.Attraction
	for _, element := range apiattractions.Elements {
		name := element.Tags["name:en"]
		if name == "" {
			name = element.Tags["name"]
		}

		website := element.Tags["website"]
		if website == "" {
			website = element.Tags["contact:website"]
		}

		rating := 3.5 + rand.Float64()*(5.0-3.5)

		tourismType := element.Tags["tourism"]
		var entryFee float64
		switch tourismType {
		case "museum", "gallery", "zoo":
			entryFee = 15.0 + rand.Float64()*(45.0-15.0)
		case "theme_park":
			entryFee = 50.0 + rand.Float64()*(200.0-50.0)
		case "viewpoint", "attraction":
			if rand.Float64() > 0.5 {
				entryFee = 10.0 + rand.Float64()*(25.0-10.0)
			} else {
				entryFee = 0
			}
		case "monument", "artwork", "ruins", "picnic_site":
			entryFee = 0

		default:
			entryFee = 0
		}

		if name == "" || website == "" || tourismType == "" {
			continue
		}

		newAttraction := models.Attraction{
			CityID:    cityID,
			Name:      name,
			Category:  tourismType,
			Latitude:  element.Lat,
			Longitude: element.Lon,
			Rating:    rating,
			EntryFee:  entryFee,
			Website:   website,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		attractions = append(attractions, newAttraction)
	}

	l.Info("Successfully processed attractions", "total_found", len(apiattractions.Elements), "added_to_db", len(attractions))
	return attractions, nil
}
