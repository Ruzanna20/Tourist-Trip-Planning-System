package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"
	"travel-planning/internal/notifications"
	"travel-planning/models"

	"github.com/segmentio/kafka-go"
)

type TripProcessor interface {
	GetTripByID(id int) (*models.Trip, error)
	GenerateOptions(tripID int) ([]models.TripOption, error)
	FinalizeTripPlan(tripID int, tier string, hotelID int, outboundID int, inboundID int) error
}

type NotificationProcessor interface {
	CreateNotification(userID int, tripID int, message string, msgType string) (int, error)
}

type Consumer struct {
	reader       *kafka.Reader
	tripService  TripProcessor
	notifService NotificationProcessor
	hub          *notifications.Hub
}

func NewConsumer(
	brokers []string,
	topic string,
	groupID string,
	tripService TripProcessor,
	notifService NotificationProcessor,
	hub *notifications.Hub,
) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		tripService:  tripService,
		notifService: notifService,
		hub:          hub,
	}
}

func (c *Consumer) Start(ctx context.Context) {
	slog.Info("Kafka Consumer started", "topic", c.reader.Config().Topic)

	for {
		m, err := c.reader.ReadMessage(ctx)
		if err != nil {
			slog.Error("Failed to read message from Kafka", "error", err)
			break
		}

		var event map[string]interface{}
		if err := json.Unmarshal(m.Value, &event); err != nil {
			slog.Error("Failed to unmarshal event", "error", err)
			continue
		}

		tripIDVal, ok1 := event["trip_id"].(float64)
		userIDVal, ok2 := event["user_id"].(float64)
		if !ok1 || !ok2 {
			slog.Error("Invalid event payload: missing trip_id or user_id")
			continue
		}

		tripID := int(tripIDVal)
		userID := int(userIDVal)

		slog.Info("Consumer picked up trip request", "trip_id", tripID)

		slog.Debug("Attempting to generate options...", "trip_id", tripID)
		_, err = c.tripService.GenerateOptions(tripID)
		if err == nil {
			msg := "Your trip options are ready. Check the details here."

			slog.Info("Options generated, saving notification...", "trip_id", tripID)
			notificationID, errNotify := c.notifService.CreateNotification(userID, tripID, msg, "TRIP_READY")

			if errNotify == nil {
				payload := map[string]interface{}{
					"id":         notificationID,
					"message":    msg,
					"trip_id":    tripID,
					"type":       "TRIP_READY",
					"is_read":    false,
					"created_at": time.Now().Format(time.RFC3339),
				}
				data, _ := json.Marshal(payload)
				c.hub.SendNotification(userID, string(data), tripID)
			} else {
				slog.Error("Failed to save notification via NotificationService", "error", errNotify)
			}
		} else {
			slog.Error("Failed to generate trip options", "trip_id", tripID, "error", err)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
