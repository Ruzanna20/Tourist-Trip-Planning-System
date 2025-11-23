package jobservice

import (
	"log"
	"time"
	"travel-planning/repository"
	"travel-planning/services"
)

type HotelJob struct {
	cityRepo        *repository.CityRepository
	hotelRepo       *repository.HotelRepository
	hotelAPIService *services.HotelAPIService
}

func NewHotelJob(cityRepo *repository.CityRepository, hotelRepo *repository.HotelRepository, hotelAPIService *services.HotelAPIService) *HotelJob {
	return &HotelJob{
		cityRepo:        cityRepo,
		hotelRepo:       hotelRepo,
		hotelAPIService: hotelAPIService,
	}
}

func (job *HotelJob) RunJob() {
	log.Println("Starting Hotel Job")

	cityLocations, err := job.cityRepo.GetAllCityLocations()
	if err != nil {
		log.Printf("CRITICAL: Failed to get city locations for job: %v\n", err)
		return
	}

	if len(cityLocations) == 0 {
		log.Println("no cities found in db")
		return
	}

	for _, cityLoc := range cityLocations {
		if cityLoc.Latitude == 0 || cityLoc.Longitude == 0 {
			continue
		}

		hotels, err := job.hotelAPIService.FetchHotelsByCity(
			cityLoc.ID, cityLoc.Latitude, cityLoc.Longitude,
		)
		if err != nil {
			log.Printf("ERROR fetching hotels for %s: %v", cityLoc.Name, err)
			time.Sleep(4 * time.Second)
			continue
		}

		for _, hotel := range hotels {
			_, err := job.hotelRepo.Upsert(hotel)
			if err != nil {
				log.Printf("CRITICAL DB ERROR inserting hotel '%s': %v", hotel.Name, err)
				continue
			}
		}

		time.Sleep(4 * time.Second)
	}

	log.Println("Hotel Job Completed.")
}
