package services

import (
	"fmt"
	"log"
	"time"
	"travel-planning/models"
	"travel-planning/repository"
)

type DataSeeder struct {
	countryRepo    *repository.CountryRepository
	cityRepo       *repository.CityRepository
	attractionRepo *repository.AttractionRepository
	hotelRepo      *repository.HotelRepository
	restaurantRepo *repository.RestaurantRepository
	flightRepo     *repository.FlightRepository

	countryAPIService    *CountryAPIService
	cityAPIService       *CityAPIService
	attractionAPIService *AttractionAPIService
	hotelAPIService      *HotelAPIService
	restaurantAPIService *RestaurantAPIService
	flightAPIService     *FlightAPIService
	googleAPIService     *GoogleService
}

func NewDataSeeder(
	countryRepo *repository.CountryRepository,
	cityRepo *repository.CityRepository,
	attractionRepo *repository.AttractionRepository,
	hotelRepo *repository.HotelRepository,
	restaurantRepo *repository.RestaurantRepository,
	flightRepo *repository.FlightRepository,
	countryAPIService *CountryAPIService,
	cityAPIService *CityAPIService,
	attractionAPIService *AttractionAPIService,
	hotelAPIService *HotelAPIService,
	restaurantAPIService *RestaurantAPIService,
	flightAPIService *FlightAPIService,
	googleAPIService *GoogleService,
) *DataSeeder {
	return &DataSeeder{
		countryRepo:          countryRepo,
		cityRepo:             cityRepo,
		attractionRepo:       attractionRepo,
		hotelRepo:            hotelRepo,
		restaurantRepo:       restaurantRepo,
		flightRepo:           flightRepo,
		countryAPIService:    countryAPIService,
		cityAPIService:       cityAPIService,
		attractionAPIService: attractionAPIService,
		hotelAPIService:      hotelAPIService,
		restaurantAPIService: restaurantAPIService,
		flightAPIService:     flightAPIService,
		googleAPIService:     googleAPIService,
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
			time.Sleep(1500 * time.Millisecond)
			continue
		}

		for _, hotel := range hotels {
			_, err = s.hotelRepo.Upsert(hotel)
			if err != nil {
				log.Printf("Failed to insert hotel %s: %v", hotel.Name, err)
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
			log.Printf("Error fetching restaurants for %s: %v", cityLoc.Name, err)
			time.Sleep(1500 * time.Millisecond)
			continue
		}

		for _, restaurant := range restaurants {
			_, err := s.restaurantRepo.Upsert(restaurant)
			if err != nil {
				log.Printf("Failed to insert restaurant %s: %v", restaurant.Name, err)
				continue
			}
		}

		time.Sleep(3 * time.Second)

	}
	log.Println("Ending restaurants seeding proccess...")
	return nil
}

func (s *DataSeeder) processFlightRoute(fromCityID, toCityID int, fromIata, toIata string) error {
	flightOffer, err := s.flightAPIService.FindBestFlightOffer(fromIata, toIata)

	if err != nil {
		return fmt.Errorf("flight search failed: %w", err)
	}

	if flightOffer != nil {
		flightOffer.FromCityID = fromCityID
		flightOffer.ToCityID = toCityID

		if _, err := s.flightRepo.Upsert(flightOffer); err != nil {
			return fmt.Errorf("critical DB error upserting flight route: %w", err)
		}
	} else {
		log.Printf("INFO: No flight found for %s -> %s.", fromIata, toIata)
	}
	return nil
}

func (s *DataSeeder) SeedFlights() error {
	log.Println("Starting Flights Seeding...")

	iataMap, err := s.flightAPIService.CityLocationsToIATA()
	if err != nil {
		return fmt.Errorf("failed to map cities to IATA codes: %w", err)
	}

	var cityIDs []int
	for cityID, code := range iataMap {
		if code != "" {
			cityIDs = append(cityIDs, cityID)
		}
	}

	const limit = 5
	totalCities := len(cityIDs)
	log.Printf("Total possible flight routes to check: %d x %d (max) = %d routes.", totalCities, limit, totalCities*limit)

	for i := 0; i < totalCities; i++ {
		fromCityID := cityIDs[i]
		fromIata := iataMap[fromCityID]

		if i%50 == 0 {
			log.Printf("Progress: Starting search for city #%d of %d (IATA: %s)...", i, totalCities, fromIata) // <<< Աշխատանքի Ընթացք
		}

		maxRange := i + limit
		if maxRange > totalCities {
			maxRange = totalCities
		}

		for j := i + 1; j < maxRange; j++ {
			toCityID := cityIDs[j]
			toIata := iataMap[toCityID]

			if fromIata == toIata {
				continue
			}

			if err := s.processFlightRoute(fromCityID, toCityID, fromIata, toIata); err != nil {
				log.Printf("Route failed (%s -> %s): %v", fromIata, toIata, err)
			}
			time.Sleep(3 * time.Second)

			if err := s.processFlightRoute(toCityID, fromCityID, toIata, fromIata); err != nil {
				log.Printf("Route failed (%s -> %s): %v", toIata, fromIata, err)
			}
			time.Sleep(3 * time.Second)
		}
	}

	log.Println("Ending flights seeding proccess...")
	return nil
}

func (s *DataSeeder) SeedAttractions() error {
	log.Println("Starting Attraction Seeding process...")

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

		attractionData, err := s.attractionAPIService.FetchAttractionByCity(cityLoc.ID, cityLoc.Latitude, cityLoc.Longitude)
		if err != nil {
			log.Printf("ERROR fetching attractions for %s: %v", cityLoc.Name, err)
			time.Sleep(3 * time.Second)
			continue
		}

		for _, attraction := range attractionData {
			attraction.CityID = cityLoc.ID
			_, err := s.attractionRepo.Upsert(&attraction)
			if err != nil {
				log.Printf("Failed to insert attraction %s: %v", attraction.Name, err)
				continue
			}
		}
		time.Sleep(3 * time.Second)
	}
	log.Println("Ending attractions seeding proccess...")
	return nil
}
