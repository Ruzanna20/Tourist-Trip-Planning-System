package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"travel-planning/services"
)

type CustomClaims = services.CustomClaims

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthHandlers struct {
	AuthService *services.AuthService
}

func NewAuthHandlers(authService *services.AuthService) *AuthHandlers {
	return &AuthHandlers{
		AuthService: authService,
	}
}

// LoginHandler godoc
// @Summary User login
// @Description Enter your email address and password to receive a JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param credentials body Credentials true "Login information"
// @Success 200 {object} Response
// @Failure 401 {string} string "Invalid information"
// @Router /login [post]
func (h *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Warn("Method not allowed", "method", r.Method, "path", r.URL.Path)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var creds Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		slog.Warn("Failed to decode login credentials", "error", err)
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}

	if creds.Username == "" || creds.Password == "" {
		slog.Warn("Login attempt with missing credentials", "username", creds.Username)
		http.Error(w, "Username and password required", http.StatusBadRequest)
		return
	}

	token, refreshToken, err := h.AuthService.Login(creds.Username, creds.Password)
	if err != nil {
		slog.Warn("Unauthorized login attempt", "username", creds.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	slog.Info("User logged in successfully", "username", creds.Username)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Token:        token,
		RefreshToken: refreshToken,
	})
}

// RefreshHandler godoc
// @Summary Refresh access token
// @Description Get a new access token using a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refresh_token body RefreshRequest true "Refresh token request"
// @Success 200 {object} Response
// @Failure 401 {string} string "Invalid refresh token"
// @Router /refresh [post]
func (h *AuthHandlers) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Warn("Method not allowed for refresh", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Failed to decode refresh request", "error", err)
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}

	newAccessToken, err := h.AuthService.RefreshToken(req.RefreshToken)
	if err != nil {
		slog.Error("Error generating new access token", "error", err)
		http.Error(w, "Error generating new access token", http.StatusInternalServerError)
		return
	}

	slog.Info("Access token refreshed successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Token: newAccessToken,
	})
}
