package models

import "time"

// Price represents the current price quote for an asset
// Maps to the 'price_quotes' database table
// This table stores only the LATEST price for each asset (upsert behavior)
type Price struct {
	AssetID      int       `json:"asset_id"`
	CurrentPrice float64   `json:"current_price"`
	BidPrice     float64   `json:"bid_price"`
	AskPrice     float64   `json:"ask_price"`
	Volume24h    float64   `json:"volume_24h"`
	Change24h    float64   `json:"change_24h"`
	High24h      float64   `json:"high_24h"`
	Low24h       float64   `json:"low_24h"`
	MarketCap    float64   `json:"market_cap"`
	Source       string    `json:"source"`        // e.g., "coingecko", "binance"
	LastUpdated  time.Time `json:"last_updated"`
}
