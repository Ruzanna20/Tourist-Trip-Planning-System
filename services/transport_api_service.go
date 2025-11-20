package services

import "travel-planning/repository"

type TransportationAPIService struct {
	amadeus  *AmadeusService
	cityRepo *repository.CityRepository
}

func NewTransportationAPIService(amadeus *AmadeusService, cityRepo *repository.CityRepository) *TransportationAPIService {
	return &TransportationAPIService{
		amadeus:  amadeus,
		cityRepo: cityRepo,
	}
}
