package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
	"travel-planning/services"
)

type CustomClaims = services.CustomClaims

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Token        string                 `json:"token,omitempty"`
	RefreshToken string                 `json:"refresh_token,omitempty"`
	User         map[string]interface{} `json:"user,omitempty"`
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

	token, refreshToken, user, status, err := h.AuthService.Login(creds.Username, creds.Password)
	if err != nil {
		slog.Warn("Unauthorized login attempt", "username", creds.Username)
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	if status == "unverifed" {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "unverified",
			"email":   creds.Username,
			"message": "Your account is not verified. A verification code has been sent to your email.",
		})
		return
	}

	slog.Info("User logged in successfully", "username", creds.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Token:        token,
		RefreshToken: refreshToken,
		User: map[string]interface{}{
			"user_id":    user.UserID,
			"first_name": user.FirstName,
			"last_name":  user.LastName,
			"email":      user.Email,
		},
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

// VerifyEmailHandler godoc
// @Summary Verify user email
// @Description Verify a user's email address using the 6-digit code stored in Redis
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body struct{Email string `json:"email"`; Code string `json:"code"`} true "Verification details"
// @Success 200 {object} map[string]string "message: Email verified successfully"
// @Failure 400 {string} string "Invalid or expired verification code"
// @Router /api/auth/verify [post]
func (h *AuthHandlers) VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("VerifyEmailHandler hit!", "email", r.Method)

	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	slog.Info("RECEIVED_DATA", "email", req.Email, "code", req.Code)

	ctx := context.Background()
	savedCode, err := h.AuthService.RedisCache.GetVerificationCode(ctx, req.Email)
	if err != nil {
		slog.Warn("Verification failed: no code found for email", "email", req.Email)
		http.Error(w, "Invalid or expired verification code", http.StatusBadRequest)
		return
	}

	slog.Info("REDIS_DATA", "saved", savedCode, "error", err)

	savedCode = strings.Trim(savedCode, "\" \n\r\t")
	inputCode := strings.Trim(req.Code, "\" \n\r\t")

	slog.Info("Final Check", "stored", savedCode, "input", inputCode)

	slog.Info("VERIFICATION_DEBUG",
		"stored_raw", savedCode,
		"input_raw", inputCode,
		"match", savedCode == inputCode)

	slog.Info("Comparing codes", "stored", savedCode, "input", inputCode)
	if savedCode != inputCode {
		slog.Warn("Verification failed: code mismatch", "email", req.Email)
		http.Error(w, "Invalid or expired verification code", http.StatusBadRequest)
		return
	}

	if err := h.AuthService.UserRepo.MarkAsVerified(req.Email); err != nil {
		slog.Error("Failed to update database", "error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	slog.Info("Email verified successfully", "email", req.Email)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Email verified successfully. You can now log in.",
	})
}
