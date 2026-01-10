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

	FlightRepo     *repository.FlightRepository
	HotelRepo      *repository.HotelRepository
	AttractionRepo *repository.AttractionRepository
	RestaurantRepo *repository.RestaurantRepository

	UserPreferencesRepo *repository.UserPreferencesRepository
}

func NewTripPlanningService(
	tripRepo *repository.TripRepository,
	itineraryRepo *repository.TripItineraryRepository,
	itineraryActivitiesRepo *repository.ItineraryActivitiesRepository,
	flightRepo *repository.FlightRepository,
	hotelRepo *repository.HotelRepository,
	attractionRepo *repository.AttractionRepository,
	restaurantRepo *repository.RestaurantRepository,
	userPreferencesRepo *repository.UserPreferencesRepository) *TripPlanningService {
	return &TripPlanningService{
		TripRepo:                tripRepo,
		ItineraryRepo:           itineraryRepo,
		ItineraryActivitiesRepo: itineraryActivitiesRepo,
		FlightRepo:              flightRepo,
		HotelRepo:               hotelRepo,
		AttractionRepo:          attractionRepo,
		RestaurantRepo:          restaurantRepo,
		UserPreferencesRepo:     userPreferencesRepo,
	}
}

func (s *TripPlanningService) GenerateTripOptions(userID int, req models.TripPlanRequest) ([]models.TripOption, error) {
	startDate, err1 := time.Parse("2006-01-02", req.StartDate)
	endDate, err2 := time.Parse("2006-01-02", req.EndDate)
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("invalid date format,use Year-month-day")
	}

	nights := int(endDate.Sub(startDate).Hours() / 24)
	if nights <= 0 {
		return nil, fmt.Errorf("trip must be at least 1 night")
	}

	totalBudget := req.BudgetAmount
	logistics_budget := totalBudget * 0.50
	oneWayBudget := (logistics_budget * 0.6) / 2
	activities_budget := totalBudget * 0.30
	more_money := totalBudget * 0.20

	var options []models.TripOption
	tiers := []string{"Economy", "Balanced", "Luxury"}

	for _, tier := range tiers {
		outboundFlight, err := s.FlightRepo.GetBestFlightByTier(req.FromCityID, req.ToCityID, oneWayBudget, tier)
		if err != nil || outboundFlight == nil {
			log.Printf("Error outbound flight for tier %s: %v", tier, err)
			continue
		}

		inboundFlight, err := s.FlightRepo.GetBestFlightByTier(req.ToCityID, req.FromCityID, oneWayBudget, tier)
		if err != nil || inboundFlight == nil {
			log.Printf("Error inbound flight for tier %s: %v", tier, err)
			continue
		}

		totalFLightsCost := outboundFlight.Price + inboundFlight.Price
		remainingMoneyForHotel := logistics_budget - totalFLightsCost
		limitPerNight := remainingMoneyForHotel / float64(nights)

		hotel, err := s.HotelRepo.GetBestHotelByTier(req.ToCityID, limitPerNight, tier)
		if err != nil {
			log.Printf("Error finding hotel for tier %s: %v", tier, err)
			continue
		}

		hotelCost := hotel.PricePerNight * float64(nights)
		actualLogisticsCost := totalFLightsCost + hotelCost

		option := models.TripOption{
			Tier:             tier,
			OutBoundFlight:   outboundFlight,
			InBoundFlight:    inboundFlight,
			Hotel:            hotel,
			LogisticsBudget:  actualLogisticsCost,
			ActivitiesBudget: activities_budget,
			MoreMoney:        more_money + (logistics_budget - actualLogisticsCost),
			TotalPriceOfTrip: actualLogisticsCost + activities_budget + more_money,
		}
		options = append(options, option)
	}
	return options, nil
}

func (s *TripPlanningService) PlanTrip(userID int, req models.TripPlanRequest) (int, error) {
	startDate, err1 := time.Parse("2006-01-02", req.StartDate)
	endDate, err2 := time.Parse("2006-01-02", req.EndDate)
	if err1 != nil || err2 != nil {
		return 0, fmt.Errorf("invalid date format,use Year-month-day")
	}

	newTrip := &models.Trip{
		UserID:            userID,
		Title:             req.Name,
		StartDate:         startDate,
		EndDate:           endDate,
		DestinationCityID: req.ToCityID,
		TotalPrice:        req.BudgetAmount,
		Currency:          req.Currency,
		Status:            "Planned",
	}

	tripID, err := s.TripRepo.Insert(newTrip)
	if err != nil {
		return 0, fmt.Errorf("failed to insert trip in db: %w", err)
	}

	duration := endDate.Sub(startDate).Hours() / 24
	if duration <= 0 {
		return 0, fmt.Errorf("end date must be after start date")
	}
	numDays := int(duration) + 1
	for day := 0; day <= numDays; day++ {
		itineraryDate := startDate.AddDate(0, 0, day)
		itinerary := &models.TripItinerary{
			TripID:    tripID,
			DayNumber: day + 1,
			Date:      itineraryDate,
			Notes:     fmt.Sprintf("Plan for day %d.", day+1),
		}

		_, err := s.ItineraryRepo.Insert(itinerary)
		if err != nil {
			log.Printf("CRITICAL: Failed to insert itinerary for Trip %d, Day %d: %v", tripID, day+1, err)
			return 0, fmt.Errorf("critical failure during itinerary setup: %w", err)
		}

	}
	return tripID, nil
}
