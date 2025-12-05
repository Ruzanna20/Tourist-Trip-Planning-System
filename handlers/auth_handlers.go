package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"travel-planning/repository"
	"travel-planning/services"
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
}

func NewAuthHandlers(jwtService *services.JWTService, cityRepo *repository.CityRepository) *AuthHandlers {
	return &AuthHandlers{
		JWTService: jwtService,
		CityRepo:   cityRepo,
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

	if creds.Username == AdminUser {
		const AdminPassword = "password"

		if creds.Password == AdminPassword {
			token, err := h.JWTService.GenerateToken(1)
			if err != nil {
				http.Error(w, "Error generating token", http.StatusInternalServerError)
				return
			}

			refreshToken, err := h.JWTService.GenerateRefreshToken(1)
			if err != nil {
				http.Error(w, "Error generating token", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(Response{
				Message:      "Login successful",
				Token:        token,
				RefreshToken: refreshToken,
			})
			return
		}
	}
	http.Error(w, "Invalid credentials", http.StatusUnauthorized)
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

func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message: "Server is running and healthy",
	})
}
