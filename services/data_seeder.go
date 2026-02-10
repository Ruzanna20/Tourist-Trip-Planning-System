package services

import (
	"fmt"
	"log/slog"
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
	slog.Info("Starting country seeding process")

	countriesToSeed, err := s.countryAPIService.FetchAllCountries()
	if err != nil {
		slog.Error("Country API fetch failed", "error", err)
		return fmt.Errorf("country API fetch failed: %w", err)
	}

	slog.Info("Fetched countries from API, starting DB insertion", "count", len(countriesToSeed))

	for _, country := range countriesToSeed {
		_, err := s.countryRepo.Upsert(&country)
		if err != nil {
			slog.Error("Failed to insert country", "name", country.Name, "code", country.Code, "error", err)
			continue
		}
	}
	slog.Info("Country seeding process completed")
	return nil
}

func (s *DataSeeder) SeedCities() error {
	slog.Info("Starting City Seeding process")

	countries, err := s.countryRepo.GetAll()
	if err != nil {
		slog.Error("Failed to get countries from DB", "error", err)
		return fmt.Errorf("failed tp get countries in db:%w", err)
	}

	for _, country := range countries {
		l := slog.With("country", country.Name, "code", country.Code)
		l.Info("Fetching cities for country")

		apiCities, err := s.cityAPIService.FetchCitiesByCountry(country.Code)
		if err != nil {
			l.Error("Error fetching cities", "error", err)
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
				l.Error("Failed to insert city", "city_name", newCity.Name, "error", err)
				continue
			}
		}
		l.Info("Finished cities for country", "cities_count", len(apiCities))
		time.Sleep(4 * time.Second)
	}

	slog.Info("City seeding process completed")
	return nil
}

func (s *DataSeeder) SeedHotels() error {
	slog.Info("Starting Hotels Seeding process")

	cityLocations, err := s.cityRepo.GetAllCityLocations()
	if err != nil {
		slog.Error("Failed to get cities locations", "error", err)
		return fmt.Errorf("failed to get city locations: %w", err)
	}
	if len(cityLocations) == 0 {
		slog.Error("Failed to get cities in db", "error", err)
		return fmt.Errorf("no cities found in database")
	}

	for _, cityLoc := range cityLocations {
		l := slog.With("city", cityLoc.Name, "id", cityLoc.ID)

		if cityLoc.Latitude == 0 || cityLoc.Longitude == 0 {
			l.Warn("Skipping city: Invalid zero coordinates")
			continue
		}

		l.Info("Fetching hotels for city")
		hotels, err := s.hotelAPIService.FetchHotelsByCity(
			cityLoc.ID, cityLoc.Latitude, cityLoc.Longitude,
		)
		if err != nil {
			l.Error("Error fetching hotels", "error", err)
			time.Sleep(2 * time.Second)
			continue
		}

		for _, hotel := range hotels {
			_, err = s.hotelRepo.Upsert(hotel)
			if err != nil {
				l.Error("Failed to insert hotel", "hotel_name", hotel.Name, "error", err)
				continue
			}
		}
		time.Sleep(1500 * time.Millisecond)
	}
	slog.Info("Hotels seeding process completed")
	return nil
}

func (s *DataSeeder) SeedRestaurants() error {
	slog.Info("Starting Restaurants Seeding")

	cityLocations, err := s.cityRepo.GetAllCityLocations()
	if err != nil {
		slog.Error("Failed to get cities locations", "error", err)
		return fmt.Errorf("failed to get city locations: %w", err)
	}
	if len(cityLocations) == 0 {
		slog.Error("Failed to get cities in db", "error", err)
		return fmt.Errorf("no cities found in database")
	}

	for _, cityLoc := range cityLocations {
		l := slog.With("city", cityLoc.Name, "id", cityLoc.ID)

		if cityLoc.Latitude == 0 || cityLoc.Longitude == 0 {
			l.Warn("Skipping city: Invalid zero coordinates")
			continue
		}

		l.Info("Fetching restaurants for city")
		restaurants, err := s.restaurantAPIService.FetchRestaurantsByCity(
			cityLoc.ID, cityLoc.Latitude, cityLoc.Longitude,
		)
		if err != nil {
			l.Error("Error fetching restaurants", "error", err)
			continue
		}

		for _, restaurant := range restaurants {
			_, err := s.restaurantRepo.Upsert(restaurant)
			if err != nil {
				l.Error("Failed to insert restaurant", "hotel_name", restaurant.Name, "error", err)
				continue
			}
		}

		time.Sleep(3 * time.Second)

	}
	slog.Info("Restaurants seeding process completed")
	return nil
}

func (s *DataSeeder) processFlightRoute(fromCityID, toCityID int, fromIata, toIata string) error {
	l := slog.With("from", fromIata, "to", toIata)

	flightOffer, err := s.flightAPIService.FindBestFlightOffer(fromIata, toIata)

	if err != nil {
		l.Error("Flight search API failed", "error", err)
		return fmt.Errorf("flight search failed: %w", err)
	}

	if flightOffer != nil {
		flightOffer.FromCityID = fromCityID
		flightOffer.ToCityID = toCityID

		if _, err := s.flightRepo.Upsert(flightOffer); err != nil {
			l.Error("Critical DB error upserting flight route", "error", err)
			return fmt.Errorf("critical DB error upserting flight route: %w", err)
		}
		l.Debug("Flight route saved to DB", "price", flightOffer.Price)
	} else {
		l.Info("No flight found for route")
	}
	return nil
}

func (s *DataSeeder) SeedFlights() error {
	slog.Info("Starting Flights Seeding process")

	iataMap, err := s.flightAPIService.CityLocationsToIATA()
	if err != nil {
		slog.Error("Failed to map cities to IATA codes", "error", err)
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
	slog.Info("Flight routes check planned", "total_cities", totalCities, "limit_per_city", limit)

	for i := 0; i < totalCities; i++ {
		fromCityID := cityIDs[i]
		fromIata := iataMap[fromCityID]

		if i%50 == 0 {
			slog.Info("Flights seeding progress", "current_index", i, "total", totalCities, "origin_iata", fromIata)
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
				slog.Warn("Route failed", "from", fromIata, "to", toIata, "error", err)
			}
			time.Sleep(3 * time.Second)

			if err := s.processFlightRoute(toCityID, fromCityID, toIata, fromIata); err != nil {
				slog.Warn("Route failed", "from", toIata, "to", fromIata, "error", err)
			}
			time.Sleep(3 * time.Second)
		}
	}

	slog.Info("Ending flights seeding process")
	return nil
}

func (s *DataSeeder) SeedAttractions() error {
	slog.Info("Starting Attraction Seeding process")

	cityLocations, err := s.cityRepo.GetAllCityLocations()
	if err != nil {
		slog.Error("Failed to get cities locations", "error", err)
		return fmt.Errorf("failed to get city locations: %w", err)
	}
	if len(cityLocations) == 0 {
		slog.Error("Failed to get cities in db", "error", err)
		return fmt.Errorf("no cities found in database")
	}

	for _, cityLoc := range cityLocations {
		l := slog.With("city", cityLoc.Name, "id", cityLoc.ID)

		if cityLoc.Latitude == 0 || cityLoc.Longitude == 0 {
			l.Warn("Skipping city: Invalid zero coordinates")
			continue
		}

		l.Info("Fetching attractions for city")
		attractionData, err := s.attractionAPIService.FetchAttractionByCity(cityLoc.ID, cityLoc.Latitude, cityLoc.Longitude)
		if err != nil {
			l.Error("Failed to fetch attractions from API", "error", err)
			time.Sleep(3 * time.Second)
			continue
		}

		for _, attraction := range attractionData {
			attraction.CityID = cityLoc.ID
			_, err := s.attractionRepo.Upsert(&attraction)
			if err != nil {
				l.Error("Failed to insert attraction into DB", "attraction_name", attraction.Name, "error", err)
				continue
			}
		}
		l.Info("Successfully seeded attractions for city", "count", len(attractionData))
		time.Sleep(3 * time.Second)
	}
	slog.Info("Attraction seeding process completed")
	return nil
}
