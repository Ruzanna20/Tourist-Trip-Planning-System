package jobservice

import (
	"log"
	"travel-planning/services"
)

type FlightJob struct {
	seeder *services.DataSeeder
}

func NewFlightJob(seeder *services.DataSeeder) *FlightJob {
	return &FlightJob{
		seeder: seeder,
	}
}

func (job *FlightJob) RunJob() {
	log.Println("Starting Flight Job")

	if err := job.seeder.SeedFlights(); err != nil {
		log.Printf("CRITICAL ERROR during Flight Job:%v", err)
	} else {
		log.Println("Flight Job completed successfully")
	}
}
