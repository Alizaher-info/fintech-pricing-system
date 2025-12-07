package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
	"services/pricing-fetcher/pkg/logger"
)

// PriceUpdateMessage represents a price update event
type PriceUpdateMessage struct {
	Symbol      string    `json:"symbol"`
	Price       float64   `json:"price"`
	Change24h   float64   `json:"change_24h"`
	Volume24h   float64   `json:"volume_24h"`
	High24h     float64   `json:"high_24h"`
	Low24h      float64   `json:"low_24h"`
	MarketCap   float64   `json:"market_cap"`
	Timestamp   time.Time `json:"timestamp"`
}

// Producer handles publishing messages to Kafka
type Producer struct {
	writer *kafka.Writer
}

// NewProducer creates a new Kafka producer
func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		Async:        false, // Synchronous writes for reliability
	}

	logger.Info(fmt.Sprintf("Kafka producer initialized - Brokers: %v, Topic: %s", brokers, topic))
	
	return &Producer{
		writer: writer,
	}
}

// PublishPriceUpdate sends a price update message to Kafka
func (p *Producer) PublishPriceUpdate(ctx context.Context, msg PriceUpdateMessage) error {
	// Convert message to JSON
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Create Kafka message
	kafkaMsg := kafka.Message{
		Key:   []byte(msg.Symbol), // Use symbol as key for partitioning
		Value: messageBytes,
		Time:  time.Now(),
	}

	// Publish message
	err = p.writer.WriteMessages(ctx, kafkaMsg)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	logger.Info(fmt.Sprintf("Published price update to Kafka - %s: $%.2f", msg.Symbol, msg.Price))
	return nil
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	if p.writer != nil {
		logger.Info("Closing Kafka producer...")
		return p.writer.Close()
	}
	return nil
}
