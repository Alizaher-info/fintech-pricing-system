package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IBM/sarama"
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
	producer sarama.SyncProducer
	topic    string
}

// NewProducer creates a new Kafka producer using Sarama
func NewProducer(brokers []string, topic string) (*Producer, error) {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal // Wait for leader acknowledgment
	config.Producer.Retry.Max = 3
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner // Hash by key for consistent partitioning

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	}

	logger.Info(fmt.Sprintf("Kafka producer initialized - Brokers: %v, Topic: %s", brokers, topic))
	
	return &Producer{
		producer: producer,
		topic:    topic,
	}, nil
}

// PublishPriceUpdate sends a price update message to Kafka
func (p *Producer) PublishPriceUpdate(ctx context.Context, msg PriceUpdateMessage) error {
	// Convert message to JSON
	messageBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Create Kafka message
	kafkaMsg := &sarama.ProducerMessage{
		Topic: p.topic,
		Key:   sarama.StringEncoder(msg.Symbol), // Use symbol as key for partitioning
		Value: sarama.ByteEncoder(messageBytes),
	}

	// Publish message synchronously
	partition, offset, err := p.producer.SendMessage(kafkaMsg)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	logger.Info(fmt.Sprintf("Published price update to Kafka - %s: $%.2f (partition: %d, offset: %d)", msg.Symbol, msg.Price, partition, offset))
	return nil
}

// Close closes the Kafka producer
func (p *Producer) Close() error {
	if p.producer != nil {
		logger.Info("Closing Kafka producer...")
		return p.producer.Close()
	}
	return nil
}
