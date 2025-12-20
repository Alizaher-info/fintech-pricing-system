package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config holds all application configuration
type Config struct {
	Server ServerConfig
	Kafka  KafkaConfig
	Redis  RedisConfig
	Log    LogConfig
}

// ServerConfig holds gRPC server configuration
type ServerConfig struct {
	GRPCPort string
}

// KafkaConfig holds Kafka consumer configuration
type KafkaConfig struct {
	Brokers     []string
	Topic       string
	GroupID     string
	WorkerCount int // Number of worker goroutines for concurrent processing
}

// RedisConfig holds Redis cache configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
	PriceTTL int // TTL for price cache in seconds
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		Server: ServerConfig{
			GRPCPort: getEnv("GRPC_PORT", "50052"),
		},
		Kafka: KafkaConfig{
			Brokers:     strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
			Topic:       getEnv("KAFKA_TOPIC", "price-updates"),
			GroupID:     getEnv("KAFKA_GROUP_ID", "market-trading-group"),
			WorkerCount: getEnvAsInt("KAFKA_WORKER_COUNT", 10),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvAsInt("REDIS_DB", 0),
			PriceTTL: getEnvAsInt("REDIS_PRICE_TTL", 60),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	return config, nil
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvAsInt retrieves an environment variable as integer or returns a default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// Validate validates the configuration
func (c *Config) Validate() error {
	if len(c.Kafka.Brokers) == 0 {
		return fmt.Errorf("kafka brokers not configured")
	}

	if c.Kafka.Topic == "" {
		return fmt.Errorf("kafka topic not configured")
	}

	if c.Server.GRPCPort == "" {
		return fmt.Errorf("gRPC port not configured")
	}

	return nil
}
