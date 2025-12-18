package services

import (
	"fmt"
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
	user, err := h.UserRepo.GetByEmail(email)
	if err != nil || user == nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	token, err := h.JWTService.GenerateToken(user.UserID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err := h.JWTService.GenerateRefreshToken(user.UserID)
	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}

func (h *AuthService) RefreshToken(refreshToken string) (string, error) {
	claims, err := h.JWTService.ValidateToken(refreshToken)
	if err != nil {
		return "", err
	}

	newAccessToken, err := h.JWTService.GenerateToken(claims.UserID)
	if err != nil {
		return "", nil
	}

	return newAccessToken, err
}
