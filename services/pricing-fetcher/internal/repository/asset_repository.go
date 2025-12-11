package repository

import (
	"database/sql"
	"fmt"
	
	"services/pricing-fetcher/internal/database"
	"services/pricing-fetcher/internal/models"
)

// AssetRepository handles all database operations for the assets table
type AssetRepository struct {
	db *database.DB
}

// NewAssetRepository creates a new AssetRepository instance
func NewAssetRepository(db *database.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

// GetAllAssets retrieves all active assets from the database
func (r *AssetRepository) GetAllAssets() ([]models.Asset, error) {
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

// GetAssetBySymbol retrieves a single asset by its symbol
func (r *AssetRepository) GetAssetBySymbol(symbol string) (*models.Asset, error) {
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

// GetAssetByID retrieves a single asset by its ID
func (r *AssetRepository) GetAssetByID(id int) (*models.Asset, error) {
	query := `SELECT id, symbol, name, asset_type FROM assets WHERE id = $1 AND is_active = true`
	
	var asset models.Asset
	err := r.db.QueryRow(query, id).Scan(&asset.ID, &asset.Symbol, &asset.Name, &asset.AssetType)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("asset not found with id: %d", id)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query asset: %w", err)
	}

	return &asset, nil
}
