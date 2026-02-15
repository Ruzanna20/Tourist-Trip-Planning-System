package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"travel-planning/models"

	"github.com/segmentio/kafka-go"
)

type TripProcessor interface {
	GetTripByID(id int) (*models.Trip, error)
	GenerateOptions(trip *models.Trip) ([]models.TripOption, error)
	FinalizeTripPlan(tripID int, tier string, hotelID int, outboundID int, inboundID int) error
}

type Consumer struct {
	reader  *kafka.Reader
	service TripProcessor
}

func NewConsumer(brokers []string, topic string, groupID string, service TripProcessor) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   topic,
			GroupID: groupID,
		}),
		service: service,
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
		slog.Info("Consumer picked up trip request", "trip_id", tripID)

		trip, err := c.service.GetTripByID(tripID)
		if err != nil {
			slog.Error("Failed to fetch trip from DB", "trip_id", tripID, "error", err)
			continue
		}

		options, err := c.service.GenerateOptions(trip)
		if err != nil {
			slog.Error("Failed to generate options", "trip_id", tripID, "error", err)
			continue
		}
		best := options[0]

		err = c.service.FinalizeTripPlan(
			tripID,
			best.Tier,
			best.Hotel.HotelID,
			best.OutBoundFlight.FlightID,
			best.InBoundFlight.FlightID,
		)

		if err != nil {
			slog.Error("Failed to finalize trip plan via Kafka worker", "trip_id", tripID, "error", err)
		} else {
			slog.Info("Successfully processed trip plan from Kafka", "trip_id", tripID)
		}

	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
