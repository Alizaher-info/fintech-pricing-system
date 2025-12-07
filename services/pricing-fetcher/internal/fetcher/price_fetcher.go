package fetcher

import (
	"context"
	"fmt"
	"time"

	"services/pricing-fetcher/internal/kafka"
	"services/pricing-fetcher/internal/models"
	"services/pricing-fetcher/internal/repository"
	"services/pricing-fetcher/pkg/logger"
)

type PriceFetcher struct {
	coinGeckoClient *CoinGeckoClient
	priceRepo       *repository.PriceRepository
	kafkaManager    *kafka.ProducerManager
	topicName       string
}

func NewPriceFetcher(priceRepo *repository.PriceRepository, kafkaManager *kafka.ProducerManager, topicName string) *PriceFetcher {
	return &PriceFetcher{
		coinGeckoClient: NewCoinGeckoClient(),
		priceRepo:       priceRepo,
		kafkaManager:    kafkaManager,
		topicName:       topicName,
	}
}

// FetchAndSaveAll fetches prices for all crypto assets and saves them to the database
func (pf *PriceFetcher) FetchAndSaveAll() error {
	logger.Info("Starting price fetch cycle...")

	// Get all assets from database
	assets, err := pf.priceRepo.GetAllAssets()
	if err != nil {
		return fmt.Errorf("failed to get assets: %w", err)
	}

	// Filter only crypto assets
	var cryptoSymbols []string
	assetMap := make(map[string]*models.Asset)
	
	for _, asset := range assets {
		if asset.AssetType == "crypto" && IsCryptoSymbol(asset.Symbol) {
			cryptoSymbols = append(cryptoSymbols, asset.Symbol)
			assetMap[asset.Symbol] = &asset
		}
	}

	logger.Info(fmt.Sprintf("Fetching prices for %d crypto assets...", len(cryptoSymbols)))

	// Fetch prices from CoinGecko
	prices, err := pf.coinGeckoClient.GetPrices(cryptoSymbols)
	if err != nil {
		return fmt.Errorf("failed to fetch prices: %w", err)
	}

	// Save each price to database
	successCount := 0
	for _, priceData := range prices {
		asset, exists := assetMap[priceData.Symbol]
		if !exists {
			logger.Error(fmt.Sprintf("Asset not found for symbol: %s", priceData.Symbol))
			continue
		}

		price := &models.Price{
			AssetID:      asset.ID,
			CurrentPrice: priceData.CurrentPrice,
			BidPrice:     priceData.CurrentPrice * 0.9995, // Simulate bid (0.05% below)
			AskPrice:     priceData.CurrentPrice * 1.0005, // Simulate ask (0.05% above)
			Volume24h:    priceData.Volume,
			Change24h:    priceData.Change24h,
			High24h:      priceData.High24h,
			Low24h:       priceData.Low24h,
			MarketCap:    priceData.MarketCap,
			Source:       "coingecko",
			LastUpdated:  time.Now(),
		}

		if err := pf.priceRepo.UpsertPrice(price); err != nil {
			logger.Error(fmt.Sprintf("Failed to save price for %s: %v", priceData.Symbol, err))
			continue
		}

		// Publish to Kafka after successful database save
		if pf.kafkaManager != nil {
			kafkaMsg := kafka.PriceUpdateMessage{
				Symbol:    priceData.Symbol,
				Price:     priceData.CurrentPrice,
				Change24h: priceData.Change24h,
				Volume24h: priceData.Volume,
				High24h:   priceData.High24h,
				Low24h:    priceData.Low24h,
				MarketCap: priceData.MarketCap,
				Timestamp: time.Now(),
			}
			
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			if err := pf.kafkaManager.PublishPriceUpdate(ctx, pf.topicName, kafkaMsg); err != nil {
				logger.Error(fmt.Sprintf("Failed to publish to Kafka for %s: %v", priceData.Symbol, err))
			}
			cancel()
		}

		logger.Info(fmt.Sprintf("✓ %s: $%.2f | 24h: %.2f%% | Vol: $%.0f",
			priceData.Symbol,
			priceData.CurrentPrice,
			priceData.Change24h,
			priceData.Volume,
		))

		successCount++
	}

	logger.Info(fmt.Sprintf("Price fetch completed: %d/%d successful", successCount, len(cryptoSymbols)))
	return nil
}

// FetchSingle fetches and saves price for a single asset
func (pf *PriceFetcher) FetchSingle(symbol string) error {
	asset, err := pf.priceRepo.GetAssetBySymbol(symbol)
	if err != nil {
		return fmt.Errorf("failed to get asset: %w", err)
	}

	if !IsCryptoSymbol(symbol) {
		return fmt.Errorf("symbol %s is not a supported cryptocurrency", symbol)
	}

	priceData, err := pf.coinGeckoClient.GetPrice(symbol)
	if err != nil {
		return fmt.Errorf("failed to fetch price: %w", err)
	}

	price := &models.Price{
		AssetID:      asset.ID,
		CurrentPrice: priceData.CurrentPrice,
		BidPrice:     priceData.CurrentPrice * 0.9995,
		AskPrice:     priceData.CurrentPrice * 1.0005,
		Volume24h:    priceData.Volume,
		Change24h:    priceData.Change24h,
		High24h:      priceData.High24h,
		Low24h:       priceData.Low24h,
		MarketCap:    priceData.MarketCap,
		Source:       "coingecko",
		LastUpdated:  time.Now(),
	}

	if err := pf.priceRepo.UpsertPrice(price); err != nil {
		return fmt.Errorf("failed to save price: %w", err)
	}

	logger.Info(fmt.Sprintf("Updated %s: $%.2f", symbol, priceData.CurrentPrice))
	return nil
}
