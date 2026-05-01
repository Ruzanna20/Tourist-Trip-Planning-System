package services

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"time"
	"travel-planning/internal/cache"
	"travel-planning/models"
	"travel-planning/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepo    *repository.UserRepository
	JWTService  *JWTService
	RedisCache  *cache.RedisCache
	MailService *MailService
}

func NewAuthService(userRepo *repository.UserRepository, jwtService *JWTService, redisCache *cache.RedisCache, mailService *MailService) *AuthService {
	return &AuthService{
		UserRepo:    userRepo,
		JWTService:  jwtService,
		RedisCache:  redisCache,
		MailService: mailService,
	}
}

func generateCode() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func (h *AuthService) Login(email, password string) (string, string, *models.User, string, error) {
	l := slog.With("email", email)
	l.Debug("Login attempt started")

	user, err := h.UserRepo.GetByEmail(email)
	if err != nil || user == nil {
		l.Warn("Login failed: user not found")
		return "", "", nil, "", fmt.Errorf("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		l.Warn("Login failed: incorrect password", "user_id", user.UserID)
		return "", "", nil, "", fmt.Errorf("invalid credentials")
	}

	if !user.IsVerified {
		code := generateCode()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := h.RedisCache.SaveVerificationCode(ctx, email, code)
		if err != nil {
			l.Error("Failed to save code to Redis", "error", err)
			return "", "", nil, "", fmt.Errorf("internal server error")
		}

		go func() {
			err := h.MailService.SendVerificationEmail(email, code)
			if err != nil {
				slog.Error("Failed to send verification email", "email", email, "error", err)
			} else {
				slog.Info("Verification email sent successfully", "email", email)
			}
		}()

		return "", "", user, "unverifed", nil

	}
	token, err := h.JWTService.GenerateToken(user.UserID)
	if err != nil {
		l.Error("Failed to generate access token", "user_id", user.UserID, "error", err)
		return "", "", nil, "", err
	}
	refreshToken, err := h.JWTService.GenerateRefreshToken(user.UserID)
	if err != nil {
		l.Error("Failed to generate refresh token", "user_id", user.UserID, "error", err)
		return "", "", nil, "", err
	}

	l.Info("User logged in successfully", "user_id", user.UserID)
	return token, refreshToken, user, "verifed", nil
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
