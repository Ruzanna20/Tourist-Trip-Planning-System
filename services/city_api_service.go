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
)

const cityAPIURL = "https://overpass-api.de/api/interpreter"
const CitiesLimit = "5"
const MinPopoluation = "[population~\"^[1-9][0-9]{4,}$\"]"

type CityAPIService struct {
	client *http.Client
}

func NewCityAPIService() *CityAPIService {
	return &CityAPIService{
		client: &http.Client{Timeout: 45 * time.Second},
	}
}

type CityAPIResponse struct {
	Elements []struct {
		Lat  float64           `json:"lat"`
		Lon  float64           `json:"lon"`
		Tags map[string]string `json:"tags"`
	} `json:"elements"`
}

func (s *CityAPIService) FetchCitiesByCountry(countryCode string) ([]models.City, error) {
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

	resp, err := s.client.Post(
		cityAPIURL,
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to make overpass API request for %s: %w", countryCode, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed: %d", resp.StatusCode)
	}
	defer resp.Body.Close()

	var apicities CityAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apicities); err != nil {
		return nil, fmt.Errorf("failed to decode API response: %w", err)
	}

	var cities []models.City
	for _, element := range apicities.Elements {
		name := element.Tags["name:en"]
		if name == "" {
			continue
		}
		populationStr := element.Tags["population"]

		if name == "" || element.Lat == 0 || element.Lon == 0 || populationStr == "" {
			continue
		}

		population, _ := strconv.Atoi(populationStr)

		newCity := models.City{
			Name:        name,
			Latitude:    element.Lat,
			Longitude:   element.Lon,
			Description: fmt.Sprintf("Top city in %s (Pop: %d)", countryCode, population),
		}
		cities = append(cities, newCity)
	}
	log.Printf("Successfully fetched %d cities for %s from OSM.", len(cities), countryCode)
	return cities, nil
}
