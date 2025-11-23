package services

import (
	"fmt"
	"log"
	"time"
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
	attractionAPIService     *AttractionAPIService
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
	attractionAPIService *AttractionAPIService,
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
		attractionAPIService:     attractionAPIService,
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
			time.Sleep(1500 * time.Millisecond)
			continue
		}

		for _, hotel := range hotels {
			_, err := s.hotelRepo.Upsert(hotel)
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
			log.Printf("Error fetching hotels for %s: %v", cityLoc.Name, err)
			time.Sleep(1500 * time.Millisecond)
			continue
		}

		for _, restaurant := range restaurants {
			_, err := s.restaurantRepo.Upsert(restaurant)
			if err != nil {
				log.Printf("Failed to insert hotel %s: %v", restaurant.Name, err)
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

	iataMap, err := s.transportationAPIService.CityLocationsToIATA()
	if err != nil {
		return fmt.Errorf("failed to map cities to IATA codes: %w", err)
	}

	var CityIDs []int
	for cityID, code := range iataMap {
		if code != "" {
			CityIDs = append(CityIDs, cityID)
		}
	}

	for i, FromCityID := range CityIDs {
		for j := i + 1; j < len(CityIDs) && j < i+50; j++ {
			ToCityID := CityIDs[j]

			fromCityID := iataMap[FromCityID]
			toCityID := iataMap[ToCityID]
			if fromCityID == toCityID {
				continue
			}

			flightOffer, err := s.transportationAPIService.FindBestFlightOffer(fromCityID, toCityID)

			if err != nil {
				log.Printf("ERROR flight search %s -> %s: %v", fromCityID, toCityID, err)
				time.Sleep(3 * time.Second)
				continue
			}

			if flightOffer != nil {
				flightOffer.FromCityID = FromCityID
				flightOffer.ToCityID = ToCityID

				if _, err := s.transportationRepo.Upsert(flightOffer); err != nil {
					log.Printf("CRITICAL DB ERROR upserting flight route: %v", err)
				}
			}
			time.Sleep(3 * time.Second)
		}
	}

	log.Println("Ending transportation seeding proccess...")
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
