package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	"travel-planning/models"
)

const AttractionAPIUrl = "https://overpass-api.de/api/interpreter"
const AttractionLimit = "20"
const AttractionFilter = "tourism=attraction|museum|monument|viewpoint|zoo|theme_park|gallery|aquarium|historic=*"

type AttractionAPIService struct {
	client         *http.Client
	searchRadiusKm int
}

func NewAttractionAPIService() *AttractionAPIService {
	return &AttractionAPIService{
		client:         &http.Client{Timeout: 90 * time.Second},
		searchRadiusKm: 15,
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
	query := fmt.Sprintf(`
    [out:json][timeout:60];
    (
    	node(around:5000, %.6f, %.6f)[~"tourism|historic|amenity|leisure"~"attraction|museum|monument|castle|theatre|park"];
    	//way(around:5000, %.6f, %.6f)[~"tourism|historic|amenity|leisure"~"attraction|museum|monument|castle|theatre|park"];
    );
    out center %s;`, lat, lon, lat, lon, AttractionLimit)

	data := url.Values{}
	data.Set("data", query)

	resp, err := s.client.Post(
		AttractionAPIUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to make overpass API request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var apiattractions AttractionAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiattractions); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	var attractions []models.Attraction
	for _, element := range apiattractions.Elements {
		name := element.Tags["name:en"]
		if name == "" {
			continue
		}

		category := element.Tags["tourism"]
		if category == "" {
			category = element.Tags["historic"]
		}

		newAttraction := models.Attraction{
			CityID:      cityID,
			Name:        name,
			Category:    category,
			Latitude:    element.Lat,
			Longitude:   element.Lon,
			Rating:      0.0,
			EntryFee:    0.0,
			Website:     element.Tags["website"],
			Description: element.Tags["description"],

			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		attractions = append(attractions, newAttraction)
	}

	log.Printf("Successfully processed %d attractions for city %d.", len(attractions), cityID)
	return attractions, nil
}
