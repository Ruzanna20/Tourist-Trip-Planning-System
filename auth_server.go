package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"travel-planning/database"
	"travel-planning/repository"
	"travel-planning/services"
)

type CustomClaims = services.CustomClaims

const AdminUser = "admin"

var jwtService *services.JWTService
var jwtSecret string

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type AppHandlers struct {
	HotelRepo      *repository.HotelRepository
	CityRepo       *repository.CityRepository
	AttractionRepo *repository.AttractionRepository
	CountryRepo    *repository.CountryRepository
}

func NewAppHandlers(HotelRepo *repository.HotelRepository,
	CityRepo *repository.CityRepository,
	AttractionRepo *repository.AttractionRepository,
	CountryRepo *repository.CountryRepository) *AppHandlers {
	return &AppHandlers{
		HotelRepo:      HotelRepo,
		CityRepo:       CityRepo,
		AttractionRepo: AttractionRepo,
		CountryRepo:    CountryRepo,
	}
}

func init() {
	jwtSecret = os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Println("JWT_SECRET not found in .env. Using default development secret.")
		jwtSecret = "default-development-secret-must-be-changed"
	}

	expiry := os.Getenv("JWT_EXPIRY_HOURS")
	if expiry == "" {
		expiry = "24"
	}

	jwtService = services.NewJWTService(jwtSecret, expiry)
	log.Println("JWT Service initialized.")
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	if creds.Username == AdminUser {
		const AdminPassword = "password"

		if creds.Password == AdminPassword {
			token, err := jwtService.GenerateToken(1)
			if err != nil {
				http.Error(w, "Error generating token", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Response{
				Message: "Login successful",
				Token:   token,
			})
			return
		}
	}
	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message: "Server is running and healthy",
	})
}

func (h *AppHandlers) protectedHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID format", http.StatusInternalServerError)
	}

	cities, err := h.CityRepo.GetAllCityLocations()
	if err != nil {
		log.Printf("DB error fetching cities:%v", err)
		http.Error(w, "Invalid User ID format", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message: fmt.Sprintf("User %d is authorized. Fetched %d cities.", userID, len(cities)),
	})
}
func main() {
	port := ":8080"

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
	appHandlers := NewAppHandlers(hotelRepo, cityRepo, attractionRepo, countryRepo)

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/health", healthHandler)
	http.HandleFunc("/protected", jwtService.AuthMiddleware(appHandlers.protectedHandler))

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
