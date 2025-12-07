package kafka

import (
	"context"
	"fmt"
	"sync"

	"services/pricing-fetcher/pkg/logger"
)

// ProducerManager manages multiple Kafka producers for different topics
type ProducerManager struct {
	producers map[string]*Producer
	brokers   []string
	mu        sync.RWMutex // Thread-safe access to producers map
}

// NewProducerManager creates a new producer manager
func NewProducerManager(brokers []string) *ProducerManager {
	logger.Info(fmt.Sprintf("Initializing Kafka Producer Manager - Brokers: %v", brokers))
	
	return &ProducerManager{
		producers: make(map[string]*Producer),
		brokers:   brokers,
	}
}

// GetProducer returns existing producer for topic or creates new one
func (pm *ProducerManager) GetProducer(topic string) *Producer {
	pm.mu.RLock()
	producer, exists := pm.producers[topic]
	pm.mu.RUnlock()
	
	if exists {
		return producer
	}
	
	// Create new producer for this topic
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	// Double-check in case another goroutine created it
	if producer, exists := pm.producers[topic]; exists {
		return producer
	}
	
	producer = NewProducer(pm.brokers, topic)
	pm.producers[topic] = producer
	
	logger.Info(fmt.Sprintf("Created new Kafka producer for topic: %s", topic))
	return producer
}

// PublishPriceUpdate publishes a price update message to specified topic
func (pm *ProducerManager) PublishPriceUpdate(ctx context.Context, topic string, msg PriceUpdateMessage) error {
	producer := pm.GetProducer(topic)
	return producer.PublishPriceUpdate(ctx, msg)
}

// PublishMessage publishes any message to specified topic (generic)
func (pm *ProducerManager) PublishMessage(ctx context.Context, topic string, key string, value interface{}) error {
	_ = pm.GetProducer(topic)  // Get producer but not used yet
	
	// Use the producer's publish method
	// For now, we'll use PublishPriceUpdate as base
	// You can extend this for other message types later
	
	return fmt.Errorf("generic publish not yet implemented - use PublishPriceUpdate")
}

// GetActiveTopics returns list of topics that have active producers
func (pm *ProducerManager) GetActiveTopics() []string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	
	topics := make([]string, 0, len(pm.producers))
	for topic := range pm.producers {
		topics = append(topics, topic)
	}
	
	return topics
}

// CloseAll closes all active producers
func (pm *ProducerManager) CloseAll() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	
	logger.Info(fmt.Sprintf("Closing all Kafka producers (%d active topics)", len(pm.producers)))
	
	var lastErr error
	for topic, producer := range pm.producers {
		logger.Info(fmt.Sprintf("Closing producer for topic: %s", topic))
		if err := producer.Close(); err != nil {
			logger.Error(fmt.Sprintf("Failed to close producer for topic %s: %v", topic, err))
			lastErr = err
		}
	}
	
	// Clear the map
	pm.producers = make(map[string]*Producer)
	
	if lastErr != nil {
		return fmt.Errorf("some producers failed to close: %w", lastErr)
	}
	
	logger.Info("All Kafka producers closed successfully")
	return nil
}

// GetProducerCount returns number of active producers
func (pm *ProducerManager) GetProducerCount() int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return len(pm.producers)
}
