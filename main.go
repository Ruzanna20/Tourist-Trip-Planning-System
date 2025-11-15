package main

import (
	"log"
	"travel-planning/database"

	// jobservice "travel-planning/jobService"
	"travel-planning/repository"
	"travel-planning/services"
)

func main() {
	log.Println("Start")

	db, err := database.NewDB()
	if err != nil {
		log.Fatalf("FATAL: DB connection failed: %v", err)
	}
	defer db.Close()

	sqlConn := db.GetConn()
	countryRepo := repository.NewCountryRepository(sqlConn)
	cityRepo := repository.NewCityRepository(sqlConn)
	attractionRepo := repository.NewAttractionRepository(sqlConn)
	seeder := services.NewDataSeeder(countryRepo, cityRepo, attractionRepo)

	//country
	log.Println("Running Data Seeder for Countries...")

	err = seeder.SeedCountries()
	if err != nil {
		log.Fatalf("CRITICAL: Data Seeding failed. Error: %v", err)
	}

	// countryJob := jobservice.NewCountryJob(countryRepo)

	// go func ()  {
	// 	Interval := 24 * time.Hour
	// 	log.Printf("Country Job run every %s",Interval)

	// 	ticker := time.NewTicker(Interval)
	// 	defer ticker.Stop()

	// 	countryJob.RunJob()

	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			countryJob.RunJob()
	// 		}
	// 	}
	// }()
	log.Println("SUCCESS:All country data has been processed and saved.")

	//city
	log.Println("Running Data Seeder for Cities...")

	citiesFilePath := "data/worldcities.csv"
	if err = seeder.SeedCities(citiesFilePath); err != nil {
		log.Fatalf("CRITICAL: Data seedng failed.Error: %v", err)
	}
	log.Println("SUCCESS: All city data has been processed and saved.")

	//attraction
	attractionFilePath := "data/destinations.csv"
	if err = seeder.SeedAttractions(attractionFilePath); err != nil {
		log.Fatalf("CRITICAL: Data seedng failed.Error: %v", err)
	}
	log.Println("SUCCESS: All attraction data has been processed and saved.")

	log.Println("Finished seeding")
}
