package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"services/pricing-fetcher/internal/models"
)

type CoinGeckoClient struct {
	httpClient *http.Client
	baseURL    string
}

// Asset symbol to CoinGecko ID mapping
var symbolToCoinGeckoID = map[string]string{
	"BTC": "bitcoin",
	"ETH": "ethereum",
	"SOL": "solana",
}

// CoinGecko API response structure
type coinGeckoMarketResponse struct {
	ID                string  `json:"id"`
	Symbol            string  `json:"symbol"`
	Name              string  `json:"name"`
	CurrentPrice      float64 `json:"current_price"`
	MarketCap         float64 `json:"market_cap"`
	TotalVolume       float64 `json:"total_volume"`
	High24h           float64 `json:"high_24h"`
	Low24h            float64 `json:"low_24h"`
	PriceChange24h    float64 `json:"price_change_24h"`
	PriceChangePercent24h float64 `json:"price_change_percentage_24h"`
}

func NewCoinGeckoClient() *CoinGeckoClient {
	return &CoinGeckoClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: "https://api.coingecko.com/api/v3",
	}
}

// GetPrice fetches the current price for a single cryptocurrency
func (c *CoinGeckoClient) GetPrice(symbol string) (*models.CoinGeckoResponse, error) {
	coinID, exists := symbolToCoinGeckoID[symbol]
	if !exists {
		return nil, fmt.Errorf("unsupported crypto symbol: %s", symbol)
	}

	// Use markets endpoint for comprehensive data
	url := fmt.Sprintf("%s/coins/markets?vs_currency=usd&ids=%s&order=market_cap_desc&per_page=1&page=1", c.baseURL, coinID)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch price for %s: %w", coinID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var marketData []coinGeckoMarketResponse
	if err := json.NewDecoder(resp.Body).Decode(&marketData); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(marketData) == 0 {
		return nil, fmt.Errorf("no data returned for %s", coinID)
	}

	data := marketData[0]
	response := &models.CoinGeckoResponse{
		Symbol:       symbol,
		CurrentPrice: data.CurrentPrice,
		Volume:       data.TotalVolume,
		Change24h:    data.PriceChangePercent24h,
		High24h:      data.High24h,
		Low24h:       data.Low24h,
		MarketCap:    data.MarketCap,
	}

	return response, nil
}

// GetPrices fetches prices for multiple cryptocurrencies
func (c *CoinGeckoClient) GetPrices(symbols []string) ([]*models.CoinGeckoResponse, error) {
	var results []*models.CoinGeckoResponse

	for _, symbol := range symbols {
		price, err := c.GetPrice(symbol)
		if err != nil {
			// Log error but continue with other symbols
			fmt.Printf("[ERROR] Failed to fetch %s: %v\n", symbol, err)
			continue
		}

		results = append(results, price)
		
		// Rate limiting: CoinGecko free tier allows ~50 calls/minute
		// Sleep 1.5 seconds between calls (40 calls/minute to be safe)
		time.Sleep(1500 * time.Millisecond)
	}

	return results, nil
}

// IsCryptoSymbol checks if a symbol is a supported cryptocurrency
func IsCryptoSymbol(symbol string) bool {
	_, exists := symbolToCoinGeckoID[symbol]
	return exists
}
