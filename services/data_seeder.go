package services

import (
	"fmt"
	"log"
	"math"
	"time"
	dataprocessor "travel-planning/data_processor"
	"travel-planning/models"
	"travel-planning/repository"
)

type DataSeeder struct {
	countryRepo    *repository.CountryRepository
	cityRepo       *repository.CityRepository
	attractionRepo *repository.AttractionRepository
	hotelRepo      *repository.HotelRepository
	restaurantRepo *repository.RestaurantRepository

	cityAPIService       *CityAPIService
	hotelAPIService      *HotelAPIService
	restaurantAPIService *RestaurantAPIService
}

func NewDataSeeder(
	countryRepo *repository.CountryRepository,
	cityRepo *repository.CityRepository,
	attractionRepo *repository.AttractionRepository,
	hotelRepo *repository.HotelRepository,
	restaurantRepo *repository.RestaurantRepository,
	cityAPIService *CityAPIService,
	hotelAPIService *HotelAPIService,
	restaurantAPIService *RestaurantAPIService) *DataSeeder {
	return &DataSeeder{
		countryRepo:          countryRepo,
		cityRepo:             cityRepo,
		attractionRepo:       attractionRepo,
		hotelRepo:            hotelRepo,
		restaurantRepo:       restaurantRepo,
		cityAPIService:       cityAPIService,
		hotelAPIService:      hotelAPIService,
		restaurantAPIService: restaurantAPIService,
	}
}

func (ds *DataSeeder) SeedCountries() error {
	log.Println("Starting country seeding process...")

	countriesToSeed, err := FetchAllCountries()
	if err != nil {
		return fmt.Errorf("country API fetch failed: %w", err)
	}

	log.Printf("Fetched countries from API.Starting db insertion")

	for _, country := range countriesToSeed {
		lastInsertedID, err := ds.countryRepo.Upsert(&country)
		if err != nil {
			log.Printf("ERROR.Failed to insert country(%v) %s (%s): %v", lastInsertedID, country.Name, country.Code, err)
			continue
		}
	}
	log.Println("Ending country seeding proccess...")
	return nil
}

func (s *DataSeeder) SeedCities() error {
	log.Println("Starting City Seeding process...")

	countries, err := s.countryRepo.GetAll()
	if err != nil {
		return fmt.Errorf("failed tp get countries in db:%w", err)
	}

	countryMap := make(map[string]int)
	for _, country := range countries {
		countryMap[country.Code] = country.CountryID
	}

	for _, country := range countries {
		apiCities, err := s.cityAPIService.FetchCitiesByCountry(country.Code)
		if err != nil {
			log.Printf("Error fetching cities for %s(%s):%v", country.Name, country.Code, err)
			continue
		}

		for _, data := range apiCities {
			countryID := countryMap[country.Code]

			newCity := &models.City{
				CountryID:   countryID,
				Name:        data.Name,
				Latitude:    data.Latitude,
				Longitude:   data.Longitude,
				Description: "",
			}

			if _, err := s.cityRepo.Upsert(newCity); err != nil {
				log.Printf("failed to upsert city %s", newCity.Name)
				time.Sleep(1500 * time.Millisecond)
				continue
			}

			if _, err := s.cityRepo.Upsert(newCity); err != nil {
				log.Printf("Failed to insert city %s: %v", newCity.Name, err)
				continue
			}
		}
		time.Sleep(1500 * time.Millisecond)

	}
	log.Println("Ending city seeding proccess...")
	return nil
}

const earthRadiusKm = 6371

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	var degToRad = func(deg float64) float64 { return deg * (math.Pi / 180) }

	rLat1 := degToRad(lat1)
	rLon1 := degToRad(lon1)
	rLat2 := degToRad(lat2)
	rLon2 := degToRad(lon2)

	dLat := rLat2 - rLat1
	dLon := rLon2 - rLon1

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(rLat1)*math.Cos(rLat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	//arc length(distance)
	return earthRadiusKm * c
}

func (s *DataSeeder) SeedAttractions(filePath string) error {
	log.Println("Starting Attraction Seeding process...")

	attractionData, err := dataprocessor.FetchAttractionFromEuropeanTour(filePath)
	if err != nil {
		return fmt.Errorf("attraction data processing failed: %w", err)
	}

	cityLocations, err := s.cityRepo.GetAllCityLocations()
	if err != nil {
		return fmt.Errorf("failed to get city locations: %w", err)
	}
	if len(cityLocations) == 0 {
		return fmt.Errorf("no cities found in database")
	}

	for _, data := range attractionData {
		var closestCityID int
		//The largest possible value
		minDistance := float64(math.MaxFloat64)

		for _, city := range cityLocations {
			distance := haversine(data.Latitude, data.Longitude, city.Latitude, city.Longitude)
			if distance < minDistance && distance < 100 {
				minDistance = distance
				closestCityID = city.ID
			}
		}

		if closestCityID == 0 {
			continue
		}

		newAttraction := &models.Attraction{
			CityID:       closestCityID,
			Name:         data.Name,
			Category:     data.Category,
			Latitude:     data.Latitude,
			Longitude:    data.Longitude,
			Rating:       data.Rating,
			EntryFee:     data.EntryFee,
			Currency:     "USD",
			OpeningHours: "",
			Description:  data.Description,
			ImageURL:     "",
			Website:      "",
		}

		_, err := s.attractionRepo.Upsert(newAttraction)
		if err != nil {
			log.Printf("Failed to insert attraction %s: %v", newAttraction.Name, err)
			continue
		}
	}
	log.Println("Ending attractions seeding proccess...")
	return nil
}

func (s *DataSeeder) SeedHotels() error {
	log.Println("Starting Hotels Seeding process...")

	cityLocations, err := s.cityRepo.GetAllCityLocations()
	if err != nil {
		return fmt.Errorf("failed to get city locations: %w", err)
	}
	if len(cityLocations) == 0 {
		return fmt.Errorf("no cities found in database")
	}

	for _, cityLoc := range cityLocations {
		if cityLoc.Latitude == 0 || cityLoc.Longitude == 0 {
			log.Printf("Skipping %s (ID %d): Invalid zero coordinates (Lat:%.4f, Lon:%.4f).",
				cityLoc.Name, cityLoc.ID, cityLoc.Latitude, cityLoc.Longitude)
			continue
		}

		hotels, err := s.hotelAPIService.FetchHotelsByCity(
			cityLoc.ID, cityLoc.Latitude, cityLoc.Longitude,
		)
		if err != nil {
			log.Printf("Error fetching hotels for %s: %v", cityLoc.Name, err)
			continue
		}

		for _, hotel := range hotels {
			_, err := s.hotelRepo.Upsert(hotel)
			if err != nil {
				log.Printf("Failed to insert hotel %s: %v", hotel.Name, err)
				time.Sleep(1500 * time.Millisecond)
				continue
			}
		}
		time.Sleep(1500 * time.Millisecond)

	}
	log.Println("Ending hotels seeding proccess...")
	return nil
}

func (s *DataSeeder) SeedRestaurants() error {
	log.Println("Starting Restaurants Seeding process...")

	cityLocations, err := s.cityRepo.GetAllCityLocations()
	if err != nil {
		return fmt.Errorf("failed to get city locations: %w", err)
	}
	if len(cityLocations) == 0 {
		return fmt.Errorf("no cities found in database")
	}

	for _, cityLoc := range cityLocations {
		if cityLoc.Latitude == 0 || cityLoc.Longitude == 0 {
			log.Printf("Skipping %s (ID %d): Invalid zero coordinates (Lat:%.4f, Lon:%.4f).",
				cityLoc.Name, cityLoc.ID, cityLoc.Latitude, cityLoc.Longitude)
			continue
		}

		restaurants, err := s.restaurantAPIService.FetchRestaurantsByCity(
			cityLoc.ID, cityLoc.Latitude, cityLoc.Longitude,
		)
		if err != nil {
			log.Printf("Error fetching hotels for %s: %v", cityLoc.Name, err)
			continue
		}

		for _, restaurant := range restaurants {
			_, err := s.restaurantRepo.Upsert(restaurant)
			if err != nil {
				log.Printf("Failed to insert hotel %s: %v", restaurant.Name, err)
				time.Sleep(1500 * time.Millisecond)
				continue
			}
		}

		time.Sleep(1500 * time.Millisecond)

	}
	log.Println("Ending restaurants seeding proccess...")
	return nil
}
