package repository

import (
	"database/sql"
	"fmt"
	"time"
	
	"services/pricing-fetcher/internal/database"
	"services/pricing-fetcher/internal/models"
)

type PriceRepository struct {
	db *database.DB
}

func NewPriceRepository(db *database.DB) *PriceRepository {
	return &PriceRepository{db: db}
}

// GetAllAssets retrieves all assets from the database
func (r *PriceRepository) GetAllAssets() ([]models.Asset, error) {
	query := `SELECT id, symbol, name, asset_type FROM assets WHERE is_active = true`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to query assets: %w", err)
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var asset models.Asset
		if err := rows.Scan(&asset.ID, &asset.Symbol, &asset.Name, &asset.AssetType); err != nil {
			return nil, fmt.Errorf("failed to scan asset: %w", err)
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

// UpsertPrice inserts or updates a price quote
func (r *PriceRepository) UpsertPrice(price *models.Price) error {
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
		return fmt.Errorf("failed to upsert price: %w", err)
	}

	return nil
}

// GetAssetBySymbol retrieves an asset by its symbol
func (r *PriceRepository) GetAssetBySymbol(symbol string) (*models.Asset, error) {
	query := `SELECT id, symbol, name, asset_type FROM assets WHERE symbol = $1 AND is_active = true`
	
	var asset models.Asset
	err := r.db.QueryRow(query, symbol).Scan(&asset.ID, &asset.Symbol, &asset.Name, &asset.AssetType)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("asset not found: %s", symbol)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query asset: %w", err)
	}

	return &asset, nil
}
