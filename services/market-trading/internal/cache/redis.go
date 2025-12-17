package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"services/market-trading/internal/models"
	"services/market-trading/pkg/logger"

	"github.com/go-redis/redis/v8"
)

// RedisCache handles Redis caching for price data
type RedisCache struct {
	client   *redis.Client
	priceTTL time.Duration
}

// NewRedisCache creates a new Redis cache client
func NewRedisCache(host, port, password string, db int, priceTTL int) (*RedisCache, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       db,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Infof("Connected to Redis at %s:%s", host, port)

	return &RedisCache{
		client:   client,
		priceTTL: time.Duration(priceTTL) * time.Second,
	}, nil
}

// SetPrice stores a price update in Redis with TTL
func (r *RedisCache) SetPrice(ctx context.Context, price *models.PriceUpdate) error {
	key := fmt.Sprintf("price:%s", price.Symbol)

	// Serialize price to JSON
	data, err := json.Marshal(price)
	if err != nil {
		return fmt.Errorf("failed to marshal price: %w", err)
	}

	// Store in Redis with TTL
	if err := r.client.Set(ctx, key, data, r.priceTTL).Err(); err != nil {
		return fmt.Errorf("failed to set price in Redis: %w", err)
	}

	return nil
}

// GetPrice retrieves a price from Redis
func (r *RedisCache) GetPrice(ctx context.Context, symbol string) (*models.PriceUpdate, error) {
	key := fmt.Sprintf("price:%s", symbol)

	// Get from Redis
	data, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return nil, nil // Not found (not an error)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get price from Redis: %w", err)
	}

	// Deserialize JSON
	var price models.PriceUpdate
	if err := json.Unmarshal([]byte(data), &price); err != nil {
		return nil, fmt.Errorf("failed to unmarshal price: %w", err)
	}

	return &price, nil
}

// GetAllPrices retrieves all cached prices
func (r *RedisCache) GetAllPrices(ctx context.Context) ([]*models.PriceUpdate, error) {
	// Scan for all price keys
	keys, err := r.client.Keys(ctx, "price:*").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to scan keys: %w", err)
	}

	prices := make([]*models.PriceUpdate, 0, len(keys))

	// Get all prices
	for _, key := range keys {
		data, err := r.client.Get(ctx, key).Result()
		if err != nil {
			logger.Warnf("Failed to get price for key %s: %v", key, err)
			continue
		}

		var price models.PriceUpdate
		if err := json.Unmarshal([]byte(data), &price); err != nil {
			logger.Warnf("Failed to unmarshal price for key %s: %v", key, err)
			continue
		}

		prices = append(prices, &price)
	}

	return prices, nil
}

// DeletePrice removes a price from Redis
func (r *RedisCache) DeletePrice(ctx context.Context, symbol string) error {
	key := fmt.Sprintf("price:%s", symbol)
	return r.client.Del(ctx, key).Err()
}

// Close closes the Redis connection
func (r *RedisCache) Close() error {
	return r.client.Close()
}

// Health checks if Redis is healthy
func (r *RedisCache) Health(ctx context.Context) error {
	return r.client.Ping(ctx).Err()
}
