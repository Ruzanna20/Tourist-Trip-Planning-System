package services

import (
	"database/sql"
	"fmt"
	"log"
	"math"
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

func (s *TripPlanningService) GenerateOptions(trip *models.Trip) ([]models.TripOption, error) {
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

func (s *TripPlanningService) PopulateItineraryDetails(tripID int, cityID int, tier string, totalActivitiesBudget float64, days int, startDate time.Time) error {
	dailyBudget := totalActivitiesBudget / float64(days)
	dailyAttractionLimit := dailyBudget * 0.70

	allAttractions, err := s.AttractionRepo.GetBestAttractionsByTier(cityID, dailyAttractionLimit, tier)
	if err != nil || len(allAttractions) == 0 {
		return fmt.Errorf("no attractions found for city %d", cityID)
	}

	allRestaurants, err := s.RestaurantRepo.GetBestRestaurantByTier(cityID, tier)
	if err != nil || len(allRestaurants) == 0 {
		return fmt.Errorf("no restaurants found")
	}

	usedAttractions := make(map[int]bool)
	itineraries, _ := s.ItineraryRepo.GetItineraryDaysByTripID(tripID)

	for day := 1; day <= days; day++ {
		if len(itineraries) < day {
			continue
		}
		currentDayItineraryID := itineraries[day-1].ItineraryID

		var firstAttraction *models.Attraction
		for i := range allAttractions {
			if !usedAttractions[allAttractions[i].AttractionID] && allAttractions[i].EntryFee <= dailyAttractionLimit*0.5 {
				firstAttraction = &allAttractions[i]
				usedAttractions[firstAttraction.AttractionID] = true
				break
			}
		}

		if firstAttraction == nil {
			continue
		}

		var bestRestaurant *models.Restaurant
		minDistanceRest := 999999.9
		for i := range allRestaurants {
			d := calculateDistance(firstAttraction.Latitude, firstAttraction.Longitude, allRestaurants[i].Latitude, allRestaurants[i].Longitude)
			if d < minDistanceRest {
				minDistanceRest = d
				bestRestaurant = &allRestaurants[i]
			}
		}

		var secondAttraction *models.Attraction
		minDistanceAttr := 999999.9
		remainingAttrBudget := dailyAttractionLimit - firstAttraction.EntryFee

		for i := range allAttractions {
			if usedAttractions[allAttractions[i].AttractionID] || allAttractions[i].EntryFee > remainingAttrBudget {
				continue
			}
			d := calculateDistance(bestRestaurant.Latitude, bestRestaurant.Longitude, allAttractions[i].Latitude, allAttractions[i].Longitude)
			if d < minDistanceAttr {
				minDistanceAttr = d
				secondAttraction = &allAttractions[i]
			}

			if secondAttraction != nil {
				usedAttractions[secondAttraction.AttractionID] = true
			}

			s.saveActivity(int64(currentDayItineraryID), "Attraction", firstAttraction.AttractionID, 1, allAttractions)
			if bestRestaurant != nil {
				s.saveActivity(int64(currentDayItineraryID), "Restaurant", bestRestaurant.RestaurantID, 2, allAttractions)
			}
			if secondAttraction != nil {
				s.saveActivity(int64(currentDayItineraryID), "Attraction", secondAttraction.AttractionID, 3, allAttractions)
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

func (s *TripPlanningService) saveActivity(itineraryID int64, aType string, entityID int, order int, allAttractions []models.Attraction) {
	activity := &models.ItineraryActivity{
		ItineraryID:  itineraryID,
		ActivityType: aType,
		OrderNumber:  order,
	}
	if aType == "Attraction" {
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
			var backupName string
			minDist := 5.0
			for _, a := range allAttractions {
				if a.AttractionID != entityID {
					d := calculateDistance(lat, lon, a.Latitude, a.Longitude)
					if d < minDist {
						minDist = d
						backupName = a.Name
					}
				}
			}

			if backupName != "" {
				activity.Notes = fmt.Sprintf("Backup plan: Visit %s if the primary location is unavailable.", backupName)
			}
		} else {
			activity.Notes = "Backup plan: Explore the local historic district or nearby public parks."
		}
	} else if aType == "Restaurant" {
		activity.RestaurantID = sql.NullInt64{Int64: int64(entityID), Valid: true}
	}
	s.ItineraryActivitiesRepo.Insert(activity)
}

func (s *TripPlanningService) FinalizeTripPlan(tripID int, tier string) error {
	trip, err := s.TripRepo.GetTripByID(tripID)
	if err != nil {
		return err
	}

	activitiesBudget := trip.TotalPrice * 0.30
	DurationDays := int(trip.EndDate.Sub(trip.StartDate).Hours() / 24)
	return s.PopulateItineraryDetails(
		trip.TripID,
		trip.DestinationCityID,
		tier,
		activitiesBudget,
		DurationDays,
		trip.StartDate,
	)
}
