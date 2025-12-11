package websocket

import (
	"context"
	"fmt"
	"sync"
	"time"

	"services/pricing-fetcher/internal/kafka"
	"services/pricing-fetcher/internal/models"
	"services/pricing-fetcher/internal/repository"
	"services/pricing-fetcher/pkg/logger"
)

// WebSocketManager manages WebSocket connections and price updates
type WebSocketManager struct {
	binanceClient    *BinanceClient
	assetRepo        *repository.AssetRepository
	priceQuoteRepo   *repository.PriceQuoteRepository
	priceHistoryRepo *repository.PriceHistoryRepository
	kafkaManager     *kafka.ProducerManager
	topicName        string
	
	// Track symbols we're monitoring
	symbols          []string
	assetMap         map[string]*models.Asset
	mu               sync.RWMutex
	
	// OHLC aggregation
	aggregator       *OHLCAggregator
	
	// Stop channel
	stopChan         chan struct{}
	wg               sync.WaitGroup
}

// NewWebSocketManager creates a new WebSocket manager
func NewWebSocketManager(
	assetRepo *repository.AssetRepository,
	priceQuoteRepo *repository.PriceQuoteRepository,
	priceHistoryRepo *repository.PriceHistoryRepository,
	kafkaManager *kafka.ProducerManager,
	topicName string,
) *WebSocketManager {
	return &WebSocketManager{
		assetRepo:        assetRepo,
		priceQuoteRepo:   priceQuoteRepo,
		priceHistoryRepo: priceHistoryRepo,
		kafkaManager:     kafkaManager,
		topicName:        topicName,
		assetMap:         make(map[string]*models.Asset),
		aggregator:       NewOHLCAggregator(4*time.Second, 0.001), // 4 seconds, 0.1% threshold
		stopChan:         make(chan struct{}),
	}
}

// Start initializes WebSocket connection and starts listening for price updates
func (wm *WebSocketManager) Start() error {
	logger.Info("Starting WebSocket Manager...")
	
	// Get all crypto assets from database
	assets, err := wm.assetRepo.GetAllAssets()
	if err != nil {
		return fmt.Errorf("failed to get assets: %w", err)
	}
	
	// Filter crypto symbols and build asset map
	wm.mu.Lock()
	for _, asset := range assets {
		if asset.AssetType == "crypto" {
			// Convert to lowercase for Binance format (btcusdt, ethusdt, solusdt)
			wm.symbols = append(wm.symbols, asset.Symbol)
			wm.assetMap[asset.Symbol] = &asset
		}
	}
	wm.mu.Unlock()
	
	if len(wm.symbols) == 0 {
		return fmt.Errorf("no crypto assets found in database")
	}
	
	logger.Info(fmt.Sprintf("Monitoring %d crypto assets via WebSocket", len(wm.symbols)))
	
	// Create Binance WebSocket client
	wm.binanceClient = NewBinanceClient(wm.symbols)
	
	// Connect to Binance WebSocket
	if err := wm.binanceClient.Connect(); err != nil {
		return fmt.Errorf("failed to connect to Binance WebSocket: %w", err)
	}
	
	logger.Info("Connected to Binance WebSocket")
	
	// Start processing price updates
	wm.wg.Add(1)
	go wm.processPriceUpdates()
	
	return nil
}

// processPriceUpdates listens to price updates from WebSocket and aggregates them
func (wm *WebSocketManager) processPriceUpdates() {
	defer wm.wg.Done()
	
	logger.Info("WebSocket price processor started")
	
	// Create ticker for OHLC aggregation (every 4 seconds)
	ticker := time.NewTicker(wm.aggregator.GetAggregationInterval())
	defer ticker.Stop()
	
	for {
		select {
		case <-wm.stopChan:
			logger.Info("Stopping WebSocket price processor...")
			return
			
		case priceUpdate := <-wm.binanceClient.PriceUpdates():
			// Add trade to aggregator (does NOT save to database yet)
			wm.aggregator.AddTrade(priceUpdate.Symbol, priceUpdate.Price, priceUpdate.Volume)
			
		case <-ticker.C:
			// Time to save aggregated candles (every 4 seconds)
			wm.saveAggregatedCandles()
			
		case err := <-wm.binanceClient.Errors():
			logger.Error(fmt.Sprintf("WebSocket error: %v", err))
			// TODO: Implement reconnection logic here
		}
	}
}

// saveAggregatedCandles saves aggregated OHLC candles for all symbols
func (wm *WebSocketManager) saveAggregatedCandles() {
	wm.mu.RLock()
	symbols := wm.symbols
	wm.mu.RUnlock()
	
	for _, symbol := range symbols {
		// Get the completed candle
		candle := wm.aggregator.GetAndResetCandle(symbol)
		if candle == nil {
			// No trades received for this symbol in the last 4 seconds
			continue
		}
		
		// Get asset info
		wm.mu.RLock()
		asset, exists := wm.assetMap[symbol]
		wm.mu.RUnlock()
		
		if !exists {
			logger.Error(fmt.Sprintf("Asset not found for symbol: %s", symbol))
			continue
		}
		
		// Save OHLC to price_history table
		wm.savePriceHistory(asset, candle)
		
		// Check if price changed enough to update price_quotes
		if wm.aggregator.ShouldUpdatePrice(symbol, candle.Close) {
			wm.savePriceQuote(asset, candle)
			wm.aggregator.UpdateLastSavedPrice(symbol, candle.Close)
		}
		
		// Publish to Kafka
		wm.publishToKafka(symbol, candle)
		
		logger.Info(fmt.Sprintf("✓ [WS] %s: $%.2f (O:%.2f H:%.2f L:%.2f) | %d trades", 
			symbol, candle.Close, candle.Open, candle.High, candle.Low, candle.Trades))
	}
}

// savePriceHistory saves OHLC candle to price_history table
func (wm *WebSocketManager) savePriceHistory(asset *models.Asset, candle *OHLCCandle) {
	price := &models.Price{
		AssetID:      asset.ID,
		CurrentPrice: candle.Close,
		BidPrice:     candle.Close * 0.9995,
		AskPrice:     candle.Close * 1.0005,
		Volume24h:    candle.Volume,
		Change24h:    0,
		High24h:      candle.High,
		Low24h:       candle.Low,
		MarketCap:    0,
		Source:       "binance_ws",
		LastUpdated:  time.Now(),
	}
	
	if err := wm.priceHistoryRepo.InsertHistory(price); err != nil {
		logger.Error(fmt.Sprintf("Failed to save price history for %s: %v", candle.Symbol, err))
	}
}

// savePriceQuote updates current price in price_quotes table
func (wm *WebSocketManager) savePriceQuote(asset *models.Asset, candle *OHLCCandle) {
	price := &models.Price{
		AssetID:      asset.ID,
		CurrentPrice: candle.Close,
		BidPrice:     candle.Close * 0.9995,
		AskPrice:     candle.Close * 1.0005,
		Volume24h:    candle.Volume,
		Change24h:    0,
		High24h:      candle.High,
		Low24h:       candle.Low,
		MarketCap:    0,
		Source:       "binance_ws",
		LastUpdated:  time.Now(),
	}
	
	if err := wm.priceQuoteRepo.UpsertPrice(price); err != nil {
		logger.Error(fmt.Sprintf("Failed to update price quote for %s: %v", candle.Symbol, err))
	}
}

// publishToKafka publishes price update to Kafka
func (wm *WebSocketManager) publishToKafka(symbol string, candle *OHLCCandle) {
	if wm.kafkaManager == nil {
		return
	}
	
	kafkaMsg := kafka.PriceUpdateMessage{
		Symbol:    symbol,
		Price:     candle.Close,
		Change24h: 0,
		Volume24h: candle.Volume,
		High24h:   candle.High,
		Low24h:    candle.Low,
		MarketCap: 0,
		Timestamp: time.Now(),
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := wm.kafkaManager.PublishPriceUpdate(ctx, wm.topicName, kafkaMsg); err != nil {
		logger.Error(fmt.Sprintf("Failed to publish to Kafka for %s: %v", symbol, err))
	}
}

// Stop gracefully stops the WebSocket manager
func (wm *WebSocketManager) Stop() {
	logger.Info("Stopping WebSocket Manager...")
	
	// Signal stop to processing goroutine
	close(wm.stopChan)
	
	// Close Binance client
	if wm.binanceClient != nil {
		wm.binanceClient.Close()
	}
	
	// Wait for all goroutines to finish
	wm.wg.Wait()
	
	logger.Info("WebSocket Manager stopped")
}

// IsConnected returns the connection status
func (wm *WebSocketManager) IsConnected() bool {
	if wm.binanceClient == nil {
		return false
	}
	return wm.binanceClient.IsConnected()
}
