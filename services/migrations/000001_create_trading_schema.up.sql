-- Trading Platform Database Schema
-- Created: November 30, 2025
-- Database: PostgreSQL 16

-- 1. ASSETS Table - Master table for all tradable assets
CREATE TABLE IF NOT EXISTS assets (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    asset_type VARCHAR(20) NOT NULL,
    exchange VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT check_asset_type CHECK (asset_type IN ('crypto', 'stock', 'forex'))
);

CREATE INDEX idx_assets_type ON assets(asset_type);
CREATE INDEX idx_assets_symbol ON assets(symbol);

-- 2. PRICE_QUOTES Table - Real-time current prices
CREATE TABLE IF NOT EXISTS price_quotes (
    id SERIAL PRIMARY KEY,
    asset_id INTEGER NOT NULL,
    current_price DECIMAL(20, 8) NOT NULL,
    bid_price DECIMAL(20, 8), -- Highest bid price
    ask_price DECIMAL(20, 8), -- Lowest ask price
    volume_24h DECIMAL(30, 2), -- 24-hour trading volume to check if trading is active good indicator
    change_24h DECIMAL(10, 4),
    high_24h DECIMAL(20, 8),
    low_24h DECIMAL(20, 8),
    market_cap DECIMAL(30, 2), -- Market capitalization for stocks/crypto to assess size
    source VARCHAR(50),
    last_updated TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT fk_price_quotes_asset FOREIGN KEY (asset_id) 
        REFERENCES assets(id) ON DELETE CASCADE,
    CONSTRAINT unique_asset_quote UNIQUE(asset_id)
);

CREATE INDEX idx_quotes_updated ON price_quotes(last_updated);

-- 3. PRICE_HISTORY Table - Historical OHLC candle data
CREATE TABLE IF NOT EXISTS price_history (
    id SERIAL PRIMARY KEY,
    asset_id INTEGER NOT NULL,
    open_price DECIMAL(20, 8) NOT NULL,
    high_price DECIMAL(20, 8) NOT NULL,
    low_price DECIMAL(20, 8) NOT NULL,
    close_price DECIMAL(20, 8) NOT NULL,
    volume DECIMAL(30, 2),
    interval VARCHAR(10) NOT NULL, -- e.g., '1m', '5m', '1h', '1d'
    timestamp TIMESTAMP NOT NULL,
    
    CONSTRAINT fk_price_history_asset FOREIGN KEY (asset_id) 
        REFERENCES assets(id) ON DELETE CASCADE,
    CONSTRAINT unique_asset_interval_timestamp UNIQUE(asset_id, interval, timestamp)
);

CREATE INDEX idx_history_asset_time ON price_history(asset_id, interval, timestamp);

-- 4. ORDER_BOOK Table - User orders (buy/sell orders)
-- Note: user_id references users in MySQL (no FK here, managed by application layer)
CREATE TABLE IF NOT EXISTS order_book (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,              -- User ID from MySQL database (no FK)
    asset_id INTEGER NOT NULL,
    side VARCHAR(10) NOT NULL,             -- 'bid' (buy) or 'ask' (sell)
    price DECIMAL(20, 8) NOT NULL,
    quantity DECIMAL(30, 8) NOT NULL,
    quantity_filled DECIMAL(30, 8) DEFAULT 0,  -- Amount already traded
    total DECIMAL(30, 8) NOT NULL,         -- price * quantity
    status VARCHAR(20) DEFAULT 'open',     -- 'open', 'filled', 'cancelled', 'partial'
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT fk_order_book_asset FOREIGN KEY (asset_id) 
        REFERENCES assets(id) ON DELETE CASCADE,
    CONSTRAINT check_order_side CHECK (side IN ('bid', 'ask')),
    CONSTRAINT check_order_status CHECK (status IN ('open', 'filled', 'cancelled', 'partial'))
);

CREATE INDEX idx_orderbook_user ON order_book(user_id, status);
CREATE INDEX idx_orderbook_asset_side ON order_book(asset_id, side, price);
CREATE INDEX idx_orderbook_status ON order_book(status, created_at);

-- 5. WATCHLIST Table - User's favorite assets
-- Note: user_id references users in MySQL database (no FK here, managed by application layer)
CREATE TABLE IF NOT EXISTS watchlist (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,              -- User ID from MySQL database (no FK)
    asset_id INTEGER NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT fk_watchlist_asset FOREIGN KEY (asset_id) 
        REFERENCES assets(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_asset_watchlist UNIQUE(user_id, asset_id)
);

CREATE INDEX idx_watchlist_user ON watchlist(user_id);

-- 6. PRICE_ALERTS Table - Price notifications
-- Note: user_id references users in MySQL database (no FK here, managed by application layer)
CREATE TABLE IF NOT EXISTS price_alerts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,              -- User ID from MySQL database (no FK)
    asset_id INTEGER NOT NULL,
    condition VARCHAR(10) NOT NULL,
    target_price DECIMAL(20, 8) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    triggered_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT fk_price_alerts_asset FOREIGN KEY (asset_id) 
        REFERENCES assets(id) ON DELETE CASCADE,
    CONSTRAINT check_alert_condition CHECK (condition IN ('above', 'below')),
    CONSTRAINT check_alert_status CHECK (status IN ('active', 'triggered', 'cancelled'))
);

CREATE INDEX idx_alerts_active ON price_alerts(asset_id, status) 
    WHERE status = 'active';

-- Insert sample data
INSERT INTO assets (symbol, name, asset_type, exchange) VALUES
('BTC', 'Bitcoin', 'crypto', 'Binance'),
('ETH', 'Ethereum', 'crypto', 'Binance'),
('SOL', 'Solana', 'crypto', 'Binance'),
('AAPL', 'Apple Inc.', 'stock', 'NASDAQ'),
('TSLA', 'Tesla Inc.', 'stock', 'NASDAQ'),
('GOOGL', 'Alphabet Inc.', 'stock', 'NASDAQ'),
('EUR/USD', 'Euro/US Dollar', 'forex', 'Forex'),
('GBP/USD', 'British Pound/US Dollar', 'forex', 'Forex')
ON CONFLICT (symbol) DO NOTHING;
