package repository

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"
	"travel-planning/models"
)

type ItineraryActivitiesRepository struct {
	db *sql.DB
}

func NewItineraryActivitiesRepository(db *sql.DB) *ItineraryActivitiesRepository {
	return &ItineraryActivitiesRepository{db: db}
}

func (r *ItineraryActivitiesRepository) Insert(tx *sql.Tx, activity *models.ItineraryActivity) (int64, error) {
	query := `INSERT INTO itinerary_activities (itinerary_id, activity_type, hotel_id, attraction_id, restaurant_id, 
	flight_id, order_number, start_time, end_time, notes, created_at
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING activity_id;`

	var activityID int64
	currTime := time.Now()
	err := tx.QueryRow(
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
		slog.Error("Failed to insert itinerary activity",
			"itinerary_id", activity.ItineraryID,
			"activity_type", activity.ActivityType,
			"order", activity.OrderNumber,
			"error", err,
		)
		return 0, fmt.Errorf("failed to insert itinerary activity: %w", err)
	}

	slog.Debug("Itinerary activity inserted successfully",
		"activity_id", activityID,
		"type", activity.ActivityType,
	)
	return activityID, nil
}

func (r *ItineraryActivitiesRepository) GetActivitiesByItineraryID(itineraryID int) ([]*models.ItineraryActivity, error) {
	query := `
		SELECT
			ia.activity_id,
            ia.itinerary_id,
            ia.activity_type,
            ia.hotel_id,
            ia.attraction_id,
            ia.restaurant_id,
            ia.flight_id,
            ia.order_number,
            ia.start_time,
            ia.end_time,
            ia.notes,
            ia.created_at,
			CASE
				WHEN ia.activity_type = 'hotel'      THEN COALESCE(h.name, '')
				WHEN ia.activity_type = 'attraction' THEN COALESCE(a.name, '')
				WHEN ia.activity_type = 'restaurant' THEN COALESCE(r.name, '')
				WHEN ia.activity_type = 'flight'     THEN COALESCE(f.airline, '')
				ELSE ''
			END AS entity_name,
			CASE
				WHEN ia.activity_type = 'hotel'      THEN COALESCE(h.address, '')
				WHEN ia.activity_type = 'attraction' THEN COALESCE(a.category, '')
				WHEN ia.activity_type = 'restaurant' THEN COALESCE(r.cuisine, '')
				WHEN ia.activity_type = 'flight'     THEN COALESCE(CAST(f.duration_minutes AS TEXT) || ' min', '')
				ELSE ''
			END AS entity_detail,
			CASE
				WHEN ia.activity_type = 'hotel'      THEN COALESCE(h.description, '')
				WHEN ia.activity_type = 'attraction' THEN COALESCE(CAST(a.entry_fee AS TEXT), '')
				WHEN ia.activity_type = 'restaurant' THEN COALESCE(r.price_range, '')
				WHEN ia.activity_type = 'flight'     THEN COALESCE(CAST(f.price AS TEXT), '')
				ELSE ''
			END AS entity_extra,
			CASE
				WHEN ia.activity_type = 'hotel'      THEN COALESCE(h.rating, 0)
				WHEN ia.activity_type = 'attraction' THEN COALESCE(a.rating, 0)
				WHEN ia.activity_type = 'restaurant' THEN COALESCE(r.rating, 0)
				ELSE 0
			END AS entity_rating
		FROM itinerary_activities ia
		LEFT JOIN hotels      h ON ia.hotel_id      = h.hotel_id
		LEFT JOIN attractions a ON ia.attraction_id = a.attraction_id
		LEFT JOIN restaurants r ON ia.restaurant_id = r.restaurant_id
		LEFT JOIN flights     f ON ia.flight_id     = f.flight_id
		WHERE ia.itinerary_id = $1
		ORDER BY ia.order_number ASC`

	rows, err := r.db.Query(query, itineraryID)
	if err != nil {
		slog.Error("Error querying itinerary activities", "itinerary_id", itineraryID, "error", err)
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
			&activity.EntityName,
			&activity.EntityDetail,
			&activity.EntityExtra,
			&activity.EntityRating,
			&activity.CreatedAt,
		)

		if err != nil {
			slog.Warn("Error scanning itinerary activity row", "error", err)
			continue
		}
		activities = append(activities, activity)
	}

	if err = rows.Err(); err != nil {
		slog.Error("Rows iteration error in activities fetching", "error", err)
		return nil, fmt.Errorf("rows iteration error:%w", err)
	}

	slog.Debug("Fetched itinerary activities", "count", len(activities), "itinerary_id", itineraryID)
	return activities, nil
}
