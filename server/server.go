package server

import (
	"log/slog"
	"net/http"
	"strconv"

	"travel-planning/handlers"
	"travel-planning/internal/notifications"
	"travel-planning/services"

	_ "travel-planning/docs"

	corsHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type AppServer struct {
	AuthHandlers         *handlers.AuthHandlers
	ResourceHandlers     *handlers.ResourceHandlers
	ReviewHandlers       *handlers.ReviewHandlers
	UserHandlers         *handlers.UserHandlers
	TripHandlers         *handlers.TripHandlers
	JWTService           *services.JWTService
	NotificationsHub     *notifications.Hub
	NotificationHandlers *handlers.NotificationHandlers
}

func NewAppServer(
	authH *handlers.AuthHandlers,
	resourceH *handlers.ResourceHandlers,
	reviewH *handlers.ReviewHandlers,
	userH *handlers.UserHandlers,
	tripH *handlers.TripHandlers,
	jwtS *services.JWTService,
	notificationsHub *notifications.Hub,
	notificationH *handlers.NotificationHandlers,
) *AppServer {
	return &AppServer{
		AuthHandlers:         authH,
		ResourceHandlers:     resourceH,
		ReviewHandlers:       reviewH,
		UserHandlers:         userH,
		TripHandlers:         tripH,
		JWTService:           jwtS,
		NotificationsHub:     notificationsHub,
		NotificationHandlers: notificationH,
	}
}

func (s *AppServer) Start(port string) {
	slog.Info("Starting Application Server", "port", port)

	r := mux.NewRouter()
	r.Use(mux.CORSMethodMiddleware(r))

	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
	authMiddleware := s.JWTService.AuthMiddleware

	// Auth
	r.HandleFunc("/login", s.AuthHandlers.LoginHandler).Methods("POST")
	r.HandleFunc("/refresh", s.AuthHandlers.RefreshHandler).Methods("POST")
	r.HandleFunc("/api/auth/verify", s.AuthHandlers.VerifyEmailHandler).Methods("POST", "OPTIONS")

	// Resources
	r.HandleFunc("/api/cities", authMiddleware(s.ResourceHandlers.GetAllCitiesHandler)).Methods("GET")
	r.HandleFunc("/api/countries", authMiddleware(s.ResourceHandlers.GetAllCountriesHandler)).Methods("GET")
	r.HandleFunc("/api/attractions", authMiddleware(s.ResourceHandlers.GetAllAttractionssHandler)).Methods("GET")
	r.HandleFunc("/api/hotels", authMiddleware(s.ResourceHandlers.GetAllHotelsHandler)).Methods("GET")
	r.HandleFunc("/api/restaurants", authMiddleware(s.ResourceHandlers.GetAllRestaurantssHandler)).Methods("GET")
	r.HandleFunc("/api/flights", authMiddleware(s.ResourceHandlers.GetAllFlightsHandler)).Methods("GET")

	// Notifications
	r.HandleFunc("/api/notifications", authMiddleware(s.NotificationHandlers.GetMyNotifications)).Methods("GET")
	r.HandleFunc("/api/notifications/{id}/read", authMiddleware(s.NotificationHandlers.MarkAsRead)).Methods("POST")
	r.HandleFunc("/api/notifications/unread-count", authMiddleware(s.NotificationHandlers.GetUnreadCount)).Methods("GET")
	r.HandleFunc("/api/notifications/{id}", authMiddleware(s.NotificationHandlers.DeleteNotification)).Methods("DELETE")

	// Reviews
	r.HandleFunc("/api/reviews", authMiddleware(s.ReviewHandlers.GetUserReviewsHandler)).Methods("GET")
	r.HandleFunc("/api/reviews", authMiddleware(s.ReviewHandlers.CreateReviewHandler)).Methods("POST")
	r.HandleFunc("/api/reviews/{id}", authMiddleware(s.ReviewHandlers.DeleteReviewHandler)).Methods("DELETE")

	// Trips
	r.HandleFunc("/api/trips", authMiddleware(s.TripHandlers.GetUserTripsHandler)).Methods("GET")
	r.HandleFunc("/api/trips/{id}", authMiddleware(s.TripHandlers.DeleteTripHandler)).Methods("DELETE")
	r.HandleFunc("/api/trips/{id}/generate-options", authMiddleware(s.TripHandlers.GenerateTripOptions)).Methods("POST")
	r.HandleFunc("/api/trips/{id}/select-option", authMiddleware(s.TripHandlers.SelectTripOption)).Methods("POST")
	r.HandleFunc("/api/trips/create", authMiddleware(s.TripHandlers.CreateTripHandler)).Methods("POST")

	r.HandleFunc("/api/users/me/visited", authMiddleware(s.ResourceHandlers.GetVisitedEntitiesHandler)).Methods("GET")
	r.HandleFunc("/api/trips/{id}/complete", authMiddleware(s.TripHandlers.CompleteTripHandler)).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/itinerary/activities/{id}/swap", authMiddleware(s.TripHandlers.SwapActivityHandler)).Methods("POST", "OPTIONS")

	// Itinerary & Activities
	r.HandleFunc("/api/trips/{id}/itinerary", authMiddleware(s.TripHandlers.GetTripItineraryHandler)).Methods("GET")
	r.HandleFunc("/api/itineraries/{id}/activities", authMiddleware(s.TripHandlers.GetActivitiesHandler)).Methods("GET")

	// Users
	r.HandleFunc("/api/users/register", s.UserHandlers.RegisterUserHandler).Methods("POST")
	r.HandleFunc("/api/users/preferences", authMiddleware(s.UserHandlers.GetPreferencesHandler)).Methods("GET")
	r.HandleFunc("/api/users/preferences", authMiddleware(s.UserHandlers.SetPreferencesHandler)).Methods("POST")

	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		userIDStr := r.URL.Query().Get("userID")
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			slog.Error("Invalid userID in WebSocket connection", "error", err)
			return
		}
		s.NotificationsHub.HandleWS(w, r, userID)
	}).Methods("GET")

	slog.Info("Routes registered successfully")

	corsHandler := corsHandlers.CORS(
		corsHandlers.AllowedOrigins([]string{"http://localhost:5173"}),
		corsHandlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		corsHandlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
		corsHandlers.AllowCredentials(),
	)(r)

	if err := http.ListenAndServe(port, corsHandler); err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}
