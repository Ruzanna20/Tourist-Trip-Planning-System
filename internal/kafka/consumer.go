package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"travel-planning/internal/notifications"
	"travel-planning/models"

	"github.com/segmentio/kafka-go"
)

type TripProcessor interface {
	GetTripByID(id int) (*models.Trip, error)
	GenerateOptions(tripID int) ([]models.TripOption, error)
	FinalizeTripPlan(tripID int, tier string, hotelID int, outboundID int, inboundID int) error
	SaveNotification(userID int, tripID int, message string, msgType string) error
}

type Consumer struct {
	reader  *kafka.Reader
	service TripProcessor
	hub     *notifications.Hub
}

func NewConsumer(brokers []string, topic string, groupID string, service TripProcessor, hub *notifications.Hub) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		service: service,
		hub:     hub,
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

		tripID := int(event["trip_id"].(float64))
		userID := int(event["user_id"].(float64))

		slog.Info("Consumer picked up trip request", "trip_id", tripID)

		_, err = c.service.GenerateOptions(tripID)
		if err == nil {
			msg := "Ձեր ուղևորության տարբերակները պատրաստ են։ Նայեք մանրամասները այստեղ։"

			errNotify := c.service.SaveNotification(userID, tripID, msg, "TRIP_READY")
			if errNotify != nil {
				slog.Error("Failed to persist notification", "error", errNotify)
			}

			c.hub.SendNotification(userID, msg, tripID)
			slog.Info("Real-time notification saved and sent", "user_id", userID)
		}
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
