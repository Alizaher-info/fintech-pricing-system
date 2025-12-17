package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"services/market-trading/internal/models"
	"services/market-trading/pkg/logger"

	"github.com/IBM/sarama"
)

// Consumer handles Kafka message consumption
type Consumer struct {
	client        sarama.ConsumerGroup
	topic         string
	priceChan     chan *models.PriceUpdate
	priceHandlers []PriceHandler
	workerCount   int
	ready         chan bool
	wg            sync.WaitGroup
}

// PriceHandler is a function that handles incoming price updates
type PriceHandler func(price *models.PriceUpdate)

// NewConsumer creates a new Kafka consumer with worker pool
func NewConsumer(brokers []string, groupID, topic string, workerCount int) (*Consumer, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V3_5_0_0
	
	// Consumer group configuration - optimized for scaling
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky // Better for scaling
	config.Consumer.Offsets.Initial = sarama.OffsetNewest // Start from latest if no offset exists
	config.Consumer.Offsets.AutoCommit.Enable = true      // Auto-commit offsets
	config.Consumer.Return.Errors = true

	client, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	if workerCount <= 0 {
		workerCount = 5 // Default worker pool size
	}

	return &Consumer{
		client:        client,
		topic:         topic,
		priceChan:     make(chan *models.PriceUpdate, 100), // Buffered channel
		priceHandlers: make([]PriceHandler, 0),
		workerCount:   workerCount,
		ready:         make(chan bool),
	}, nil
}

// RegisterPriceHandler registers a handler function for price updates
// Safe to call before Start()
func (c *Consumer) RegisterPriceHandler(handler PriceHandler) {
	c.priceHandlers = append(c.priceHandlers, handler)
}

// startWorkers starts worker goroutines to process price updates concurrently
func (c *Consumer) startWorkers(ctx context.Context) {
	logger.Infof("Starting %d worker goroutines", c.workerCount)
	
	for i := 0; i < c.workerCount; i++ {
		c.wg.Add(1)
		go func(workerID int) {
			defer c.wg.Done()
			logger.Infof("Worker %d started", workerID)
			
			for {
				select {
				case price := <-c.priceChan:
					// Call all registered handlers for this price update
					for _, handler := range c.priceHandlers {
						handler(price)
					}
				case <-ctx.Done():
					logger.Infof("Worker %d shutting down", workerID)
					return
				}
			}
		}(i)
	}
}

// Start starts consuming messages from Kafka with worker pool
func (c *Consumer) Start(ctx context.Context) error {
	logger.Infof("Starting Kafka consumer for topic: %s", c.topic)

	// Start worker pool
	c.startWorkers(ctx)

	handler := &consumerGroupHandler{
		consumer: c,
		ready:    c.ready,
	}

	// Start consuming in a goroutine
	go func() {
		for {
			// Consume should be called inside an infinite loop
			if err := c.client.Consume(ctx, []string{c.topic}, handler); err != nil {
				logger.Errorf("Error from consumer: %v", err)
			}

			// Check if context was cancelled
			if ctx.Err() != nil {
				close(c.priceChan) // Close channel to signal workers
				return
			}

			c.ready = make(chan bool)
		}
	}()

	// Wait until consumer is ready
	<-c.ready
	logger.Info("Kafka consumer is ready")

	return nil
}

// Close closes the Kafka consumer and waits for workers to finish
func (c *Consumer) Close() error {
	logger.Info("Closing Kafka consumer...")
	
	if err := c.client.Close(); err != nil {
		return fmt.Errorf("failed to close consumer: %w", err)
	}
	
	// Wait for all workers to finish processing
	c.wg.Wait()
	logger.Info("All workers stopped")
	logger.Info("Kafka consumer closed")
	return nil
}

// consumerGroupHandler implements sarama.ConsumerGroupHandler
type consumerGroupHandler struct {
	consumer *Consumer
	ready    chan bool
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	close(h.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim processes messages from a partition and sends to worker pool
func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			if message == nil {
				return nil
			}

			// Parse the price update message
			var priceUpdate models.PriceUpdate
			if err := json.Unmarshal(message.Value, &priceUpdate); err != nil {
				logger.Errorf("Failed to unmarshal price update: %v", err)
				session.MarkMessage(message, "") // Mark as processed even if error
				continue
			}

			// Send to worker pool via channel (non-blocking with buffer)
			select {
			case h.consumer.priceChan <- &priceUpdate:
				// Successfully sent to workers
			case <-session.Context().Done():
				return nil
			}

			// Mark message as processed (auto-commit will handle offset)
			session.MarkMessage(message, "")

		case <-session.Context().Done():
			return nil
		}
	}
}
