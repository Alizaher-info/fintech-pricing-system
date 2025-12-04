package models

import "time"

type Asset struct {
	ID        int
	Symbol    string
	Name      string
	AssetType string
}

type Price struct {
	AssetID      int
	CurrentPrice float64
	BidPrice     float64
	AskPrice     float64
	Volume24h    float64
	Change24h    float64
	High24h      float64
	Low24h       float64
	MarketCap    float64
	Source       string
	LastUpdated  time.Time
}

type CoinGeckoResponse struct {
	Symbol       string
	CurrentPrice float64
	Volume       float64
	Change24h    float64
	High24h      float64
	Low24h       float64
	MarketCap    float64
}
