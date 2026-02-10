package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"travel-planning/models"
	"travel-planning/services"
)

type UserHandlers struct {
	UserService *services.UserService
}

func NewUserHandlers(userService *services.UserService) *UserHandlers {
	return &UserHandlers{
		UserService: userService,
	}
}

// RegisterUserHandler godoc
// @Summary New user registration
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.UserRegistrationRequest true "Registration information"
// @Success 201 {object} map[string]interface{} "user_id"
// @Router /api/users/register [post]
func (h *UserHandlers) RegisterUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		slog.Warn("Method not allowed for registration", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req models.UserRegistrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("Failed to decode registration request", "error", err)
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" {
		slog.Warn("Registration attempt with missing fields", "email", req.Email)
		http.Error(w, "First name, Last Name, email, and password are required.", http.StatusBadRequest)
		return
	}

	userID, err := h.UserService.RegisterUser(req)
	if err != nil {
		slog.Error("User registration failed", "email", req.Email, "error", err)
		http.Error(w, "Failed to register user due to an internal error.", http.StatusInternalServerError)
		return
	}

	slog.Info("User registered successfully via API", "user_id", userID, "email", req.Email)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userID,
	})
}

// SetPreferencesHandler godoc
// @Summary Saving user preferences
// @Security BearerAuth
// @Tags Users
// @Accept json
// @Produce json
// @Param preferences body models.UserPreferences true "Preferences"
// @Success 201 {object} map[string]interface{} "preference_id"
// @Router /api/users/preferences [post]
func (s *UserHandlers) SetPreferencesHandler(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.Header.Get("X-User-ID")
	userID, _ := strconv.Atoi(userIDStr)
	l := slog.With("user_id", userID, "path", r.URL.Path)

	if r.Method != http.MethodPost {
		l.Warn("Method not allowed for preferences", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if userID <= 0 {
		l.Error("Unauthorized preference update attempt")
		http.Error(w, "Authentication error: Invalid User ID", http.StatusUnauthorized)
		return
	}

	var req models.UserPreferences
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		l.Warn("Failed to decode preferences body", "error", err)
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}

	l.Info("Updating user preferences")
	prefID, err := s.UserService.SavePreferences(userID, req)
	if err != nil {
		l.Error("Failed to save user preferences", "error", err)
		http.Error(w, "Failed to save user preferences", http.StatusInternalServerError)
		return
	}

	l.Info("User preferences saved successfully", "preference_id", prefID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"preference_id": prefID,
	})
}
