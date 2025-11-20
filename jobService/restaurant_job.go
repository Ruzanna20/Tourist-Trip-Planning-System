package jobservice

import (
	"log"
	"time"
	"travel-planning/repository"
	"travel-planning/services"
)

type RestaurantJob struct {
	cityRepo             *repository.CityRepository
	restaurantRepo       *repository.RestaurantRepository
	restaurantAPIService *services.RestaurantAPIService
}

func NewRestaurantJob(cityRepo *repository.CityRepository, restaurantRepo *repository.RestaurantRepository, restaurantAPIService *services.RestaurantAPIService) *RestaurantJob {
	return &RestaurantJob{
		cityRepo:             cityRepo,
		restaurantRepo:       restaurantRepo,
		restaurantAPIService: restaurantAPIService,
	}
}

func (job *RestaurantJob) RunJob() {
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

		restaurants, err := job.restaurantAPIService.FetchRestaurantsByCity(
			cityLoc.ID, cityLoc.Latitude, cityLoc.Longitude,
		)
		if err != nil {
			log.Printf("ERROR fetching hotels for %s: %v", cityLoc.Name, err)
			continue
		}

		for _, restaurant := range restaurants {
			_, err := job.restaurantRepo.Upsert(restaurant)
			if err != nil {
				log.Printf("CRITICAL DB ERROR inserting restaurant '%s': %v", restaurant.Name, err)
				continue
			}
		}

		time.Sleep(2 * time.Second)
	}

	log.Println("Restaurant Job Completed.")
}
