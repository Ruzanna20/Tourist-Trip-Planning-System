package services

import (
	"fmt"
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
	if req.Email == "" {
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
			return 0, fmt.Errorf("user with email %s already exists", req.Email)
		}
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userID, nil
}

func (s *UserService) SavePreferences(userID int, preferences models.UserPreferences) (int, error) {
	if preferences.BudgetMin <= 0 || preferences.BudgetMax <= preferences.BudgetMin {
		return 0, fmt.Errorf("invalid budget range")
	}

	preferences.UserID = userID
	prefID, err := s.UserPreferencesRepo.Upsert(&preferences)
	if err != nil {
		return 0, fmt.Errorf("failed to save preferences: %w", err)
	}

	return prefID, nil
}

func (s *UserService) GetUserPreferences(userID int) (*models.UserPreferences, error) {
	prefs, err := s.UserPreferencesRepo.GetByUserID(userID)
	if err != nil {
		return nil, fmt.Errorf("error fetching preferences: %w", err)
	}
	return prefs, nil
}
