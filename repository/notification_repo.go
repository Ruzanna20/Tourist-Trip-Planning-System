package repository

import (
	"database/sql"
	"fmt"
	"travel-planning/models"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{
		db: db,
	}
}

func (r *NotificationRepository) SaveNotification(userID int, tripID int, message string, msgType string) (int, error) {
	var id int
	query := `
		INSERT INTO notifications (user_id, trip_id, message, type, is_read, created_at)
		VALUES ($1, $2, $3, $4, false, NOW())
		RETURNING id
	`
	err := r.db.QueryRow(query, userID, tripID, message, msgType).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("failed to save notification: %w", err)
	}
	return id, nil
}

func (r *NotificationRepository) GetUserNotifications(userID int) ([]models.Notification, error) {
	query := `
        SELECT 
            n.id, n.user_id, n.trip_id, n.message, n.type, n.is_read, n.created_at,
            COALESCE(t.title, '') as trip_title
        FROM notifications n
        LEFT JOIN trips t ON n.trip_id = t.trip_id
        WHERE n.user_id = $1 
        ORDER BY n.created_at DESC LIMIT 50`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer rows.Close()

	var notifs []models.Notification
	for rows.Next() {
		var n models.Notification
		err := rows.Scan(
			&n.ID,
			&n.UserID,
			&n.TripID,
			&n.Message,
			&n.Type,
			&n.IsRead,
			&n.CreatedAt,
			&n.TripTitle,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifs = append(notifs, n)
	}

	return notifs, nil
}

func (r *NotificationRepository) MarkAsRead(notificationID int, userID int) error {
	query := `UPDATE notifications SET is_read = true WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(query, notificationID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}
	return nil
}

func (r *NotificationRepository) GetUnreadCount(userID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = false`
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count unread notifications: %w", err)
	}
	return count, nil
}

func (r *NotificationRepository) DeleteNotification(notificationID int, userID int) error {
	query := `DELETE FROM notifications WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(query, notificationID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("notification not found or already deleted")
	}

	return nil
}
