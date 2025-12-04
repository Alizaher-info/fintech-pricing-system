package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	
	"services/pricing-api/internal/config"
	"services/pricing-api/internal/database"
	"services/pricing-api/internal/repository"
	"services/pricing-api/pkg/logger"
)

func main() {
	// Initialize logger
	logger.Init()
	logger.Info("Starting Pricing API Service...")

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load configuration: %v", err))
		os.Exit(1)
	}
	logger.Info("Configuration loaded successfully")

	// Connect to database
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to database: %v", err))
		os.Exit(1)
	}
	defer db.Close()
	logger.Info(fmt.Sprintf("Connected to PostgreSQL: trading_platform"))

	// Initialize repository
	priceRepo := repository.NewPriceRepository(db)

	// Test: Get all assets
	assets, err := priceRepo.GetAllAssets()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to fetch assets: %v", err))
		os.Exit(1)
	}
	logger.Info(fmt.Sprintf("Found %d assets in database", len(assets)))
	
	for _, asset := range assets {
		logger.Info(fmt.Sprintf("  - %s (%s) [%s]", asset.Symbol, asset.Name, asset.AssetType))
	}

	// Setup graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	logger.Info("Service is running. Press Ctrl+C to stop.")
	
	// Wait for shutdown signal
	<-sigChan
	logger.Info("Shutting down gracefully...")
}
