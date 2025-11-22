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
	countryRepo        *repository.CountryRepository
	cityRepo           *repository.CityRepository
	attractionRepo     *repository.AttractionRepository
	hotelRepo          *repository.HotelRepository
	restaurantRepo     *repository.RestaurantRepository
	transportationRepo *repository.TransportationRepository

	countryAPIService        *CountryAPIService
	cityAPIService           *CityAPIService
	hotelAPIService          *HotelAPIService
	restaurantAPIService     *RestaurantAPIService
	transportationAPIService *TransportationAPIService
}

func NewDataSeeder(
	countryRepo *repository.CountryRepository,
	cityRepo *repository.CityRepository,
	attractionRepo *repository.AttractionRepository,
	hotelRepo *repository.HotelRepository,
	restaurantRepo *repository.RestaurantRepository,
	transportationRepo *repository.TransportationRepository,
	countryAPIService *CountryAPIService,
	cityAPIService *CityAPIService,
	hotelAPIService *HotelAPIService,
	restaurantAPIService *RestaurantAPIService,
	transportationAPIService *TransportationAPIService,
) *DataSeeder {
	return &DataSeeder{
		countryRepo:              countryRepo,
		cityRepo:                 cityRepo,
		attractionRepo:           attractionRepo,
		hotelRepo:                hotelRepo,
		restaurantRepo:           restaurantRepo,
		transportationRepo:       transportationRepo,
		countryAPIService:        countryAPIService,
		cityAPIService:           cityAPIService,
		hotelAPIService:          hotelAPIService,
		restaurantAPIService:     restaurantAPIService,
		transportationAPIService: transportationAPIService,
	}
}

func (s *DataSeeder) SeedCountries() error {
	log.Println("Starting country seeding process...")

	countriesToSeed, err := s.countryAPIService.FetchAllCountries()
	if err != nil {
		return fmt.Errorf("country API fetch failed: %w", err)
	}

	log.Printf("Fetched countries from API.Starting db insertion")

	for _, country := range countriesToSeed {
		lastInsertedID, err := s.countryRepo.Upsert(&country)
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

	for _, country := range countries {
		apiCities, err := s.cityAPIService.FetchCitiesByCountry(country.Code)
		if err != nil {
			log.Printf("Error fetching cities for %s(%s):%v", country.Name, country.Code, err)
			time.Sleep(4 * time.Second)
			continue
		}

		for _, data := range apiCities {
			newCity := &models.City{
				CountryID:   country.CountryID,
				Name:        data.Name,
				Latitude:    data.Latitude,
				Longitude:   data.Longitude,
				Description: data.Description,
			}

			if _, err := s.cityRepo.Upsert(newCity); err != nil {
				log.Printf("Failed to insert city %s: %v", newCity.Name, err)
				continue
			}
		}
		time.Sleep(4 * time.Second)
	}
	log.Println("Ending city seeding proccess...")
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

func (s *DataSeeder) SeedTransportation() error {
	log.Println("Starting Transportation Seeding...")

	iataMap, err := s.transportationAPIService.MapCityLocationsToIATA()
	if err != nil {
		return fmt.Errorf("failed to map cities to IATA codes: %w", err)
	}

	var validCityIDs []int
	for cityID, code := range iataMap {
		if code != "" {
			validCityIDs = append(validCityIDs, cityID)
		}
	}

	totalRoutesProcessed := 0
	sleepInterval := 3 * time.Second

	for i, originID := range validCityIDs {
		for j := i + 1; j < len(validCityIDs) && j < i+50; j++ {
			destinationID := validCityIDs[j]

			originIata := iataMap[originID]
			destinationIata := iataMap[destinationID]
			if originIata == destinationIata {
				continue
			}

			flightOffer, err := s.transportationAPIService.FindBestFlightOffer(originIata, destinationIata)

			if err != nil {
				log.Printf("ERROR flight search %s -> %s: %v", originIata, destinationIata, err)
				time.Sleep(sleepInterval)
				continue
			}

			if flightOffer != nil {
				flightOffer.FromCityID = originID
				flightOffer.ToCityID = destinationID

				if _, err := s.transportationRepo.Upsert(flightOffer); err != nil {
					log.Printf("CRITICAL DB ERROR upserting flight route: %v", err)
				}
				totalRoutesProcessed++
			}
			time.Sleep(sleepInterval)
		}
	}

	log.Printf("Transportation Seeding Completed. Total routes processed: %d", totalRoutesProcessed)
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
