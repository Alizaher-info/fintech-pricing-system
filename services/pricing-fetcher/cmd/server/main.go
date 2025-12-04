package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"services/pricing-fetcher/internal/config"
	"services/pricing-fetcher/internal/database"
	"services/pricing-fetcher/internal/fetcher"
	"services/pricing-fetcher/internal/repository"
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

	// Initialize repository and fetcher
	priceRepo := repository.NewPriceRepository(db)
	priceFetcher := fetcher.NewPriceFetcher(priceRepo)

	// Verify assets in database
	assets, err := priceRepo.GetAllAssets()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to fetch assets: %v", err))
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf("Found %d assets in database", len(assets)))

	// Create and start background worker
	priceWorker := worker.NewPriceWorker(priceFetcher, cfg.FetchInterval)
	priceWorker.Start()

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	logger.Info("Service is running. Press Ctrl+C to stop.")
	
	// Wait for shutdown signal
	<-sigChan
	logger.Info("Shutdown signal received...")
	
	// Stop background worker gracefully
	priceWorker.Stop()
	
	logger.Info("Service stopped gracefully.")
}
