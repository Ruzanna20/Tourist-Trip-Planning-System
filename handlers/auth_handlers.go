package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"travel-planning/models"
	"travel-planning/repository"
	"travel-planning/services"

	"golang.org/x/crypto/bcrypt"
)

type CustomClaims = services.CustomClaims

const AdminUser = "admin"

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Message      string `json:"message"`
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthHandlers struct {
	JWTService *services.JWTService
	CityRepo   *repository.CityRepository
	UserRepo   *repository.UserRepository
}

func NewAuthHandlers(jwtService *services.JWTService, cityRepo *repository.CityRepository, userRepo *repository.UserRepository) *AuthHandlers {
	return &AuthHandlers{
		JWTService: jwtService,
		CityRepo:   cityRepo,
		UserRepo:   userRepo,
	}
}

func (h *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.UserRepo.GetByEmail(creds.Username)
	if err != nil {
		log.Printf("DB error on login: %v", err)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}
	if user == nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, err := h.JWTService.GenerateToken(user.UserID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}
	refreshToken, err := h.JWTService.GenerateRefreshToken(user.UserID)
	if err != nil {
		http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message:      "Login successful",
		Token:        token,
		RefreshToken: refreshToken,
	})
}

func (h *AuthHandlers) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}

	claims, err := h.JWTService.ValidateToken(req.RefreshToken)
	if err != nil {
		http.Error(w, "Invalid refresh token. Please log in.", http.StatusUnauthorized)
		return
	}

	newAccessToken, err := h.JWTService.GenerateToken(claims.UserID)
	if err != nil {
		log.Printf("Error generating new access token: %v", err)
		http.Error(w, "Error generating new access token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message: "Token refreshed successfully",
		Token:   newAccessToken,
	})
}

func (h *AppHandlers) ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid User ID format", http.StatusInternalServerError)
		return
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

func (h *AppHandlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.UserRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		http.Error(w, "First name, Last Name, email, and password are required.", http.StatusBadRequest)
		return
	}

	newUser := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	userID, err := h.UserRepo.Insert(newUser, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			http.Error(w, "User already registered.", http.StatusConflict)
			return
		}
		log.Printf("Registration failed: %v", err)
		http.Error(w, "Failed to register user due to an internal error.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "User registered successfully.",
		"user_id": userID,
	})
}
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message: "Server is running and healthy",
	})
}
