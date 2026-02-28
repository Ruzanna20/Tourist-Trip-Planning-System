package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"math"
	"strings"
	"time"
	"travel-planning/internal/kafka"
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

	KafkaProducer *kafka.Producer
}

func NewTripPlanningService(
	tripRepo *repository.TripRepository,
	itineraryRepo *repository.TripItineraryRepository,
	itineraryActivitiesRepo *repository.ItineraryActivitiesRepository,
	flightRepo *repository.FlightRepository,
	hotelRepo *repository.HotelRepository,
	attractionRepo *repository.AttractionRepository,
	restaurantRepo *repository.RestaurantRepository,
	userPreferencesRepo *repository.UserPreferencesRepository,
	KafkaProducer *kafka.Producer) *TripPlanningService {
	return &TripPlanningService{
		TripRepo:                tripRepo,
		ItineraryRepo:           itineraryRepo,
		ItineraryActivitiesRepo: itineraryActivitiesRepo,
		FlightRepo:              flightRepo,
		HotelRepo:               hotelRepo,
		AttractionRepo:          attractionRepo,
		RestaurantRepo:          restaurantRepo,
		UserPreferencesRepo:     userPreferencesRepo,
		KafkaProducer:           KafkaProducer,
	}
}

func (s *TripPlanningService) GenerateOptions(tripID int) ([]models.TripOption, error) {
	trip, err := s.TripRepo.GetTripByID(tripID)
	if err != nil {
		return nil, fmt.Errorf("trip not found: %w", err)
	}

	prefs, err := s.UserPreferencesRepo.GetByUserID(trip.UserID)
	if err != nil {
		return nil, fmt.Errorf("could not find user home city: %v", err)
	}
	originCityID := prefs.HomeCityID

	nights := int(trip.EndDate.Sub(trip.StartDate).Hours() / 24)
	if nights <= 0 {
		return nil, fmt.Errorf("trip must be at least 1 night")
	}

	totalBudget := trip.TotalPrice
	logistics_budget := totalBudget * 0.50
	oneWayBudget := (logistics_budget * 0.6) / 2
	activities_budget := totalBudget * 0.30
	more_money := totalBudget * 0.20

	var options []models.TripOption
	tiers := []string{"Economy", "Balanced", "Luxury"}

	for _, tier := range tiers {
		outboundFlight, err := s.FlightRepo.GetBestFlightByTier(originCityID, trip.DestinationCityID, oneWayBudget, tier)
		if err != nil || outboundFlight == nil {
			log.Printf("Error outbound flight for tier %s: %v", tier, err)
			continue
		}

		inboundFlight, err := s.FlightRepo.GetBestFlightByTier(trip.DestinationCityID, originCityID, oneWayBudget, tier)
		if err != nil || inboundFlight == nil {
			log.Printf("Error inbound flight for tier %s: %v", tier, err)
			continue
		}

		totalFLightsCost := outboundFlight.Price + inboundFlight.Price
		remainingMoneyForHotel := logistics_budget - totalFLightsCost
		limitPerNight := remainingMoneyForHotel / float64(nights)

		if limitPerNight <= 0 {
			log.Printf("Skip tier %s: no money left for hotel", tier)
			continue
		}

		hotel, err := s.HotelRepo.GetBestHotelByTier(trip.DestinationCityID, limitPerNight, tier)
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
	if len(options) == 0 {
		return nil, fmt.Errorf("could not generate any trip options within your budget")
	}
	return options, nil
}

func (s *TripPlanningService) PlanTrip(userID int, req models.TripPlanRequest) (int, error) {
	l := slog.With("user_id", userID, "trip_name", req.Name)
	l.Info("Starting trip planning")

	startDate, err1 := time.Parse("2006-01-02", req.StartDate)
	endDate, err2 := time.Parse("2006-01-02", req.EndDate)
	if err1 != nil || err2 != nil {
		l.Error("Invalid date format", "start", req.StartDate, "end", req.EndDate)
		return 0, fmt.Errorf("invalid date format,use Year-month-day")
	}

	duration := endDate.Sub(startDate).Hours() / 24
	if duration <= 0 {
		l.Warn("Invalid trip duration", "duration_days", duration)
		return 0, fmt.Errorf("end date must be after start date")
	}

	tx, err := s.TripRepo.GetConn().Begin()
	if err != nil {
		l.Error("Failed to start transaction", "error", err)
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	newTrip := &models.Trip{
		UserID:            userID,
		Title:             req.Name,
		StartDate:         startDate,
		EndDate:           endDate,
		DestinationCityID: req.ToCityID,
		TotalPrice:        req.BudgetAmount,
		Status:            "Planned",
	}

	tripID, err := s.TripRepo.Insert(tx, newTrip)
	if err != nil {
		l.Error("DB error: failed to insert trip", "error", err)
		return 0, fmt.Errorf("failed to insert trip in db: %w", err)
	}

	l.Info("Trip header created", "trip_id", tripID)

	numDays := int(duration) + 1
	for day := 0; day < numDays; day++ {
		itineraryDate := startDate.AddDate(0, 0, day)
		itinerary := &models.TripItinerary{
			TripID:    tripID,
			DayNumber: day + 1,
			Date:      itineraryDate,
			Notes:     fmt.Sprintf("Plan for day %d.", day+1),
		}

		_, err := s.ItineraryRepo.Insert(tx, itinerary)
		if err != nil {
			l.Error("DB error: failed to insert itinerary", "day", day+1, "error", err)
			return 0, fmt.Errorf("critical failure during itinerary setup: %w", err)
		}

	}

	if err := tx.Commit(); err != nil {
		l.Error("Failed to commit transaction", "error", err)
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	if s.KafkaProducer != nil {
		go s.KafkaProducer.PublishTripReques(context.Background(), tripID, userID, req.ToCityID)
	}

	return tripID, nil
}

func (s *TripPlanningService) PopulateItineraryDetails(tx *sql.Tx,
	tripID int,
	cityID int,
	tier string,
	totalActivitiesBudget float64,
	days int,
	startDate time.Time,
	hotelID int,
	outboundFlightID int,
	inboundFlightID int) error {

	l := slog.With("trip_id", tripID, "tier", tier)
	l.Debug("Populating itinerary details")

	itineraries, err := s.ItineraryRepo.GetItineraryDaysByTripID(tripID)
	if err != nil || len(itineraries) == 0 {
		return fmt.Errorf("no itinerary days found for trip %d", tripID)
	}

	totalDays := len(itineraries)
	dailyBudget := totalActivitiesBudget / float64(days)
	dailyAttractionLimit := dailyBudget * 0.70

	allAttractions, err := s.AttractionRepo.GetBestAttractionsByTier(cityID, dailyAttractionLimit, tier)
	if err != nil || len(allAttractions) == 0 {
		l.Warn("No attractions found for criteria", "city_id", cityID, "limit", dailyAttractionLimit)
		return fmt.Errorf("no attractions found for city %d", cityID)
	}

	allRestaurants, err := s.RestaurantRepo.GetBestRestaurantByTier(cityID, tier)
	if err != nil || len(allRestaurants) == 0 {
		slog.Warn("No restaurants found for this city and tier",
			"trip_id", tripID, "city_id", cityID, "tier", tier)
	}

	usedAttractions := make(map[int]bool)

	for i, dayPlan := range itineraries {
		dayNum := i + 1
		currentDayID := int64(dayPlan.ItineraryID)

		switch {
		case dayNum == 1:
			s.saveActivity(tx, currentDayID, "flight", outboundFlightID, 0, allAttractions, dayPlan.Date)
			s.saveActivity(tx, currentDayID, "hotel", hotelID, 1, allAttractions, dayPlan.Date)
			s.saveActivity(tx, currentDayID, "event", 0, 2, nil, dayPlan.Date)

		case dayNum == totalDays:
			s.saveActivity(tx, currentDayID, "event", 0, 1, nil, dayPlan.Date)
			s.saveActivity(tx, currentDayID, "flight", inboundFlightID, 2, allAttractions, dayPlan.Date)

		default:
			var lastLat, lastLon float64
			var firstAttraction *models.Attraction

			for j := range allAttractions {
				if !usedAttractions[allAttractions[j].AttractionID] {
					firstAttraction = &allAttractions[j]
					usedAttractions[firstAttraction.AttractionID] = true
					s.saveActivity(tx, currentDayID, "attraction", firstAttraction.AttractionID, 1, allAttractions, dayPlan.Date)
					lastLat, lastLon = firstAttraction.Latitude, firstAttraction.Longitude
					break
				}
			}

			if firstAttraction != nil {
				var bestRest *models.Restaurant
				minD := 999999.0
				for j := range allRestaurants {
					d := calculateDistance(lastLat, lastLon, allRestaurants[j].Latitude, allRestaurants[j].Longitude)
					if d < minD {
						minD = d
						bestRest = &allRestaurants[j]
					}
				}
				if bestRest != nil {
					s.saveActivity(tx, currentDayID, "restaurant", bestRest.RestaurantID, 2, nil, dayPlan.Date)
					lastLat, lastLon = bestRest.Latitude, bestRest.Longitude
				}

				for j := range allAttractions {
					if !usedAttractions[allAttractions[j].AttractionID] && (firstAttraction.EntryFee+allAttractions[j].EntryFee <= dailyAttractionLimit) {
						usedAttractions[allAttractions[j].AttractionID] = true
						s.saveActivity(tx, currentDayID, "attraction", allAttractions[j].AttractionID, 3, allAttractions, dayPlan.Date)
						break
					}
				}
				s.saveActivity(tx, currentDayID, "hotel", hotelID, 4, nil, dayPlan.Date)
			}
		}
	}
	return nil
}

func calculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return R * c
}

func (s *TripPlanningService) saveActivity(tx *sql.Tx, itineraryID int64, aType string, entityID int, order int, allAttractions []models.Attraction, dayDate time.Time) {
	aType = strings.ToLower(aType)

	year, month, day := dayDate.Date()
	location := dayDate.Location()

	var startTime, endTime time.Time
	switch order {
	case 0:
		startTime = time.Date(year, month, day, 9, 0, 0, 0, location)
		endTime = startTime.Add(2 * time.Hour)
	case 1:
		startTime = time.Date(year, month, day, 11, 30, 0, 0, location)
		endTime = startTime.Add(1*time.Hour + 30*time.Minute)
	case 2:
		startTime = time.Date(year, month, day, 14, 0, 0, 0, location)
		endTime = startTime.Add(1 * time.Hour)
	case 3:
		startTime = time.Date(year, month, day, 16, 0, 0, 0, location)
		endTime = startTime.Add(2 * time.Hour)
	case 4:
		startTime = time.Date(year, month, day, 19, 0, 0, 0, location)
		endTime = startTime.Add(2 * time.Hour)
	default:
		startTime = time.Date(year, month, day, 21, 0, 0, 0, location)
		endTime = startTime.Add(1 * time.Hour)
	}

	activity := &models.ItineraryActivity{
		ItineraryID:  itineraryID,
		ActivityType: aType,
		OrderNumber:  order,
		StartTime:    startTime,
		EndTime:      endTime,
		AttractionID: sql.NullInt64{Int64: 0, Valid: false},
		RestaurantID: sql.NullInt64{Int64: 0, Valid: false},
		HotelID:      sql.NullInt64{Int64: 0, Valid: false},
		FlightID:     sql.NullInt64{Int64: 0, Valid: false},
	}

	switch aType {
	case "flight":
		if entityID > 0 {
			activity.FlightID = sql.NullInt64{Int64: int64(entityID), Valid: true}
			activity.Notes = "Please arrive at the airport 3 hours before departure."
		}
	case "hotel":
		if entityID > 0 {
			activity.HotelID = sql.NullInt64{Int64: int64(entityID), Valid: true}
			activity.Notes = "Accommodation check-in/stay."
		}
	case "restaurant":
		if entityID > 0 {
			activity.RestaurantID = sql.NullInt64{Int64: int64(entityID), Valid: true}
			activity.Notes = "Enjoy your meal. Local alternatives are available if busy."
		}
	case "attraction", "event":
		activity.ActivityType = "attraction"
		if entityID > 0 {
			activity.AttractionID = sql.NullInt64{Int64: int64(entityID), Valid: true}
			var lat, lon float64
			found := false

			for _, a := range allAttractions {
				if a.AttractionID == entityID {
					lat, lon = a.Latitude, a.Longitude
					found = true
					break
				}
			}
			if found {
				minD := 9999.0
				var backup string
				for _, a := range allAttractions {
					if a.AttractionID != entityID {
						d := calculateDistance(lat, lon, a.Latitude, a.Longitude)
						if d < minD {
							minD = d
							backup = a.Name
						}
					}
				}
				if backup != "" {
					activity.Notes = fmt.Sprintf("If this place is closed, the best nearby alternative is %s.", backup)
				}
			}
		} else {
			activity.AttractionID = sql.NullInt64{Valid: false}
			if order == 2 {
				activity.Notes = "Welcome. Enjoy a relaxing walk in a nearby park after your flight."
			} else {
				activity.Notes = "Last day. Perfect time for souvenir shopping and a final city stroll."
			}
		}

	}

	_, err := s.ItineraryActivitiesRepo.Insert(tx, activity)
	if err != nil {
		slog.Error("Failed to insert itinerary activity",
			"itinerary_id", itineraryID,
			"activity_type", aType,
			"order", order,
			"error", err)
	}
}

func (s *TripPlanningService) FinalizeTripPlan(tripID int, tier string, hotelID int, outboundFlightID int, inboundFlightID int) error {
	l := slog.With("trip_id", tripID, "tier", tier)
	l.Info("Finalizing trip plan")

	trip, err := s.TripRepo.GetTripByID(tripID)
	if err != nil {
		l.Error("Trip not found", "error", err)
		return err
	}

	tx, err := s.TripRepo.GetConn().Begin()
	if err != nil {
		l.Error("Transaction failed to start", "error", err)
		return err
	}
	defer tx.Rollback()

	activitiesBudget := trip.TotalPrice * 0.30
	DurationDays := int(trip.EndDate.Sub(trip.StartDate).Hours()/24) + 1
	err = s.PopulateItineraryDetails(
		tx,
		trip.TripID,
		trip.DestinationCityID,
		tier,
		activitiesBudget,
		DurationDays,
		trip.StartDate,
		hotelID,
		outboundFlightID,
		inboundFlightID,
	)
	if err != nil {
		l.Error("Failed to populate details", "error", err)
		return err
	}

	if err := tx.Commit(); err != nil {
		l.Error("Finalize commit failed", "error", err)
		return err
	}

	l.Info("Trip finalized and saved successfully")
	return nil
}

func (s *TripPlanningService) GetTripByID(id int) (*models.Trip, error) {
	return s.TripRepo.GetTripByID(id)
}

func (s *TripPlanningService) GetUserTrips(userID int) ([]models.Trip, error) {
	return s.TripRepo.GetAllTripsByUserID(userID)
}

func (s *TripPlanningService) GetItineraryDays(tripID int) ([]*models.TripItinerary, error) {
	return s.ItineraryRepo.GetItineraryDaysByTripID(tripID)
}

func (s *TripPlanningService) GetActivitiesByDay(itineraryID int) ([]*models.ItineraryActivity, error) {
	return s.ItineraryActivitiesRepo.GetActivitiesByItineraryID(itineraryID)
}

func (s *TripPlanningService) DeleteUserTrip(tripID, userID int) error {
	return s.TripRepo.DeleteByIDAndUserID(tripID, userID)
}
