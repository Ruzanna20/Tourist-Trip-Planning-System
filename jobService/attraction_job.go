package jobservice

import (
	"log"
	"travel-planning/services"
)

type AttractionJob struct {
	seeder *services.DataSeeder
}

func NewAttractionJob(seeder *services.DataSeeder) *AttractionJob {
	return &AttractionJob{
		seeder: seeder,
	}
}

func (job *AttractionJob) RunJob() {
	log.Println("Starting Attraction Job")

	if err := job.seeder.SeedAttractions(); err != nil {
		log.Printf("CRITICAL ERROR during Attraction Job:%v", err)
	} else {
		log.Println("Attraction Job completed successfully")
	}
}
