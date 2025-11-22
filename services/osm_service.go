package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"travel-planning/models"
)

const OverpassAPIUrl = "https://overpass-api.de/api/interpreter"
const HotelLimit = "20"

// OSMService-ը կառավարում է Overpass API-ի հետ կապը
type OSMService struct {
	client *http.Client
}

func NewOSMService() *OSMService {
	return &OSMService{
		// 60 վայրկյանի Timeout՝ 504 Timeout-ը կանխելու համար
		client: &http.Client{Timeout: 60 * time.Second},
	}
}

// OverpassResponse - OSM-ի կողմից վերադարձված JSON-ի կառուցվածքն է
type OverpassResponse struct {
	Elements []struct {
		Lat  float64           `json:"lat"`
		Lon  float64           `json:"lon"`
		Tags map[string]string `json:"tags"`
	} `json:"elements"`
}

// FetchHotelsFromOSM - Որոնում է հյուրանոցներ՝ հիմնված կոորդինատների վրա
func (s *OSMService) FetchHotelsFromOSM(lat, lon float64) ([]models.Hotel, error) {

	// 1. Ստեղծել QL հարցումը (Bounding Box 3 կմ)
	delta := 0.03
	lat_S := lat - delta
	lon_W := lon - delta
	lat_N := lat + delta
	lon_E := lon + delta

	// Միայն node-ի որոնում՝ կանխելու 504 Timeout-ը
	query := fmt.Sprintf(`
		[out:json][timeout:60]; 
		(
		  node["amenity"="hotel"](%.6f, %.6f, %.6f, %.6f);
		);
		out center %s;
	`, lat_S, lon_W, lat_N, lon_E, HotelLimit)

	// 2. Պատրաստել POST Body-ն
	bodyContent := fmt.Sprintf("data=%s", url.QueryEscape(query))

	resp, err := s.client.Post(
		OverpassAPIUrl,
		"application/x-www-form-urlencoded",
		strings.NewReader(bodyContent),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to make overpass API request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("overpass api request failed status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// 3. Վերծանել և Մոդելավորել
	var overpassRes OverpassResponse
	if err := json.NewDecoder(resp.Body).Decode(&overpassRes); err != nil {
		return nil, fmt.Errorf("failed to decode overpass response: %w", err)
	}

	var hotels []models.Hotel
	for _, element := range overpassRes.Elements {
		// Ապահովել Անգլերեն Անունը
		name := element.Tags["name:en"]
		if name == "" {
			name = element.Tags["name"]
		}

		if name == "" || element.Lat == 0 || element.Lon == 0 {
			continue
		}

		// Լցնել Amenities, Website, Phone դաշտերը Tags-ով
		amenities := element.Tags["amenity"] + ", " + element.Tags["leisure"]
		website := element.Tags["website"]
		phone := element.Tags["phone"]

		newHotel := models.Hotel{
			Name:      name,
			Address:   element.Tags["addr:full"], // Օգտագործել ամբողջական հասցե, եթե կա
			Amenities: strings.Trim(amenities, ", "),
			Website:   website,
			Phone:     phone,

			// Մնացած դաշտերը (Rating, Price) կլինեն 0 կամ դատարկ
			Stars:         0,
			Rating:        0.0,
			PricePerNight: 0.0,
			CreatedAt:     time.Now().UTC(),
			UpdatedAt:     time.Now().UTC(),
		}
		hotels = append(hotels, newHotel)
	}

	return hotels, nil
}
