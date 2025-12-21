package server

import (
	"log"
	"net/http"

	"travel-planning/handlers"
	"travel-planning/services"
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
	log.Printf("Server starting on port %s", port)

	authMiddleware := s.JWTService.AuthMiddleware

	http.HandleFunc("/login", s.AuthHandlers.LoginHandler)
	http.HandleFunc("/refresh", s.AuthHandlers.RefreshHandler)

	http.HandleFunc("/api/cities", authMiddleware(s.ResourceHandlers.GetAllCitiesHandler))
	http.HandleFunc("/api/countries", authMiddleware(s.ResourceHandlers.GetAllCountriesHandler))
	http.HandleFunc("/api/attractions", authMiddleware(s.ResourceHandlers.GetAllAttractionssHandler))
	http.HandleFunc("/api/hotels", authMiddleware(s.ResourceHandlers.GetAllHotelsHandler))
	http.HandleFunc("/api/restaurants", authMiddleware(s.ResourceHandlers.GetAllRestaurantssHandler))
	http.HandleFunc("/api/flights", authMiddleware(s.ResourceHandlers.GetAllFlightsHandler))

	http.HandleFunc("/api/reviews", authMiddleware(s.ReviewHandlers.CreateReviewHandler))

	http.HandleFunc("/api/trips/generate-options", authMiddleware(s.TripHandlers.GenerateTripOptions))
	http.HandleFunc("/api/trips/create", authMiddleware(s.TripHandlers.CreateTripHandler))
	// http.HandleFunc("/api/trips/", authMiddleware(s.TripHandlers.GetTripItineraryHandler))
	// http.HandleFunc("/api/itineraries/", authMiddleware(s.TripHandlers.GetActivitiesHandler))

	http.HandleFunc("/api/users/register", s.UserHandlers.RegisterUserHandler)
	http.HandleFunc("/api/users/preferences", authMiddleware(s.UserHandlers.SetPreferencesHandler))

	log.Fatal(http.ListenAndServe(port, nil))
}
