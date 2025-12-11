package websocket

import (
	"sync"
	"time"
)

// OHLCCandle represents aggregated OHLC data for a time interval
type OHLCCandle struct {
	Symbol    string
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Volume    float64
	Trades    int
	StartTime time.Time
}

// OHLCAggregator aggregates price updates into OHLC candles
type OHLCAggregator struct {
	// Current candles being built (one per symbol)
	currentCandles map[string]*OHLCCandle
	
	// Last saved price for each symbol (for threshold check)
	lastSavedPrices map[string]float64
	
	// Thread safety
	mu sync.RWMutex
	
	// Configuration
	aggregationInterval  time.Duration // Default: 4 seconds
	priceChangeThreshold float64       // Default: 0.1% (0.001)
}

// NewOHLCAggregator creates a new OHLC aggregator
func NewOHLCAggregator(interval time.Duration, threshold float64) *OHLCAggregator {
	return &OHLCAggregator{
		currentCandles:       make(map[string]*OHLCCandle),
		lastSavedPrices:      make(map[string]float64),
		aggregationInterval:  interval,
		priceChangeThreshold: threshold,
	}
}

// AddTrade adds a new trade to the current candle
func (oa *OHLCAggregator) AddTrade(symbol string, price float64, volume float64) {
	oa.mu.Lock()
	defer oa.mu.Unlock()
	
	candle, exists := oa.currentCandles[symbol]
	
	if !exists {
		// Create new candle
		oa.currentCandles[symbol] = &OHLCCandle{
			Symbol:    symbol,
			Open:      price,
			High:      price,
			Low:       price,
			Close:     price,
			Volume:    volume,
			Trades:    1,
			StartTime: time.Now(),
		}
		return
	}
	
	// Update existing candle
	if price > candle.High {
		candle.High = price
	}
	if price < candle.Low {
		candle.Low = price
	}
	candle.Close = price
	candle.Volume += volume
	candle.Trades++
}

// GetAndResetCandle retrieves the current candle and resets it
func (oa *OHLCAggregator) GetAndResetCandle(symbol string) *OHLCCandle {
	oa.mu.Lock()
	defer oa.mu.Unlock()
	
	candle, exists := oa.currentCandles[symbol]
	if !exists {
		return nil
	}
	
	// Create a copy to return
	result := &OHLCCandle{
		Symbol:    candle.Symbol,
		Open:      candle.Open,
		High:      candle.High,
		Low:       candle.Low,
		Close:     candle.Close,
		Volume:    candle.Volume,
		Trades:    candle.Trades,
		StartTime: candle.StartTime,
	}
	
	// Reset the candle
	delete(oa.currentCandles, symbol)
	
	return result
}

// ShouldUpdatePrice checks if price changed more than threshold
func (oa *OHLCAggregator) ShouldUpdatePrice(symbol string, newPrice float64) bool {
	oa.mu.RLock()
	lastPrice, exists := oa.lastSavedPrices[symbol]
	oa.mu.RUnlock()
	
	if !exists {
		// First time, always update
		return true
	}
	
	// Calculate percentage change
	change := (newPrice - lastPrice) / lastPrice
	if change < 0 {
		change = -change // Absolute value
	}
	
	return change >= oa.priceChangeThreshold
}

// UpdateLastSavedPrice updates the last saved price for a symbol
func (oa *OHLCAggregator) UpdateLastSavedPrice(symbol string, price float64) {
	oa.mu.Lock()
	defer oa.mu.Unlock()
	oa.lastSavedPrices[symbol] = price
}

// GetAggregationInterval returns the aggregation interval
func (oa *OHLCAggregator) GetAggregationInterval() time.Duration {
	return oa.aggregationInterval
}
