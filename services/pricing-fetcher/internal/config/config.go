package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL     string
	CoinGeckoAPIKey string
	FetchInterval   int
	ServicePort     string
	LogLevel        string
	KafkaBrokers    string
	// Topic configuration
	TopicPriceUpdates string
	TopicOrderCreated string
	TopicPriceAlerts  string
}

func Load() (*Config, error) {
	// Load .env file (ignore error if not found, use env vars instead)
	_ = godotenv.Load()

	fetchInterval, _ := strconv.Atoi(getEnv("FETCH_INTERVAL", "60"))

	return &Config{
		DatabaseURL:       getEnv("DATABASE_URL", "postgres://app:app@localhost:5432/trading_platform?sslmode=disable"),
		CoinGeckoAPIKey:   getEnv("COINGECKO_API_KEY", ""),
		FetchInterval:     fetchInterval,
		ServicePort:       getEnv("SERVICE_PORT", "50051"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		KafkaBrokers:      getEnv("KAFKA_BROKERS", "kafka:9092"),
		TopicPriceUpdates: getEnv("TOPIC_PRICE_UPDATES", "price-updates"),
		TopicOrderCreated: getEnv("TOPIC_ORDER_CREATED", "order-created"),
		TopicPriceAlerts:  getEnv("TOPIC_PRICE_ALERTS", "price-alerts"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
