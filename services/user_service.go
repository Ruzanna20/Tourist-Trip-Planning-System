package services

import (
	"fmt"
	"log/slog"
	"strings"
	"travel-planning/models"
	"travel-planning/repository"
)

type UserService struct {
	UserRepo            *repository.UserRepository
	UserPreferencesRepo *repository.UserPreferencesRepository
}

func NewUserService(userRepo *repository.UserRepository, userPreferencesRepo *repository.UserPreferencesRepository) *UserService {
	return &UserService{
		UserRepo:            userRepo,
		UserPreferencesRepo: userPreferencesRepo,
	}
}

func (s *UserService) RegisterUser(req models.UserRegistrationRequest) (int, error) {
	l := slog.With("email", req.Email)
	l.Debug("Attempting to register new user")

	if req.Email == "" {
		l.Warn("User registration failed: email is missing")
		return 0, fmt.Errorf("email is required")
	}

	newUser := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}

	userID, err := s.UserRepo.Insert(newUser, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			l.Warn("User registration failed: email already taken")
			return 0, fmt.Errorf("user with email %s already exists", req.Email)
		}
		l.Error("Database error during user registration", "error", err)
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	l.Info("New user registered successfully", "user_id", userID)
	return userID, nil
}

func (s *UserService) SavePreferences(userID int, preferences models.UserPreferences) (int, error) {
	l := slog.With("user_id", userID)
	l.Debug("Attempting to save user preferences")

	if preferences.BudgetMin <= 0 || preferences.BudgetMax <= preferences.BudgetMin {
		l.Warn("Invalid budget range provided", "min", preferences.BudgetMin, "max", preferences.BudgetMax)
		return 0, fmt.Errorf("invalid budget range")
	}

	preferences.UserID = userID
	prefID, err := s.UserPreferencesRepo.Upsert(&preferences)
	if err != nil {
		l.Error("Failed to save user preferences to DB", "error", err)
		return 0, fmt.Errorf("failed to save preferences: %w", err)
	}

	l.Info("User preferences updated successfully", "pref_id", prefID)
	return prefID, nil
}

func (s *UserService) GetUserPreferences(userID int) (*models.UserPreferences, error) {
	l := slog.With("user_id", userID)

	prefs, err := s.UserPreferencesRepo.GetByUserID(userID)
	if err != nil {
		l.Error("Failed to fetch user preferences", "error", err)
		return nil, fmt.Errorf("error fetching preferences: %w", err)
	}

	l.Debug("User preferences fetched successfully")
	return prefs, nil
}
