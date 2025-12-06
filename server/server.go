package server

import (
	"log"
	"net/http"

	"travel-planning/handlers"
	"travel-planning/services"
)

type AppServer struct {
	ResourceHandlers *handlers.AppHandlers
	AuthHandlers     *handlers.AuthHandlers
	JWTService       *services.JWTService
}

func NewAppServer(resourceH *handlers.AppHandlers, authH *handlers.AuthHandlers, jwtS *services.JWTService) *AppServer {
	return &AppServer{
		ResourceHandlers: resourceH,
		AuthHandlers:     authH,
		JWTService:       jwtS,
	}
}

func (s *AppServer) Start(port string) {
	log.Printf("Server starting on port %s", port)

	http.HandleFunc("/health", handlers.HealthHandler)
	http.HandleFunc("/login", s.AuthHandlers.LoginHandler)
	http.HandleFunc("/refresh", s.AuthHandlers.RefreshHandler)
	authMiddleware := s.JWTService.AuthMiddleware

	http.HandleFunc("/protected", authMiddleware(s.ResourceHandlers.ProtectedHandler))
	http.HandleFunc("/api/cities", authMiddleware(s.ResourceHandlers.GetAllCitiesHandler))
	http.HandleFunc("/api/countries", authMiddleware(s.ResourceHandlers.GetAllCountriesHandler))
	http.HandleFunc("/api/attractions", authMiddleware(s.ResourceHandlers.GetAllAttractionssHandler))
	http.HandleFunc("/api/hotels", authMiddleware(s.ResourceHandlers.GetAllHotelsHandler))
	http.HandleFunc("/api/restaurants", authMiddleware(s.ResourceHandlers.GetAllRestaurantssHandler))
	http.HandleFunc("/api/flights", authMiddleware(s.ResourceHandlers.GetAllFlightsHandler))
	http.HandleFunc("/api/trips", authMiddleware(s.ResourceHandlers.GetTripsHandler))

	http.HandleFunc("/api/users/register", s.ResourceHandlers.RegisterUserHandler)
	http.HandleFunc("/api/users/preferences", authMiddleware(s.ResourceHandlers.SetPreferencesHandler))

	http.HandleFunc("/api/trips/create", authMiddleware(s.ResourceHandlers.CreateTripHandler))
	log.Fatal(http.ListenAndServe(port, nil))
}
