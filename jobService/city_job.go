package jobservice

import (
	"log"
	"travel-planning/services"
)

type CityJob struct {
	seeder *services.DataSeeder
}

func NewCityJob(seeder *services.DataSeeder) *CityJob {
	return &CityJob{
		seeder: seeder,
	}
}

func (job *CityJob) RunJob() {
	log.Println("Starting City Job")

	if err := job.seeder.SeedCities(); err != nil {
		log.Printf("CRITICAL ERROR during City Job:%v", err)
	} else {
		log.Println("City Job completed successfully")
	}
}
