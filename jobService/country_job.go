package jobservice

import (
	"log"
	"travel-planning/services"
)

type CountryJob struct {
	seeder *services.DataSeeder
}

func NewCountryJob(seeder *services.DataSeeder) *CountryJob {
	return &CountryJob{
		seeder: seeder,
	}
}

func (job *CountryJob) RunJob() {
	log.Println("Starting Country Job")

	if err := job.seeder.SeedCountries(); err != nil {
		log.Printf("CRITICAL ERROR during Country Job:%v", err)
	} else {
		log.Println("Country Job completed successfully")
	}
}
