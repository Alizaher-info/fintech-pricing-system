package models

// CoinGeckoResponse represents the response structure from CoinGecko API
// This is NOT a database table - it's used for API integration only
type CoinGeckoResponse struct {
	Symbol       string  `json:"symbol"`
	CurrentPrice float64 `json:"current_price"`
	Volume       float64 `json:"total_volume"`
	Change24h    float64 `json:"price_change_percentage_24h"`
	High24h      float64 `json:"high_24h"`
	Low24h       float64 `json:"low_24h"`
	MarketCap    float64 `json:"market_cap"`
}
