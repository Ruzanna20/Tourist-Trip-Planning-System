package services

import (
	"fmt"
	"log"
	"time"
	"travel-planning/models"
	"travel-planning/repository"
)

type TripPlanningService struct {
	TripRepo                *repository.TripRepository
	ItineraryRepo           *repository.TripItineraryRepository
	ItineraryActivitiesRepo *repository.ItineraryActivitiesRepository
}

func NewTripPlanningService(tripRepo *repository.TripRepository, itineraryRepo *repository.TripItineraryRepository,
	ItineraryActivitiesRepo *repository.ItineraryActivitiesRepository) *TripPlanningService {
	return &TripPlanningService{
		TripRepo:                tripRepo,
		ItineraryRepo:           itineraryRepo,
		ItineraryActivitiesRepo: ItineraryActivitiesRepo,
	}
}

func (s *TripPlanningService) PlanTrip(userID int, req models.TripPlanRequest) (int, error) {
	startDate, err1 := time.Parse("2006-01-02", req.StartDate)
	endDate, err2 := time.Parse("2006-01-02", req.EndDate)
	if err1 != nil || err2 != nil {
		return 0, fmt.Errorf("invalid date format,use Year-month-day")
	}

	duration := endDate.Sub(startDate).Hours() / 24
	if duration <= 0 {
		return 0, fmt.Errorf("end date must be after start date")
	}
	newTrip := &models.Trip{
		UserID:            userID,
		Title:             req.Name,
		StartDate:         startDate,
		EndDate:           endDate,
		DestinationCityID: req.DestinationCityID,
		TotalPrice:        req.BudgetAmount,
		Currency:          req.Currency,
		Status:            "Planned",
	}

	tripID, err := s.TripRepo.Insert(newTrip)
	if err != nil {
		return 0, fmt.Errorf("failed to insert trip in db: %w", err)
	}

	numDays := int(duration) + 1
	for day := 0; day < numDays; day++ {
		itineraryDate := startDate.AddDate(0, 0, day)
		itinerary := &models.TripItinerary{
			TripID:    tripID,
			DayNumber: day + 1,
			Date:      itineraryDate,
			Notes:     fmt.Sprintf("Day %d itinerary notes.", day+1),
		}

		_, err := s.ItineraryRepo.Insert(itinerary)
		if err != nil {
			log.Printf("CRITICAL: Failed to insert itinerary for Trip %d, Day %d: %v", tripID, day+1, err)
			return 0, fmt.Errorf("critical failure during itinerary setup: %w", err)
		}

	}
	return tripID, nil
}
