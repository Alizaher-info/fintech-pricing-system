package models

import "time"

// PriceUpdate represents a price update message from Kafka
type PriceUpdate struct {
	Symbol    string    `json:"symbol"`
	Price     float64   `json:"price"`
	BidPrice  float64   `json:"bid_price"`
	AskPrice  float64   `json:"ask_price"`
	Volume24h float64   `json:"volume_24h"`
	Change24h float64   `json:"change_24h"`
	High24h   float64   `json:"high_24h"`
	Low24h    float64   `json:"low_24h"`
	Timestamp time.Time `json:"timestamp"`
}
