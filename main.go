package main

import (
	"context"
	"log/slog"
	"os"

	"time"

	"travel-planning/database"
	"travel-planning/handlers"
	"travel-planning/internal/cache"
	"travel-planning/internal/kafka"
	"travel-planning/internal/notifications"

	jobservice "travel-planning/jobService"
	"travel-planning/server"

	"travel-planning/repository"
	"travel-planning/services"
)

// @title Tourist Trip Planning System API
// @version 1.0
// @description This is the API server for the Tourist Trip Planning System.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func init() {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(handler)

	slog.SetDefault(logger)
}

func main() {
	slog.Info("Appliaction starting up")

	db, err := database.NewDB()
	if err != nil {
		slog.Error("FATAL: DB connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	sqlConn := db.GetConn()

	redisAddr := os.Getenv("REDIS_ADDR")
	cacheService := cache.NewRedisCache(redisAddr)

	countryRepo := repository.NewCountryRepository(sqlConn)
	cityRepo := repository.NewCityRepository(sqlConn)
	attractionRepo := repository.NewAttractionRepository(sqlConn)
	hotelRepo := repository.NewHotelRepository(sqlConn)
	restaurantRepo := repository.NewRestaurantRepository(sqlConn)
	flightRepo := repository.NewFlightRepository(sqlConn)
	userRepo := repository.NewUserRepository(sqlConn)
	userPreferencesRepo := repository.NewUserPreferencesRepository(sqlConn)
	tripRepo := repository.NewTripRepository(sqlConn)
	itineraryRepo := repository.NewTripItineraryRepository(sqlConn)
	itineraryActivitiesRepo := repository.NewItineraryActivitiesRepository(sqlConn)
	reviewRepo := repository.NewReviewRepository(sqlConn)
	notificationRepo := repository.NewNotificationRepository(sqlConn)

	amadeusService := services.NewAmadeusService()
	countryAPIService := services.NewCountryAPIService(cacheService)
	cityAPIService := services.NewCityAPIService(cacheService)
	attractionAPIService := services.NewAttractionAPIService(cacheService)
	hotelAPIService := services.NewHotelAPIService(cacheService)
	restaurantAPIService := services.NewRestaurantAPIService(cacheService)
	flightAPIService := services.NewFlightAPIService(amadeusService, cityRepo, cacheService)
	notificationService := services.NewNotificationService(notificationRepo)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		slog.Warn("JWT_SECRET not found in environment, using default development secret")
		jwtSecret = "default-development-secret-must-be-changed"
	}
	jwtService := services.NewJWTService(jwtSecret)
	mailService := services.NewMailService(
		"smtp.gmail.com",
		587,
		os.Getenv("EMAIL_USER"),
		os.Getenv("EMAIL_PASSWORD"),
	)

	authService := services.NewAuthService(userRepo, jwtService, cacheService, mailService)
	userService := services.NewUserService(userRepo, userPreferencesRepo)
	resourceService := services.NewResourceService(hotelRepo, cityRepo, attractionRepo, countryRepo, restaurantRepo, flightRepo)
	reviewService := services.NewReviewService(reviewRepo)

	kafkaProducer := kafka.NewProducer("kafka:9092")
	defer kafkaProducer.Close()

	tripPlanningService := services.NewTripPlanningService(
		tripRepo,
		itineraryRepo,
		itineraryActivitiesRepo,
		flightRepo,
		hotelRepo,
		attractionRepo,
		restaurantRepo,
		userPreferencesRepo,
		kafkaProducer,
		cacheService,
	)

	hub := notifications.NewHub()

	kafkaConsumer := kafka.NewConsumer([]string{"kafka:9092"}, "trip-requests", "trip-service-group", tripPlanningService, notificationService, hub)
	defer kafkaConsumer.Close()

	go kafkaConsumer.Start(context.Background())

	seeder := services.NewDataSeeder(
		countryRepo,
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

	slog.Info("Jobs started in background")
	interval := 24 * time.Hour

	// Country Job
	countryJob := jobservice.NewCountryJob(seeder)
	go func() {
		slog.Info("Country Job scheduled", "interval", interval)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		countryJob.RunJob()
		for range ticker.C {
			countryJob.RunJob()
		}
	}()

	// // City Job
	// cityJob := jobservice.NewCityJob(seeder)
	// go func() {
	// 	slog.Info("City Job scheduled", "interval", interval)
	// 	ticker := time.NewTicker(interval)
	// 	defer ticker.Stop()

	// 	cityJob.RunJob()
	// 	for range ticker.C {
	// 		cityJob.RunJob()
	// 	}
	// }()

	// // Attraction Job
	// attractionJob := jobservice.NewAttractionJob(seeder)
	// go func() {
	// 	slog.Info("Attraction Job scheduled", "interval", interval)
	// 	ticker := time.NewTicker(interval)
	// 	defer ticker.Stop()

	// 	attractionJob.RunJob()
	// 	for range ticker.C {
	// 		attractionJob.RunJob()
	// 	}
	// }()

	// // Hotel Job
	// hotelJob := jobservice.NewHotelJob(seeder)
	// go func() {
	// 	slog.Info("Hotel Job scheduled", "interval", interval)
	// 	ticker := time.NewTicker(interval)
	// 	defer ticker.Stop()

	// 	hotelJob.RunJob()
	// 	for range ticker.C {
	// 		hotelJob.RunJob()
	// 	}
	// }()

	// // Restaurant Job
	// restaurantJob := jobservice.NewRestaurantJob(seeder)
	// go func() {
	// 	slog.Info("Restaurant Job scheduled", "interval", interval)
	// 	ticker := time.NewTicker(interval)
	// 	defer ticker.Stop()

	// 	restaurantJob.RunJob()
	// 	for range ticker.C {
	// 		restaurantJob.RunJob()
	// 	}
	// }()

	// // Flight Job
	// flightJob := jobservice.NewFlightJob(seeder)
	// go func() {
	// 	slog.Info("Flight Job scheduled", "interval", interval)
	// 	ticker := time.NewTicker(interval)
	// 	defer ticker.Stop()

	// 	flightJob.RunJob()
	// 	for range ticker.C {
	// 		flightJob.RunJob()
	// 	}
	// }()

	authHandlers := handlers.NewAuthHandlers(authService)
	userHandlers := handlers.NewUserHandlers(userService, authService)
	resourceHandlers := handlers.NewResourceHandlers(resourceService)
	reviewHandlers := handlers.NewReviewHandlers(reviewService)

	notificationHandlers := handlers.NewNotificationHandlers(notificationService)

	tripHandlers := handlers.NewTripHandlers(tripPlanningService)

	appServer := server.NewAppServer(
		authHandlers,
		resourceHandlers,
		reviewHandlers,
		userHandlers,
		tripHandlers,
		jwtService,
		hub,
		notificationHandlers,
	)
	appServer.Start(":8080")

	select {}
}
