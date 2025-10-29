package main

import (
	"log"
	"os"
	"travel-planning/database"
	"travel-planning/repository"
	"travel-planning/services"
)

func main() {
	log.Println("Start")

	db,err := database.NewDB()
	if err != nil {
		log.Fatalf("FATAL: DB connection failed: %v",err)
	}
	defer db.Close()

	sqlConn := db.GetConn()
	countryRepo := repository.NewCountryRepository(sqlConn)
	cityRepo := repository.NewCityRepository(sqlConn)
	attractionRepo := repository.NewAttractionRepository(sqlConn)
	seeder := services.NewDataSeeder(countryRepo,cityRepo,attractionRepo)

	//country
	log.Println("Running Data Seeder for Countries...")

	err = seeder.SeedCountries()
	if err != nil {
			log.Fatalf("CRITICAL: Data Seeding failed. Error: %v", err) 
		}

	log.Println("SUCCESS:All country data has been processed and saved.")


	//city
	log.Println("Running Data Seeder for Cities...")

	citiesFilePath := "data/worldcities.csv"
	if err = seeder.SeedCities(citiesFilePath);err != nil {
		log.Fatalf("CRITICAL: Data seedng failed.Error: %v",err)
	}
	log.Println("SUCCESS: All city data has been processed and saved.")


	//attraction
	attractionFilePath := "data/destinations.csv" 
	if _, err := os.Stat(attractionFilePath); os.IsNotExist(err) {
		log.Fatalf("Attraction CSV file not found at %s.", attractionFilePath)
	}

	if err := seeder.SeedAttractions(attractionFilePath); err != nil {
		log.Fatalf("CRITICAL: Data seedng failed.Error: %v",err)
	}
	log.Println("SUCCESS: All attraction data has been processed and saved.")

	log.Println("Finished seeding")
}