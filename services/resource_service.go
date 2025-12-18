package services

import (
	"log"
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
		log.Printf("DB error fetching country data:%v", err)
		return nil, err
	}
	return countries, nil
}

func (s *ResourceService) GetAllCities() ([]models.City, error) {
	cities, err := s.CityRepo.GetAllCities()
	if err != nil {
		log.Printf("DB error fetching city data:%v", err)
		return nil, err
	}

	return cities, nil

}

func (s *ResourceService) GetAllAttractions() ([]models.Attraction, error) {
	attractions, err := s.AttractionRepo.GetAllAttractions()
	if err != nil {
		log.Printf("DB error fetching attraction data:%v", err)
		return nil, err
	}

	return attractions, nil

}

func (s *ResourceService) GetAllHotels() ([]models.Hotel, error) {
	hotels, err := s.HotelRepo.GetAllHotels()
	if err != nil {
		log.Printf("DB error fetching hotel data:%v", err)
		return nil, err
	}

	return hotels, nil

}

func (s *ResourceService) GetAllRestaurants() ([]models.Restaurant, error) {
	restaurants, err := s.RestaurantRepo.GetAllRestaurants()
	if err != nil {
		log.Printf("DB error fetching restaurant data:%v", err)
		return nil, err
	}

	return restaurants, nil

}

func (s *ResourceService) GetAllFlights() ([]models.Flight, error) {
	flights, err := s.FlightRepo.GetAllFlights()
	if err != nil {
		log.Printf("DB error fetching flight data:%v", err)
		return nil, err
	}

	return flights, nil
}
