package main

import (
	"flag"
	"log"
	"os"

	"travel-planning/database"
	"travel-planning/handlers"
	"travel-planning/server"

	"travel-planning/repository"
	"travel-planning/services"
)

func main() {
	log.Println("Start")
	seedFlag := flag.Bool("seed", false, "Set to true ran data seeding job")
	flag.Parse()

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
	flightRepo := repository.NewFlightRepository(sqlConn)

	userRepo := repository.NewUserRepository(sqlConn)
	userPreferencesRepo := repository.NewUserPreferencesRepository(sqlConn)
	tripRepo := repository.NewTripRepository(sqlConn)
	tripItineraryRepo := repository.NewTripItineraryRepository(sqlConn)
	itineraryActivitiesRepo := repository.NewItineraryActivitiesRepository(sqlConn)

	reviewRepo := repository.NewReviewRepository(sqlConn)

	amadeusService := services.NewAmadeusService()
	countryAPIService := services.NewCountryAPIService()
	cityAPIService := services.NewCityAPIService()
	attractionAPIService := services.NewAttractionAPIService()
	foursquareAPIService := services.NewFoursquareService()
	hotelAPIService := services.NewHotelAPIService(amadeusService, foursquareAPIService)
	restaurantAPIService := services.NewRestaurantAPIService()
	flightAPIService := services.NewFlightAPIService(amadeusService, cityRepo)

	tripPlanningService := services.NewTripPlanningService(tripRepo, tripItineraryRepo, itineraryActivitiesRepo)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default-development-secret-must-be-changed"
	}
	jwtService := services.NewJWTService(jwtSecret, "5")

	seeder := services.NewDataSeeder(countryRepo,
		cityRepo,
		attractionRepo,
		hotelRepo,
		restaurantRepo,
		flightRepo,
		countryAPIService,
		cityAPIService,
		attractionAPIService,
		hotelAPIService,
		restaurantAPIService,
		flightAPIService)

	if *seedFlag {
		log.Println("Starting Seeding Job...")
		// country
		err = seeder.SeedCountries()
		if err != nil {
			log.Fatalf("CRITICAL: Data Seeding failed. Error: %v", err)
		}

		// //city
		// if err = seeder.SeedCities(); err != nil {
		// 	log.Fatalf("CRITICAL: Data seedng failed.Error: %v", err)
		// }

		// //attraction
		// if err = seeder.SeedAttractions(); err != nil {
		// 	log.Fatalf("CRITICAL: Data seedng failed.Error: %v", err)

		// }

		// // hotel
		// if err = seeder.SeedHotels(); err != nil {
		// 	log.Fatalf("CRITICAL: Hotel Data seeding failed. Error: %v", err)
		// }

		// //restaurant
		// if err = seeder.SeedRestaurants(); err != nil {
		// 	log.Fatalf("CRITICAL: Restaurant Data seeding failed. Error: %v", err)
		// }

		// //flight
		// if err = seeder.SeedFlights(); err != nil {
		// 	log.Fatalf("CRITICAL: Flights Seeding failed. Error: %v", err)
		// }

		// log.Println("Seeding job finished. Exiting.")
		// os.Exit(0)
	}

	getResourceHandlers := handlers.NewAppHandlers(
		hotelRepo,
		cityRepo,
		attractionRepo,
		countryRepo,
		restaurantRepo,
		flightRepo,
		userRepo,
		userPreferencesRepo,
		tripRepo,
		reviewRepo,
		tripPlanningService)

	authHandlers := handlers.NewAuthHandlers(jwtService, cityRepo, userRepo)

	appServer := server.NewAppServer(getResourceHandlers, authHandlers, jwtService)
	appServer.Start(":8080")

	// log.Println("Jobs started in background...")
	// Interval := 24 * time.Hour

	// countryJob := jobservice.NewCountryJob(seeder)
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
	// go func() {
	// 	log.Printf("City Job run every %s", Interval)

	// 	ticker := time.NewTicker(Interval)
	// 	defer ticker.Stop()

	// 	cityJob.RunJob()
	// 	for range ticker.C {
	// 		cityJob.RunJob()
	// 	}
	// }()

	// attractionJob := jobservice.NewAttractionJob(seeder)
	// go func() {
	// 	log.Printf("Attraction Job run every %s", Interval)

	// 	ticker := time.NewTicker(Interval)
	// 	defer ticker.Stop()

	// 	attractionJob.RunJob()
	// 	for range ticker.C {
	// 		attractionJob.RunJob()
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

	// flightJob := jobservice.NewFlightJob(seeder)
	// go func() {
	// 	log.Printf("Flight Job run every %s", Interval)

	// 	ticker := time.NewTicker(Interval)
	// 	defer ticker.Stop()

	// 	flightJob.RunJob()
	// 	for range ticker.C {
	// 		flightJob.RunJob()
	// 	}
	// }()

	select {}
}
