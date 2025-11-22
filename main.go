package main

import (
	"log"

	//"time"
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
	hotelRepo := repository.NewHotelRepository(sqlConn)
	restaurantRepo := repository.NewRestaurantRepository(sqlConn)
	amadeusService := services.NewAmadeusService()
	countryAPIService := services.NewCountryAPIService()
	cityAPIService := services.NewCityAPIService()
	hotelAPIService := services.NewHotelAPIService(amadeusService)
	restaurantAPIService := services.NewRestaurantAPIService(amadeusService)
	transportationAPIService := services.TransportationAPIService(amadeusService)
	seeder := services.NewDataSeeder(countryRepo,
		cityRepo,
		attractionRepo,
		hotelRepo,
		restaurantRepo,
		countryAPIService,
		cityAPIService,
		hotelAPIService,
		restaurantAPIService,
		transportationAPIService)

	log.Println("Starting Seeding...")

	//country
	err = seeder.SeedCountries()
	if err != nil {
		log.Fatalf("CRITICAL: Data Seeding failed. Error: %v", err)
	}

	// //city
	// if err = seeder.SeedCities(); err != nil {
	// 	log.Fatalf("CRITICAL: Data seedng failed.Error: %v", err)
	// }

	// //attraction
	// attractionFilePath := "data/destinations.csv"
	// if err = seeder.SeedAttractions(attractionFilePath); err != nil {
	// 	log.Fatalf("CRITICAL: Data seedng failed.Error: %v", err)
	// }

	// // hotel
	// if err = seeder.SeedHotels(); err != nil {
	// 	log.Fatalf("CRITICAL: Hotel Data seeding failed. Error: %v", err)
	// }

	// if err = seeder.SeedRestaurants(); err != nil {
	// 	log.Fatalf("CRITICAL: Restaurant Data seeding failed. Error: %v", err)
	// }

	if err = seeder.SeedTransportation(); err != nil {
		log.Fatalf("CRITICAL: Transportation Seeding failed. Error: %v", err)
	}

	// log.Println("Jobs started in background...")
	// Interval := 24 * time.Hour

	// countryJob := jobservice.NewCountryJob(countryRepo)
	// go func() {
	// 	log.Printf("Country Job run every %s", Interval)

	// 	ticker := time.NewTicker(Interval)
	// 	defer ticker.Stop()

	// 	countryJob.RunJob()
	// 	for range ticker.C {
	// 		countryJob.RunJob()
	// 	}
	// }()

	// cityJob := jobservice.NewCityJob(seeder)
	// go func ()  {
	// 	log.Printf("City Job run every %s", Interval)

	// 	ticker := time.NewTicker(Interval)
	// 	defer ticker.Stop()

	// 	cityJob.RunJob()
	// 	for range ticker.C {
	// 		cityJob.RunJob()
	// 	}
	// }()

	// hotelJob := jobservice.NewHotelJob(cityRepo, hotelRepo, hotelAPIService)
	// go func() {
	// 	log.Printf("Hotel Job run every %s", Interval)

	// 	ticker := time.NewTicker(Interval)
	// 	defer ticker.Stop()

	// 	hotelJob.RunJob()
	// 	for range ticker.C {
	// 		hotelJob.RunJob()
	// 	}
	// }()

	// restaurantJob := jobservice.NewRestaurantJob(cityRepo, restaurantRepo, restaurantAPIService)
	// go func() {
	// 	log.Printf("Restaurant Job run every %s", Interval)

	// 	ticker := time.NewTicker(Interval)
	// 	defer ticker.Stop()

	// 	restaurantJob.RunJob()
	// 	for range ticker.C {
	// 		restaurantJob.RunJob()
	// 	}
	// }()

	// select {}
}
