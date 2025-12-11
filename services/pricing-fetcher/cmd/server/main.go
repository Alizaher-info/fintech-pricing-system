package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"services/pricing-fetcher/internal/config"
	"services/pricing-fetcher/internal/database"
	"services/pricing-fetcher/internal/fetcher"
	"services/pricing-fetcher/internal/kafka"
	"services/pricing-fetcher/internal/repository"
	"services/pricing-fetcher/internal/websocket"
	"services/pricing-fetcher/internal/worker"
	"services/pricing-fetcher/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init()
	logger.Info("Starting Pricing Fetcher Service...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load configuration: %v", err))
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf("Configuration loaded - Fetch interval: %d seconds", cfg.FetchInterval))

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to database: %v", err))
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("Connected to PostgreSQL: trading_platform")

	// Initialize Kafka Producer Manager
	brokers := strings.Split(cfg.KafkaBrokers, ",")
	kafkaManager := kafka.NewProducerManager(brokers)
	defer kafkaManager.CloseAll()
	
	logger.Info(fmt.Sprintf("Kafka topics configured:"))
	logger.Info(fmt.Sprintf("  - Price Updates: %s", cfg.TopicPriceUpdates))
	logger.Info(fmt.Sprintf("  - Order Created: %s", cfg.TopicOrderCreated))
	logger.Info(fmt.Sprintf("  - Price Alerts: %s", cfg.TopicPriceAlerts))

	// Initialize repositories (separated by concern)
	assetRepo := repository.NewAssetRepository(db)
	priceQuoteRepo := repository.NewPriceQuoteRepository(db)
	priceHistoryRepo := repository.NewPriceHistoryRepository(db)

	// Verify assets in database
	assets, err := assetRepo.GetAllAssets()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to fetch assets: %v", err))
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf("Found %d assets in database", len(assets)))

	// Try to start WebSocket manager (primary data source)
	wsManager := websocket.NewWebSocketManager(assetRepo, priceQuoteRepo, priceHistoryRepo, kafkaManager, cfg.TopicPriceUpdates)
	
	if err := wsManager.Start(); err != nil {
		logger.Error(fmt.Sprintf("Failed to start WebSocket manager: %v", err))
		logger.Info("Falling back to HTTP polling...")
		
		// Fallback: Use HTTP polling with CoinGecko
		priceFetcher := fetcher.NewPriceFetcher(assetRepo, priceQuoteRepo, priceHistoryRepo, kafkaManager, cfg.TopicPriceUpdates)
		priceWorker := worker.NewPriceWorker(priceFetcher, cfg.FetchInterval)
		priceWorker.Start()
		
		// Setup graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		
		logger.Info("Service is running with HTTP fallback. Press Ctrl+C to stop.")
		<-sigChan
		
		logger.Info("Shutdown signal received...")
		priceWorker.Stop()
		
	} else {
		logger.Info("WebSocket manager started successfully")
		
		// Monitor WebSocket connection health and fallback if needed
		go func() {
			ticker := time.NewTicker(30 * time.Second)
			defer ticker.Stop()
			
			for range ticker.C {
				if !wsManager.IsConnected() {
					logger.Error("WebSocket disconnected! Starting HTTP fallback...")
					
					// Start HTTP fallback
					priceFetcher := fetcher.NewPriceFetcher(assetRepo, priceQuoteRepo, priceHistoryRepo, kafkaManager, cfg.TopicPriceUpdates)
					priceWorker := worker.NewPriceWorker(priceFetcher, cfg.FetchInterval)
					priceWorker.Start()
					
					// Stop monitoring since we switched to HTTP
					return
				}
			}
		}()
		
		// Setup graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		
		logger.Info("Service is running with WebSocket. Press Ctrl+C to stop.")
		<-sigChan
		
		logger.Info("Shutdown signal received...")
		wsManager.Stop()
	}
	
	logger.Info("Service stopped gracefully.")
}
