package repository

import (
	"fmt"
	"time"
	
	"services/pricing-fetcher/internal/database"
	"services/pricing-fetcher/internal/models"
)

// PriceHistoryRepository handles all database operations for the price_history table
type PriceHistoryRepository struct {
	db *database.DB
}

// NewPriceHistoryRepository creates a new PriceHistoryRepository instance
func NewPriceHistoryRepository(db *database.DB) *PriceHistoryRepository {
	return &PriceHistoryRepository{db: db}
}

// InsertHistory inserts a new OHLC candle record into price history
// For 1-minute candles based on real-time price fetches
// Since we fetch every 60 seconds, each fetch creates one 1m candle
func (r *PriceHistoryRepository) InsertHistory(price *models.Price) error {
	query := `
		INSERT INTO price_history (
			asset_id, open_price, high_price, low_price, close_price, 
			volume, interval, timestamp
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	// For a single price fetch, we create a 1-minute candle
	// Open = Close = current price (single data point)
	// High = high_24h from API
	// Low = low_24h from API
	_, err := r.db.Exec(query,
		price.AssetID,
		price.CurrentPrice, // open_price
		price.High24h,      // high_price
		price.Low24h,       // low_price
		price.CurrentPrice, // close_price (same as open for single fetch)
		price.Volume24h,    // volume
		"1m",              // interval (1 minute)
		time.Now(),        // timestamp
	)

	if err != nil {
		return fmt.Errorf("failed to insert price history: %w", err)
	}

	return nil
}

// GetHistoryByAsset retrieves OHLC candle history for a specific asset within a time range
func (r *PriceHistoryRepository) GetHistoryByAsset(assetID int, startTime, endTime time.Time) ([]models.PriceHistory, error) {
	query := `
		SELECT id, asset_id, open_price, high_price, low_price, close_price, 
		       volume, interval, timestamp
		FROM price_history
		WHERE asset_id = $1 AND timestamp BETWEEN $2 AND $3
		ORDER BY timestamp ASC
	`
	
	rows, err := r.db.Query(query, assetID, startTime, endTime)
	if err != nil {
		return nil, fmt.Errorf("failed to query price history: %w", err)
	}
	defer rows.Close()

	var history []models.PriceHistory
	for rows.Next() {
		var record models.PriceHistory
		if err := rows.Scan(
			&record.ID,
			&record.AssetID,
			&record.OpenPrice,
			&record.HighPrice,
			&record.LowPrice,
			&record.ClosePrice,
			&record.Volume,
			&record.Interval,
			&record.Timestamp,
		); err != nil {
			return nil, fmt.Errorf("failed to scan price history: %w", err)
		}
		history = append(history, record)
	}

	return history, nil
}

// GetLatestHistory retrieves the most recent N candle records for an asset
func (r *PriceHistoryRepository) GetLatestHistory(assetID int, limit int) ([]models.PriceHistory, error) {
	query := `
		SELECT id, asset_id, open_price, high_price, low_price, close_price,
		       volume, interval, timestamp
		FROM price_history
		WHERE asset_id = $1
		ORDER BY timestamp DESC
		LIMIT $2
	`
	
	rows, err := r.db.Query(query, assetID, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query latest price history: %w", err)
	}
	defer rows.Close()

	var history []models.PriceHistory
	for rows.Next() {
		var record models.PriceHistory
		if err := rows.Scan(
			&record.ID,
			&record.AssetID,
			&record.OpenPrice,
			&record.HighPrice,
			&record.LowPrice,
			&record.ClosePrice,
			&record.Volume,
			&record.Interval,
			&record.Timestamp,
		); err != nil {
			return nil, fmt.Errorf("failed to scan price history: %w", err)
		}
		history = append(history, record)
	}

	return history, nil
}

// DeleteOldHistory removes price history records older than the specified retention period
// Useful for cleanup jobs to prevent database bloat
func (r *PriceHistoryRepository) DeleteOldHistory(retentionDays int) (int64, error) {
	query := `
		DELETE FROM price_history
		WHERE timestamp < NOW() - INTERVAL '%d days'
	`
	
	result, err := r.db.Exec(fmt.Sprintf(query, retentionDays))
	if err != nil {
		return 0, fmt.Errorf("failed to delete old price history: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}
