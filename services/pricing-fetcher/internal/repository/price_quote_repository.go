package repository

import (
	"fmt"
	"time"
	
	"services/pricing-fetcher/internal/database"
	"services/pricing-fetcher/internal/models"
)

// PriceQuoteRepository handles all database operations for the price_quotes table
type PriceQuoteRepository struct {
	db *database.DB
}

// NewPriceQuoteRepository creates a new PriceQuoteRepository instance
func NewPriceQuoteRepository(db *database.DB) *PriceQuoteRepository {
	return &PriceQuoteRepository{db: db}
}

// UpsertPrice inserts or updates a price quote for an asset
// This maintains only the current/latest price for each asset
func (r *PriceQuoteRepository) UpsertPrice(price *models.Price) error {
	query := `
		INSERT INTO price_quotes (
			asset_id, current_price, bid_price, ask_price,
			volume_24h, change_24h, high_24h, low_24h,
			market_cap, source, last_updated
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (asset_id) DO UPDATE SET
			current_price = EXCLUDED.current_price,
			bid_price = EXCLUDED.bid_price,
			ask_price = EXCLUDED.ask_price,
			volume_24h = EXCLUDED.volume_24h,
			change_24h = EXCLUDED.change_24h,
			high_24h = CASE
				WHEN EXCLUDED.current_price > price_quotes.high_24h
				THEN EXCLUDED.current_price
				ELSE price_quotes.high_24h
			END,
			low_24h = CASE
				WHEN EXCLUDED.current_price < price_quotes.low_24h
				THEN EXCLUDED.current_price
				ELSE price_quotes.low_24h
			END,
			market_cap = EXCLUDED.market_cap,
			source = EXCLUDED.source,
			last_updated = EXCLUDED.last_updated
	`

	_, err := r.db.Exec(query,
		price.AssetID,
		price.CurrentPrice,
		price.BidPrice,
		price.AskPrice,
		price.Volume24h,
		price.Change24h,
		price.High24h,
		price.Low24h,
		price.MarketCap,
		price.Source,
		time.Now(),
	)

	if err != nil {
		return fmt.Errorf("failed to upsert price quote: %w", err)
	}

	return nil
}

// GetCurrentPrice retrieves the current price quote for a specific asset
func (r *PriceQuoteRepository) GetCurrentPrice(assetID int) (*models.Price, error) {
	query := `
		SELECT asset_id, current_price, bid_price, ask_price,
		       volume_24h, change_24h, high_24h, low_24h,
		       market_cap, source, last_updated
		FROM price_quotes
		WHERE asset_id = $1
	`
	
	var price models.Price
	err := r.db.QueryRow(query, assetID).Scan(
		&price.AssetID,
		&price.CurrentPrice,
		&price.BidPrice,
		&price.AskPrice,
		&price.Volume24h,
		&price.Change24h,
		&price.High24h,
		&price.Low24h,
		&price.MarketCap,
		&price.Source,
		&price.LastUpdated,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get current price: %w", err)
	}

	return &price, nil
}

// GetAllCurrentPrices retrieves all current price quotes
func (r *PriceQuoteRepository) GetAllCurrentPrices() ([]models.Price, error) {
	query := `
		SELECT asset_id, current_price, bid_price, ask_price,
		       volume_24h, change_24h, high_24h, low_24h,
		       market_cap, source, last_updated
		FROM price_quotes
		ORDER BY asset_id
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query all current prices: %w", err)
	}
	defer rows.Close()

	var prices []models.Price
	for rows.Next() {
		var price models.Price
		if err := rows.Scan(
			&price.AssetID,
			&price.CurrentPrice,
			&price.BidPrice,
			&price.AskPrice,
			&price.Volume24h,
			&price.Change24h,
			&price.High24h,
			&price.Low24h,
			&price.MarketCap,
			&price.Source,
			&price.LastUpdated,
		); err != nil {
			return nil, fmt.Errorf("failed to scan price: %w", err)
		}
		prices = append(prices, price)
	}

	return prices, nil
}
