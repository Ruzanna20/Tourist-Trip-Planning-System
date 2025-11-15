package jobservice

import (
	"log"
	"travel-planning/repository"
	"travel-planning/services"
)

type CountryJob struct {
	countryRepo *repository.CountryRepository
}

func NewCountryJob(repo *repository.CountryRepository) *CountryJob {
	return &CountryJob{
		countryRepo: repo,
	}
}

func (job *CountryJob) RunJob() {
	log.Println("Starting Country Job")

	apiCountries, err := services.FetchAllCountries()
	if err != nil {
		log.Printf("CRITICAL: Failed to fetch countries from API: %v\n", err)
		return
	}

	for i, country := range apiCountries {
		_, err := job.countryRepo.Upsert(&country)
		if err != nil {
			log.Printf("ERROR processing country %s (%s): %v", country.Name, country.Code, err)
			continue
		}

		if i%50 == 0 {
			log.Printf("Processed %d countries", i+1)
		}
	}

	log.Println("Country Job Completed.")
}
