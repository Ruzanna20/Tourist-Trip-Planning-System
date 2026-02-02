package services

import (
	"log/slog"
	"travel-planning/models"
	"travel-planning/repository"
)

type ResourceService struct {
	HotelRepo      *repository.HotelRepository
	CityRepo       *repository.CityRepository
	AttractionRepo *repository.AttractionRepository
	CountryRepo    *repository.CountryRepository
	RestaurantRepo *repository.RestaurantRepository
	FlightRepo     *repository.FlightRepository
}

func NewResourceService(
	HotelRepo *repository.HotelRepository,
	CityRepo *repository.CityRepository,
	AttractionRepo *repository.AttractionRepository,
	CountryRepo *repository.CountryRepository,
	RestaurantRepo *repository.RestaurantRepository,
	FlightRepo *repository.FlightRepository,
) *ResourceService {
	return &ResourceService{
		HotelRepo:      HotelRepo,
		CityRepo:       CityRepo,
		AttractionRepo: AttractionRepo,
		CountryRepo:    CountryRepo,
		RestaurantRepo: RestaurantRepo,
		FlightRepo:     FlightRepo,
	}
}

func (s *ResourceService) GetAllCountries() ([]models.Country, error) {
	countries, err := s.CountryRepo.GetAll()
	if err != nil {
		slog.Error("Database error: failed to fetch countries", "error", err)
		return nil, err
	}
	slog.Debug("Fetched countries from database", "count", len(countries))
	return countries, nil
}

func (s *ResourceService) GetAllCities() ([]models.City, error) {
	cities, err := s.CityRepo.GetAllCities()
	if err != nil {
		slog.Error("Database error: failed to fetch cities", "error", err)
		return nil, err
	}
	slog.Debug("Fetched cities from database", "count", len(cities))
	return cities, nil
}

func (s *ResourceService) GetAllAttractions() ([]models.Attraction, error) {
	attractions, err := s.AttractionRepo.GetAllAttractions()
	if err != nil {
		slog.Error("Database error: failed to fetch attractions", "error", err)
		return nil, err
	}
	slog.Debug("Fetched attractions from database", "count", len(attractions))
	return attractions, nil
}

func (s *ResourceService) GetAllHotels() ([]models.Hotel, error) {
	hotels, err := s.HotelRepo.GetAllHotels()
	if err != nil {
		slog.Error("Database error: failed to fetch hotels", "error", err)
		return nil, err
	}
	slog.Debug("Fetched hotels from database", "count", len(hotels))
	return hotels, nil
}

func (s *ResourceService) GetAllRestaurants() ([]models.Restaurant, error) {
	restaurants, err := s.RestaurantRepo.GetAllRestaurants()
	if err != nil {
		slog.Error("Database error: failed to fetch restaurants", "error", err)
		return nil, err
	}
	slog.Debug("Fetched restaurants from database", "count", len(restaurants))
	return restaurants, nil
}

func (s *ResourceService) GetAllFlights() ([]models.Flight, error) {
	flights, err := s.FlightRepo.GetAllFlights()
	if err != nil {
		slog.Error("Database error: failed to fetch flights", "error", err)
		return nil, err
	}
	slog.Debug("Fetched flights from database", "count", len(flights))
	return flights, nil
}
