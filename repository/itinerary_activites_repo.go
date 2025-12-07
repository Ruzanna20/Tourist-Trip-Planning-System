package repository

import (
	"database/sql"
	"fmt"
	"time"
	"travel-planning/models"
)

type ItineraryActivitiesRepository struct {
	db *sql.DB
}

func NewItineraryActivitiesRepository(db *sql.DB) *ItineraryActivitiesRepository {
	return &ItineraryActivitiesRepository{db: db}
}

func (r *ItineraryActivitiesRepository) Insert(activity *models.ItineraryActivity) (int64, error) {
	query := `INSERT INTO itinerary_activities (itinerary_id, activity_type, hotel_id, attraction_id, restaurant_id, 
	flight_id, order_number, start_time, end_time, notes, created_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING activity_id;`

	var activityID int64
	currTime := time.Now()
	err := r.db.QueryRow(
		query,
		activity.ItineraryID,
		activity.ActivityType,
		activity.HotelID,
		activity.AttractionID,
		activity.RestaurantID,
		activity.FlightID,
		activity.OrderNumber,
		activity.StartTime,
		activity.EndTime,
		activity.Notes,
		currTime,
	).Scan(&activityID)

	if err != nil {
		return 0, fmt.Errorf("failed to insert itinerary activity: %w", err)
	}

	return activityID, nil
}

func (r *ItineraryActivitiesRepository) GetActivitiesByItineraryID(itineraryID int) ([]*models.ItineraryActivity, error) {
	query := `SELECT 
	            activity_id, itinerary_id, activity_type, hotel_id, attraction_id, 
	            restaurant_id, flight_id, order_number, start_time, end_time, 
	            notes, created_at
	          FROM itinerary_activities 
	          WHERE itinerary_id = $1`

	rows, err := r.db.Query(
		query,
		itineraryID)
	if err != nil {
		return nil, fmt.Errorf("error querying itinerary activities: %w", err)
	}

	defer rows.Close()

	var activities []*models.ItineraryActivity
	for rows.Next() {
		activity := &models.ItineraryActivity{}

		err := rows.Scan(
			&activity.ActivityID,
			&activity.ItineraryID,
			&activity.ActivityType,
			&activity.HotelID,
			&activity.AttractionID,
			&activity.RestaurantID,
			&activity.FlightID,
			&activity.OrderNumber,
			&activity.StartTime,
			&activity.EndTime,
			&activity.Notes,
			&activity.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("error scanning itinerary activity:%w", err)
		}
		activities = append(activities, activity)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error:%w", err)
	}

	return activities, nil
}
