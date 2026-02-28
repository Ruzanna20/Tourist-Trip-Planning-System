package server

import (
	"log/slog"
	"net/http"

	"travel-planning/handlers"
	"travel-planning/services"

	_ "travel-planning/docs"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type AppServer struct {
	AuthHandlers     *handlers.AuthHandlers
	ResourceHandlers *handlers.ResourceHandlers
	ReviewHandlers   *handlers.ReviewHandlers
	UserHandlers     *handlers.UserHandlers
	TripHandlers     *handlers.TripHandlers
	JWTService       *services.JWTService
}

func NewAppServer(
	authH *handlers.AuthHandlers,
	resourceH *handlers.ResourceHandlers,
	reviewH *handlers.ReviewHandlers,
	userH *handlers.UserHandlers,
	tripH *handlers.TripHandlers,
	jwtS *services.JWTService,
) *AppServer {
	return &AppServer{
		AuthHandlers:     authH,
		ResourceHandlers: resourceH,
		ReviewHandlers:   reviewH,
		UserHandlers:     userH,
		TripHandlers:     tripH,
		JWTService:       jwtS,
	}
}

func (s *AppServer) Start(port string) {
	slog.Info("Starting Application Server", "port", port)

	r := mux.NewRouter()

	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	authMiddleware := s.JWTService.AuthMiddleware

	// Auth
	r.HandleFunc("/login", s.AuthHandlers.LoginHandler).Methods("POST")
	r.HandleFunc("/refresh", s.AuthHandlers.RefreshHandler).Methods("POST")

	// Resources
	r.HandleFunc("/api/cities", authMiddleware(s.ResourceHandlers.GetAllCitiesHandler)).Methods("GET")
	r.HandleFunc("/api/countries", authMiddleware(s.ResourceHandlers.GetAllCountriesHandler)).Methods("GET")
	r.HandleFunc("/api/attractions", authMiddleware(s.ResourceHandlers.GetAllAttractionssHandler)).Methods("GET")
	r.HandleFunc("/api/hotels", authMiddleware(s.ResourceHandlers.GetAllHotelsHandler)).Methods("GET")
	r.HandleFunc("/api/restaurants", authMiddleware(s.ResourceHandlers.GetAllRestaurantssHandler)).Methods("GET")
	r.HandleFunc("/api/flights", authMiddleware(s.ResourceHandlers.GetAllFlightsHandler)).Methods("GET")

	// Reviews
	r.HandleFunc("/api/reviews", authMiddleware(s.ReviewHandlers.GetUserReviewsHandler)).Methods("GET")
	r.HandleFunc("/api/reviews", authMiddleware(s.ReviewHandlers.CreateReviewHandler)).Methods("POST")

	// Trips
	r.HandleFunc("/api/trips", authMiddleware(s.TripHandlers.GetUserTripsHandler)).Methods("GET")
	r.HandleFunc("/api/trips/{id}", authMiddleware(s.TripHandlers.DeleteTripHandler)).Methods("DELETE")
	r.HandleFunc("/api/trips/{id}/generate-options", authMiddleware(s.TripHandlers.GenerateTripOptions)).Methods("POST")
	r.HandleFunc("/api/trips/{id}/select-option", authMiddleware(s.TripHandlers.SelectTripOption)).Methods("POST")
	r.HandleFunc("/api/trips/create", authMiddleware(s.TripHandlers.CreateTripHandler)).Methods("POST")

	// Itinerary & Activities
	r.HandleFunc("/api/trips/{id}/itinerary", authMiddleware(s.TripHandlers.GetTripItineraryHandler)).Methods("GET")
	r.HandleFunc("/api/itineraries/{id}/activities", authMiddleware(s.TripHandlers.GetActivitiesHandler)).Methods("GET")

	// Users
	r.HandleFunc("/api/users/register", s.UserHandlers.RegisterUserHandler).Methods("POST")
	r.HandleFunc("/api/users/preferences", authMiddleware(s.UserHandlers.GetPreferencesHandler)).Methods("GET")
	r.HandleFunc("/api/users/preferences", authMiddleware(s.UserHandlers.SetPreferencesHandler)).Methods("POST")

	slog.Info("Routes registered successfully")

	if err := http.ListenAndServe(port, r); err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}
