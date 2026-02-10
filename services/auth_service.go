package services

import (
	"fmt"
	"log/slog"
	"travel-planning/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo   *repository.UserRepository
	JWTService *JWTService
}

func NewAuthService(userRepo *repository.UserRepository, jwtService *JWTService) *AuthService {
	return &AuthService{
		UserRepo:   userRepo,
		JWTService: jwtService,
	}
}

func (h *AuthService) Login(email, password string) (string, string, error) {
	l := slog.With("email", email)
	l.Debug("Login attempt started")

	user, err := h.UserRepo.GetByEmail(email)
	if err != nil || user == nil {
		l.Warn("Login failed: user not found")
		return "", "", fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		l.Warn("Login failed: incorrect password", "user_id", user.UserID)
		return "", "", fmt.Errorf("invalid credentials")
	}

	token, err := h.JWTService.GenerateToken(user.UserID)
	if err != nil {
		l.Error("Failed to generate access token", "user_id", user.UserID, "error", err)
		return "", "", err
	}
	refreshToken, err := h.JWTService.GenerateRefreshToken(user.UserID)
	if err != nil {
		l.Error("Failed to generate refresh token", "user_id", user.UserID, "error", err)
		return "", "", err
	}

	l.Info("User logged in successfully", "user_id", user.UserID)
	return token, refreshToken, err
}

func (h *AuthService) RefreshToken(refreshToken string) (string, error) {
	slog.Debug("Refresh token attempt started")

	claims, err := h.JWTService.ValidateToken(refreshToken)
	if err != nil {
		slog.Warn("Token refresh failed: invalid or expired refresh token", "error", err)
		return "", err
	}

	newAccessToken, err := h.JWTService.GenerateToken(claims.UserID)
	if err != nil {
		slog.Error("Failed to generate new access token during refresh", "user_id", claims.UserID, "error", err)
		return "", nil
	}

	slog.Info("Token refreshed successfully", "user_id", claims.UserID)
	return newAccessToken, err
}
