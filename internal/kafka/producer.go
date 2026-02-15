package kafka

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(broker string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) PublishTripReques(ctx context.Context, tripID int, userID int, cityID int) error {
	event := map[string]interface{}{
		"type":       "generate_itinerary",
		"trip_id":    tripID,
		"user_id":    userID,
		"city_id":    cityID,
		"created_at": time.Now(),
	}

	msgBytes, err := json.Marshal(event)
	if err != nil {
		slog.Error("Failed to marshal trip event", "error", err)
		return err
	}

	err = p.writer.WriteMessages(ctx, kafka.Message{
		Topic: "trip-requests",
		Key:   []byte(strconv.Itoa(tripID)),
		Value: msgBytes,
	})

	if err != nil {
		slog.Error("Kafka publish failed", "trip_id", tripID, "error", err)
		return err
	}

	slog.Info("Event sent to Kafka", "trip_id", tripID, "topic", "trip-requests")
	return nil
}

func (p *Producer) Close() {
	p.writer.Close()
}
