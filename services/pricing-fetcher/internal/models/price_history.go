package models

import "time"

// PriceHistory represents a historical OHLC (Open, High, Low, Close) candle
// Maps to the 'price_history' database table
// This table stores historical candle data for charting and technical analysis
type PriceHistory struct {
	ID         int       `json:"id"`
	AssetID    int       `json:"asset_id"`
	OpenPrice  float64   `json:"open_price"`
	HighPrice  float64   `json:"high_price"`
	LowPrice   float64   `json:"low_price"`
	ClosePrice float64   `json:"close_price"`
	Volume     float64   `json:"volume"`
	Interval   string    `json:"interval"`   // '1m', '5m', '1h', '1d'
	Timestamp  time.Time `json:"timestamp"`
}
