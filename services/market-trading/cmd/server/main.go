package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"services/market-trading/internal/cache"
	"services/market-trading/internal/config"
	grpcServer "services/market-trading/internal/grpc"
	"services/market-trading/internal/kafka"
	"services/market-trading/internal/models"
	"services/market-trading/pkg/logger"

	pb "services/market-trading/gen/market/v1"
	"google.golang.org/grpc"
)

func main() {
	logger.Info("Starting Market Trading Service...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}

	if err := cfg.Validate(); err != nil {
		logger.Fatalf("Invalid configuration: %v", err)
	}

	logger.Infof("Configuration loaded successfully")

	// Initialize Redis cache
	redisCache, err := cache.NewRedisCache(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		cfg.Redis.DB,
		cfg.Redis.PriceTTL,
	)
	if err != nil {
		logger.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisCache.Close()

	// Initialize gRPC server
	grpcSrv := grpcServer.NewServer(redisCache)

	// Initialize Kafka consumer with worker pool
	consumer, err := kafka.NewConsumer(
		cfg.Kafka.Brokers,
		cfg.Kafka.GroupID,
		cfg.Kafka.Topic,
		cfg.Kafka.WorkerCount,
	)
	if err != nil {
		logger.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Register price handler: Cache + Broadcast
	consumer.RegisterPriceHandler(func(price *models.PriceUpdate) {
		ctx := context.Background()

		// Save to Redis cache
		if err := redisCache.SetPrice(ctx, price); err != nil {
			logger.Errorf("Failed to cache price for %s: %v", price.Symbol, err)
		}

		// Broadcast to all gRPC streaming clients
		grpcSrv.BroadcastPrice(price)

		// Log every 10th update to avoid spam
		// You can adjust this or remove it based on your needs
		logger.Debugf("Price updated: %s = $%.2f", price.Symbol, price.Price)
	})

	// Start Kafka consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := consumer.Start(ctx); err != nil {
		logger.Fatalf("Failed to start Kafka consumer: %v", err)
	}

	// Start gRPC server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.Server.GRPCPort))
	if err != nil {
		logger.Fatalf("Failed to listen on port %s: %v", cfg.Server.GRPCPort, err)
	}

	grpcServerInstance := grpc.NewServer()
	pb.RegisterMarketDataServiceServer(grpcServerInstance, grpcSrv)

	// Start gRPC server in goroutine
	go func() {
		logger.Infof("gRPC server listening on port %s", cfg.Server.GRPCPort)
		if err := grpcServerInstance.Serve(lis); err != nil {
			logger.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Status log
	logger.Info("✅ Market Trading Service is running")
	logger.Infof("📊 Kafka Topic: %s", cfg.Kafka.Topic)
	logger.Infof("🔌 gRPC Port: %s", cfg.Server.GRPCPort)
	logger.Infof("💾 Redis: %s:%s", cfg.Redis.Host, cfg.Redis.Port)

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	logger.Info("Shutting down Market Trading Service...")

	// Graceful shutdown
	cancel()
	grpcServerInstance.GracefulStop()

	logger.Info("Market Trading Service stopped")
}
