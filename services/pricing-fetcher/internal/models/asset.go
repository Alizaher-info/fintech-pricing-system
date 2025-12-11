package models

// Asset represents a cryptocurrency or financial asset in the system
// Maps to the 'assets' database table
type Asset struct {
	ID        int    `json:"id"`
	Symbol    string `json:"symbol"`     // e.g., "BTC", "ETH", "SOL"
	Name      string `json:"name"`       // e.g., "Bitcoin", "Ethereum"
	AssetType string `json:"asset_type"` // e.g., "crypto", "stock", "forex"
}
