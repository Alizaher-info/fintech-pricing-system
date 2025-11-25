# Trading Platform Database Schema Guide

**Date:** November 22, 2025  
**Project:** Fintech Pricing System - Trading Platform  
**Database:** PostgreSQL  
**Purpose:** Complete reference for all 6 database tables

---

## Table of Contents
1. [ASSETS Table](#1-assets-table)
2. [PRICE_QUOTES Table](#2-price_quotes-table)
3. [PRICE_HISTORY Table](#3-price_history-table)
4. [ORDER_BOOK Table](#4-order_book-table)
5. [WATCHLIST Table](#5-watchlist-table)
6. [PRICE_ALERTS Table](#6-price_alerts-table)
7. [Complete Schema SQL](#complete-schema-sql)
8. [Glossary](#glossary)

---

## 1. ASSETS Table

### Purpose
Master table containing all tradable assets (cryptocurrencies, stocks, forex pairs).

### Table Structure
```sql
CREATE TABLE assets (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(20) UNIQUE NOT NULL,     -- e.g., BTC, AAPL, EUR/USD
    name VARCHAR(100) NOT NULL,             -- e.g., Bitcoin, Apple Inc.
    asset_type VARCHAR(20) NOT NULL,        -- 'crypto', 'stock', 'forex'
    exchange VARCHAR(50),                   -- e.g., Binance, NYSE, Forex
    is_active BOOLEAN DEFAULT true,         -- Can it be traded?
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CHECK (asset_type IN ('crypto', 'stock', 'forex'))
);

CREATE INDEX idx_assets_type ON assets(asset_type);
CREATE INDEX idx_assets_symbol ON assets(symbol);
```

### Real Examples
```
┌────┬────────┬─────────────────┬────────────┬──────────┬───────────┐
│ id │ symbol │ name            │ asset_type │ exchange │ is_active │
├────┼────────┼─────────────────┼────────────┼──────────┼───────────┤
│ 1  │ BTC    │ Bitcoin         │ crypto     │ Binance  │ true      │
│ 2  │ ETH    │ Ethereum        │ crypto     │ Binance  │ true      │
│ 3  │ AAPL   │ Apple Inc.      │ stock      │ NASDAQ   │ true      │
│ 4  │ TSLA   │ Tesla Inc.      │ stock      │ NASDAQ   │ true      │
│ 5  │ EUR/USD│ Euro/US Dollar  │ forex      │ Forex    │ true      │
└────┴────────┴─────────────────┴────────────┴──────────┴───────────┘
```

### Common Queries
```sql
-- Get all active cryptocurrencies
SELECT * FROM assets 
WHERE asset_type = 'crypto' AND is_active = true;

-- Get all stocks on NASDAQ
SELECT * FROM assets 
WHERE asset_type = 'stock' AND exchange = 'NASDAQ';

-- Search for an asset by symbol
SELECT * FROM assets 
WHERE symbol = 'BTC';
```

### Key Points
- **One asset = one row** (Bitcoin is one row, Apple is one row)
- **symbol is UNIQUE** (can't have two Bitcoin entries)
- **is_active** lets you disable assets without deleting them
- **asset_type** helps filter: crypto vs stocks vs forex

---

## 2. PRICE_QUOTES Table

### Purpose
Stores **real-time** current prices for all assets. This is what users see when they check "current price".

### Table Structure
```sql
CREATE TABLE price_quotes (
    id SERIAL PRIMARY KEY,
    asset_id INTEGER NOT NULL,
    current_price DECIMAL(20, 8) NOT NULL,  -- Current market price
    bid_price DECIMAL(20, 8),               -- Highest buy order
    ask_price DECIMAL(20, 8),               -- Lowest sell order
    volume_24h DECIMAL(30, 2),              -- Trading volume (24 hours)
    change_24h DECIMAL(10, 4),              -- % change (e.g., +2.5%)
    high_24h DECIMAL(20, 8),                -- Highest price in 24h
    low_24h DECIMAL(20, 8),                 -- Lowest price in 24h
    market_cap DECIMAL(30, 2),              -- Total value (crypto/stocks)
    source VARCHAR(50),                     -- API source (Binance, CoinGecko)
    last_updated TIMESTAMP DEFAULT NOW(),   -- When price was fetched
    
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    UNIQUE(asset_id)  -- One row per asset (updates in place)
);

CREATE INDEX idx_quotes_updated ON price_quotes(last_updated);
```

### Real Examples
```
Bitcoin Real-Time Data:
├─ Current Price: $95,000.00
├─ Bid Price: $94,995.00 (someone wants to buy at this price)
├─ Ask Price: $95,005.00 (someone wants to sell at this price)
├─ Volume 24h: $45,000,000,000 (total traded in 24 hours)
├─ Change 24h: +2.5% (price went up 2.5%)
├─ High 24h: $95,500.00 (peak price today)
├─ Low 24h: $92,000.00 (lowest price today)
└─ Source: Binance API
```

### Understanding Bid/Ask/Spread
```
Price: $100.00
       |
┌──────┼──────┐
│ BID  │ ASK  │
│$99.95│$100.05│ ← Spread = $0.10
└──────┴──────┘

BID = Buyers want to pay $99.95
ASK = Sellers want to receive $100.05
SPREAD = $0.10 (the difference)

If you BUY → You pay $100.05 (ask price)
If you SELL → You get $99.95 (bid price)
```

### Why Spread Matters
1. **Trading Cost:** You lose the spread amount instantly
2. **Market Health:** Small spread = liquid market (good!)
3. **Strategy Impact:** Day traders need tiny spreads to profit

### Common Queries
```sql
-- Get current Bitcoin price
SELECT current_price, change_24h, volume_24h
FROM price_quotes pq
JOIN assets a ON pq.asset_id = a.id
WHERE a.symbol = 'BTC';

-- Get top 10 gainers (24h)
SELECT a.symbol, a.name, pq.current_price, pq.change_24h
FROM price_quotes pq
JOIN assets a ON pq.asset_id = a.id
WHERE a.is_active = true
ORDER BY pq.change_24h DESC
LIMIT 10;

-- Calculate spread
SELECT 
    a.symbol,
    pq.bid_price,
    pq.ask_price,
    (pq.ask_price - pq.bid_price) as spread,
    ((pq.ask_price - pq.bid_price) / pq.current_price * 100) as spread_percentage
FROM price_quotes pq
JOIN assets a ON pq.asset_id = a.id;
```

### Key Points
- **One row per asset** (updates existing row, doesn't create new ones)
- **Updates every 30-60 seconds** from API
- **bid/ask prices** show market depth
- **source field** tracks where data came from (important for debugging)

---

## 3. PRICE_HISTORY Table

### Purpose
Stores **historical price data** as OHLC candles for charting and analysis.

### Table Structure
```sql
CREATE TABLE price_history (
    id SERIAL PRIMARY KEY,
    asset_id INTEGER NOT NULL,
    open_price DECIMAL(20, 8) NOT NULL,     -- Price at interval start
    high_price DECIMAL(20, 8) NOT NULL,     -- Highest price in interval
    low_price DECIMAL(20, 8) NOT NULL,      -- Lowest price in interval
    close_price DECIMAL(20, 8) NOT NULL,    -- Price at interval end
    volume DECIMAL(30, 2),                  -- Trading volume
    interval VARCHAR(10) NOT NULL,          -- '1m', '5m', '1h', '1d'
    timestamp TIMESTAMP NOT NULL,           -- When this candle started
    
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    UNIQUE(asset_id, interval, timestamp)
);

CREATE INDEX idx_history_asset_time ON price_history(asset_id, interval, timestamp);
```

### Understanding OHLC Candles
```
Each "candle" represents price movement in a time period:

        │  high ($95,500) ← Highest price reached
        ┃
      ┌─╂─┐
      │ ┃ │ open ($95,000) ← Price at start
      │ ┃ │ close ($95,200) ← Price at end
      └─╂─┘
        ┃
        │  low ($94,800) ← Lowest price reached

Candle Color:
🟢 Green = close > open (price went up!)
🔴 Red = close < open (price went down!)
```

### Real Example - Bitcoin 1-Hour Candles
```
Timestamp        │ Open      │ High      │ Low       │ Close     │ Volume
─────────────────┼───────────┼───────────┼───────────┼───────────┼──────────
2025-11-22 10:00 │ $95,000   │ $95,500   │ $94,800   │ $95,200   │ $10M  🟢
2025-11-22 11:00 │ $95,200   │ $95,800   │ $95,100   │ $95,600   │ $12M  🟢
2025-11-22 12:00 │ $95,600   │ $96,000   │ $95,400   │ $95,500   │ $8M   🔴
2025-11-22 13:00 │ $95,500   │ $95,700   │ $94,900   │ $95,000   │ $15M  🔴
```

### Interval Types
```
1m  = 1 minute   → Day traders (very short-term)
5m  = 5 minutes  → Active traders
15m = 15 minutes → Swing traders
1h  = 1 hour     → Daily analysis
4h  = 4 hours    → Medium-term trends
1d  = 1 day      → Long-term investors
1w  = 1 week     → Big picture view
```

### Common Queries
```sql
-- Get last 24 hours of hourly candles for Bitcoin
SELECT 
    timestamp,
    open_price,
    high_price,
    low_price,
    close_price,
    volume
FROM price_history ph
JOIN assets a ON ph.asset_id = a.id
WHERE a.symbol = 'BTC' 
  AND ph.interval = '1h'
  AND ph.timestamp > NOW() - INTERVAL '24 hours'
ORDER BY ph.timestamp ASC;

-- Get daily candles for last 30 days
SELECT 
    DATE(timestamp) as date,
    open_price,
    high_price,
    low_price,
    close_price
FROM price_history ph
JOIN assets a ON ph.asset_id = a.id
WHERE a.symbol = 'ETH'
  AND ph.interval = '1d'
  AND ph.timestamp > NOW() - INTERVAL '30 days'
ORDER BY ph.timestamp DESC;

-- Calculate daily price change
SELECT 
    timestamp,
    close_price,
    LAG(close_price) OVER (ORDER BY timestamp) as prev_close,
    ((close_price - LAG(close_price) OVER (ORDER BY timestamp)) 
     / LAG(close_price) OVER (ORDER BY timestamp) * 100) as change_percent
FROM price_history
WHERE asset_id = 1 AND interval = '1d';
```

### Key Points
- **Multiple rows per asset** (one for each time interval)
- **Used for charts** (TradingView-style candlestick charts)
- **UNIQUE constraint** prevents duplicate candles
- **Different intervals** for different trading strategies

---

## 4. ORDER_BOOK Table

### Purpose
Shows **all active buy/sell orders** at different prices (market depth).

### Table Structure
```sql
CREATE TABLE order_book (
    id SERIAL PRIMARY KEY,
    asset_id INTEGER NOT NULL,
    side VARCHAR(10) NOT NULL,              -- 'bid' (buy) or 'ask' (sell)
    price DECIMAL(20, 8) NOT NULL,          -- Order price
    quantity DECIMAL(30, 8) NOT NULL,       -- Amount to buy/sell
    total DECIMAL(30, 8) NOT NULL,          -- price × quantity
    timestamp TIMESTAMP DEFAULT NOW(),
    
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    CHECK (side IN ('bid', 'ask'))
);

CREATE INDEX idx_orderbook_asset_side ON order_book(asset_id, side, price);
```

### Understanding Order Book
```
BITCOIN ORDER BOOK

ASKS (Sell Orders - people want to sell):
Price      │ Quantity │ Total
───────────┼──────────┼──────────
$95,100.00 │ 0.5 BTC  │ $47,550   ← Cheapest seller
$95,150.00 │ 1.2 BTC  │ $114,180
$95,200.00 │ 0.8 BTC  │ $76,160
─────────────────────────────────
        SPREAD = $100
─────────────────────────────────
BIDS (Buy Orders - people want to buy):
Price      │ Quantity │ Total
───────────┼──────────┼──────────
$95,000.00 │ 2.0 BTC  │ $190,000  ← Highest buyer
$94,950.00 │ 1.5 BTC  │ $142,425
$94,900.00 │ 3.0 BTC  │ $284,700

If you SELL now → You get $95,000 (best bid)
If you BUY now → You pay $95,100 (best ask)
```

### Market Depth Concept
```
Deep Market (Liquid):
Bid: $100.00 (1000 BTC)
Bid: $99.95  (500 BTC)
Bid: $99.90  (300 BTC)
└─ Lots of orders = stable price ✅

Shallow Market (Illiquid):
Bid: $100.00 (1 BTC)
Bid: $95.00  (0.5 BTC)
Bid: $90.00  (0.2 BTC)
└─ Few orders = price can crash fast! ❌
```

### Common Queries
```sql
-- Get order book for Bitcoin (top 10 each side)
SELECT 
    side,
    price,
    quantity,
    total
FROM order_book
WHERE asset_id = 1
ORDER BY 
    CASE WHEN side = 'ask' THEN price END ASC,  -- Asks: lowest first
    CASE WHEN side = 'bid' THEN price END DESC  -- Bids: highest first
LIMIT 20;

-- Calculate market depth (total volume at each side)
SELECT 
    side,
    SUM(quantity) as total_quantity,
    SUM(total) as total_value,
    COUNT(*) as num_orders
FROM order_book
WHERE asset_id = 1
GROUP BY side;

-- Get best bid and ask
SELECT 
    MAX(CASE WHEN side = 'bid' THEN price END) as best_bid,
    MIN(CASE WHEN side = 'ask' THEN price END) as best_ask,
    (MIN(CASE WHEN side = 'ask' THEN price END) - 
     MAX(CASE WHEN side = 'bid' THEN price END)) as spread
FROM order_book
WHERE asset_id = 1;
```

### Key Points
- **Shows market depth** (how many people want to buy/sell)
- **Updates frequently** (every few seconds)
- **Best bid/ask** determines current spread
- **Deep order book** = healthy, liquid market

---

## 5. WATCHLIST Table

### Purpose
Lets users **save favorite assets** for quick access.

### Table Structure
```sql
CREATE TABLE watchlist (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,               -- Which user
    asset_id INTEGER NOT NULL,              -- Which asset
    notes TEXT,                             -- User's personal notes
    created_at TIMESTAMP DEFAULT NOW(),
    
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    UNIQUE(user_id, asset_id)              -- Can't add same asset twice
);

CREATE INDEX idx_watchlist_user ON watchlist(user_id);
```

### Real Example
```
User ID: 123 (John's Watchlist)

┌────┬────────┬─────────────┬──────────────────────┬────────────┐
│ id │ symbol │ name        │ notes                │ added_date │
├────┼────────┼─────────────┼──────────────────────┼────────────┤
│ 1  │ BTC    │ Bitcoin     │ Buy if drops to $90k │ 2025-11-20 │
│ 2  │ ETH    │ Ethereum    │ Long-term hold       │ 2025-11-19 │
│ 3  │ TSLA   │ Tesla       │ Check earnings       │ 2025-11-18 │
│ 4  │ AAPL   │ Apple       │ Dividend stock       │ 2025-11-15 │
└────┴────────┴─────────────┴──────────────────────┴────────────┘
```

### Frontend Display
```
╔═══════════════════════════════════════════════════╗
║          MY WATCHLIST (4 assets)                  ║
╠═══════════════════════════════════════════════════╣
║ [BTC] Bitcoin                          $95,000.00 ║
║   Crypto • +2.5% (24h) • "Buy at $90k"            ║
║   [Remove] [View Chart] [Set Alert]               ║
║───────────────────────────────────────────────────║
║ [ETH] Ethereum                          $3,200.00 ║
║   Crypto • +1.8% (24h) • "Long-term hold"         ║
║   [Remove] [View Chart] [Set Alert]               ║
╚═══════════════════════════════════════════════════╝
```

### Common Queries
```sql
-- Get user's watchlist with current prices
SELECT 
    w.id,
    a.symbol,
    a.name,
    a.asset_type,
    w.notes,
    pq.current_price,
    pq.change_24h,
    w.created_at
FROM watchlist w
JOIN assets a ON w.asset_id = a.id
LEFT JOIN price_quotes pq ON a.id = pq.asset_id
WHERE w.user_id = 123
ORDER BY w.created_at DESC;

-- Add asset to watchlist
INSERT INTO watchlist (user_id, asset_id, notes)
VALUES (123, 1, 'Waiting to buy at $90k')
ON CONFLICT (user_id, asset_id) DO NOTHING;

-- Remove from watchlist
DELETE FROM watchlist
WHERE user_id = 123 AND asset_id = 1;

-- Update notes
UPDATE watchlist
SET notes = 'Buy if drops below $85k'
WHERE user_id = 123 AND asset_id = 1;
```

### Key Points
- **Like bookmarks** for your favorite assets
- **UNIQUE constraint** prevents duplicate entries
- **Notes field** for personal reminders
- **Quick access** to assets you care about

---

## 6. PRICE_ALERTS Table

### Purpose
Sends **notifications** when price reaches a target.

### Table Structure
```sql
CREATE TABLE price_alerts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,              -- Which user
    asset_id INTEGER NOT NULL,             -- Which asset to monitor
    condition VARCHAR(10) NOT NULL,        -- 'above' or 'below'
    target_price DECIMAL(20, 8) NOT NULL,  -- Target price
    status VARCHAR(20) DEFAULT 'active',   -- 'active', 'triggered', 'cancelled'
    triggered_at TIMESTAMP,                -- When it triggered
    created_at TIMESTAMP DEFAULT NOW(),
    
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    CHECK (condition IN ('above', 'below')),
    CHECK (status IN ('active', 'triggered', 'cancelled'))
);

CREATE INDEX idx_alerts_active ON price_alerts(asset_id, status) 
    WHERE status = 'active';
```

### Real Examples

#### Example 1: Buy Alert
```
User wants to BUY Bitcoin at discount:

Alert:
├─ Asset: Bitcoin (BTC)
├─ Condition: BELOW
├─ Target: $90,000
├─ Current: $95,000
└─ Status: active

When price drops to $89,999 → 🔔 TRIGGERED!
"Bitcoin is now $89,999 - Time to buy!"
```

#### Example 2: Sell Alert
```
User owns Ethereum, wants to take profit:

Alert:
├─ Asset: Ethereum (ETH)
├─ Condition: ABOVE
├─ Target: $3,500
├─ Current: $3,200
└─ Status: active

When price rises to $3,501 → 🔔 TRIGGERED!
"Ethereum reached $3,501 - Time to sell!"
```

#### Example 3: Ladder Strategy
```
Bitcoin Trading Strategy:

Alert 1: BELOW $90,000  → "Buy opportunity!"
Alert 2: ABOVE $100,000 → "Take profit - sell 50%"
Alert 3: ABOVE $110,000 → "Sell another 25%"
Alert 4: BELOW $85,000  → "Stop loss - sell all!"
```

### Alert Status Flow
```
[User Creates Alert]
        ↓
   status = 'active'
        ↓
[Background Service checks every 60 seconds]
        ↓
   Price matches condition?
        ↓
    ┌───NO───┐     ┌───YES───┐
    │        │     │         │
Continue     status = 'triggered'
monitoring        ↓
           Send notification
           (Email/SMS/Push)
                ↓
           User dismissed
```

### Common Queries
```sql
-- Get user's active alerts
SELECT 
    pa.id,
    a.symbol,
    a.name,
    pa.condition,
    pa.target_price,
    pq.current_price,
    (pa.target_price - pq.current_price) as distance,
    pa.created_at
FROM price_alerts pa
JOIN assets a ON pa.asset_id = a.id
LEFT JOIN price_quotes pq ON a.id = pq.asset_id
WHERE pa.user_id = 123 AND pa.status = 'active'
ORDER BY ABS(pa.target_price - pq.current_price) ASC;  -- Closest first

-- Check which alerts should trigger
SELECT 
    pa.id,
    pa.user_id,
    pa.condition,
    pa.target_price,
    a.symbol,
    pq.current_price
FROM price_alerts pa
JOIN assets a ON pa.asset_id = a.id
JOIN price_quotes pq ON a.id = pq.asset_id
WHERE pa.status = 'active'
  AND (
      (pa.condition = 'above' AND pq.current_price >= pa.target_price) OR
      (pa.condition = 'below' AND pq.current_price <= pa.target_price)
  );

-- Trigger alert
UPDATE price_alerts
SET 
    status = 'triggered',
    triggered_at = NOW()
WHERE id = 123;

-- Get triggered alerts (last 7 days)
SELECT 
    a.symbol,
    pa.condition,
    pa.target_price,
    pa.triggered_at
FROM price_alerts pa
JOIN assets a ON pa.asset_id = a.id
WHERE pa.user_id = 123 
  AND pa.status = 'triggered'
  AND pa.triggered_at > NOW() - INTERVAL '7 days'
ORDER BY pa.triggered_at DESC;
```

### Key Points
- **24/7 price monitoring** (you don't have to watch constantly)
- **Multiple alerts per asset** (different price levels)
- **Status tracking** (active/triggered/cancelled)
- **History preserved** (see when alerts fired)

---

## Complete Schema SQL

```sql
-- Drop existing tables (in correct order due to foreign keys)
DROP TABLE IF EXISTS price_alerts CASCADE;
DROP TABLE IF EXISTS watchlist CASCADE;
DROP TABLE IF EXISTS order_book CASCADE;
DROP TABLE IF EXISTS price_history CASCADE;
DROP TABLE IF EXISTS price_quotes CASCADE;
DROP TABLE IF EXISTS assets CASCADE;

-- 1. ASSETS
CREATE TABLE assets (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    asset_type VARCHAR(20) NOT NULL,
    exchange VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CHECK (asset_type IN ('crypto', 'stock', 'forex'))
);

CREATE INDEX idx_assets_type ON assets(asset_type);
CREATE INDEX idx_assets_symbol ON assets(symbol);

-- 2. PRICE_QUOTES
CREATE TABLE price_quotes (
    id SERIAL PRIMARY KEY,
    asset_id INTEGER NOT NULL,
    current_price DECIMAL(20, 8) NOT NULL,
    bid_price DECIMAL(20, 8),
    ask_price DECIMAL(20, 8),
    volume_24h DECIMAL(30, 2),
    change_24h DECIMAL(10, 4),
    high_24h DECIMAL(20, 8),
    low_24h DECIMAL(20, 8),
    market_cap DECIMAL(30, 2),
    source VARCHAR(50),
    last_updated TIMESTAMP DEFAULT NOW(),
    
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    UNIQUE(asset_id)
);

CREATE INDEX idx_quotes_updated ON price_quotes(last_updated);

-- 3. PRICE_HISTORY
CREATE TABLE price_history (
    id SERIAL PRIMARY KEY,
    asset_id INTEGER NOT NULL,
    open_price DECIMAL(20, 8) NOT NULL,
    high_price DECIMAL(20, 8) NOT NULL,
    low_price DECIMAL(20, 8) NOT NULL,
    close_price DECIMAL(20, 8) NOT NULL,
    volume DECIMAL(30, 2),
    interval VARCHAR(10) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    UNIQUE(asset_id, interval, timestamp)
);

CREATE INDEX idx_history_asset_time ON price_history(asset_id, interval, timestamp);

-- 4. ORDER_BOOK
CREATE TABLE order_book (
    id SERIAL PRIMARY KEY,
    asset_id INTEGER NOT NULL,
    side VARCHAR(10) NOT NULL,
    price DECIMAL(20, 8) NOT NULL,
    quantity DECIMAL(30, 8) NOT NULL,
    total DECIMAL(30, 8) NOT NULL,
    timestamp TIMESTAMP DEFAULT NOW(),
    
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    CHECK (side IN ('bid', 'ask'))
);

CREATE INDEX idx_orderbook_asset_side ON order_book(asset_id, side, price);

-- 5. WATCHLIST
CREATE TABLE watchlist (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    asset_id INTEGER NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    UNIQUE(user_id, asset_id)
);

CREATE INDEX idx_watchlist_user ON watchlist(user_id);

-- 6. PRICE_ALERTS
CREATE TABLE price_alerts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    asset_id INTEGER NOT NULL,
    condition VARCHAR(10) NOT NULL,
    target_price DECIMAL(20, 8) NOT NULL,
    status VARCHAR(20) DEFAULT 'active',
    triggered_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    CHECK (condition IN ('above', 'below')),
    CHECK (status IN ('active', 'triggered', 'cancelled'))
);

CREATE INDEX idx_alerts_active ON price_alerts(asset_id, status) 
    WHERE status = 'active';

-- Sample data
INSERT INTO assets (symbol, name, asset_type, exchange) VALUES
('BTC', 'Bitcoin', 'crypto', 'Binance'),
('ETH', 'Ethereum', 'crypto', 'Binance'),
('AAPL', 'Apple Inc.', 'stock', 'NASDAQ'),
('TSLA', 'Tesla Inc.', 'stock', 'NASDAQ'),
('EUR/USD', 'Euro/US Dollar', 'forex', 'Forex');
```

---

## Glossary

### Trading Terms

**Asset**: Anything that can be bought or sold (Bitcoin, Apple stock, EUR/USD).

**Bid Price**: The highest price a buyer is willing to pay.

**Ask Price**: The lowest price a seller is willing to accept.

**Spread**: The difference between bid and ask prices (bid/ask gap).

**OHLC**: Open-High-Low-Close prices for a time period.

**Candle**: Visual representation of OHLC data on a chart.

**Volume**: Amount of asset traded in a time period.

**Market Cap**: Total value of all coins/shares (price × supply).

**Liquidity**: How easy it is to buy/sell without affecting price.

**Order Book**: List of all buy/sell orders at different prices.

**Market Depth**: Total volume of orders at various price levels.

**Interval**: Time period for chart candles (1m, 5m, 1h, 1d).

### Asset Types

**Crypto (Cryptocurrency)**: Digital currency (Bitcoin, Ethereum)
- Trades 24/7
- High volatility
- No central authority

**Stock**: Company shares (Apple, Tesla)
- Trades during market hours (9:30 AM - 4:00 PM EST)
- Moderate volatility
- Regulated by SEC

**Forex (Foreign Exchange)**: Currency pairs (EUR/USD, GBP/JPY)
- Trades 24/5 (Monday-Friday)
- Lower volatility
- Largest market in world

### Alert Conditions

**ABOVE**: Trigger when price goes higher than target.
- Use case: "Notify me if Bitcoin goes ABOVE $100,000"

**BELOW**: Trigger when price goes lower than target.
- Use case: "Notify me if Bitcoin drops BELOW $90,000"

### Alert Statuses

**active**: Alert is waiting to trigger.

**triggered**: Alert condition was met, notification sent.

**cancelled**: User cancelled the alert before it triggered.

---

## Database Relationships

```
┌─────────────┐
│   ASSETS    │ ← Master table (all tradable assets)
└──────┬──────┘
       │
       ├────→ PRICE_QUOTES (1:1) ← Current real-time prices
       │
       ├────→ PRICE_HISTORY (1:many) ← Historical candles
       │
       ├────→ ORDER_BOOK (1:many) ← Bid/ask orders
       │
       ├────→ WATCHLIST (1:many) ← User favorites
       │
       └────→ PRICE_ALERTS (1:many) ← Price notifications
```

**Key Relationships:**
- One asset has ONE current price quote (price_quotes)
- One asset has MANY historical candles (price_history)
- One asset has MANY order book entries (order_book)
- One asset can be in MANY user watchlists (watchlist)
- One asset can have MANY price alerts (price_alerts)

---

## Performance Tips

### Indexing Strategy
```sql
-- Always index foreign keys
CREATE INDEX ON watchlist(asset_id);
CREATE INDEX ON price_alerts(asset_id);

-- Index frequently queried columns
CREATE INDEX ON price_history(timestamp);
CREATE INDEX ON price_quotes(last_updated);

-- Composite indexes for common queries
CREATE INDEX ON price_history(asset_id, interval, timestamp);
CREATE INDEX ON order_book(asset_id, side, price);
```

### Query Optimization
```sql
-- Use JOIN instead of subqueries
✅ GOOD:
SELECT a.symbol, pq.current_price
FROM assets a
JOIN price_quotes pq ON a.id = pq.asset_id;

❌ BAD:
SELECT symbol, (SELECT current_price FROM price_quotes WHERE asset_id = a.id)
FROM assets a;

-- Limit result sets
SELECT * FROM price_history 
WHERE timestamp > NOW() - INTERVAL '24 hours'  -- ✅ Filter first
LIMIT 100;  -- ✅ Limit results
```

### Data Retention
```sql
-- Delete old order book data (keep only 24 hours)
DELETE FROM order_book
WHERE timestamp < NOW() - INTERVAL '24 hours';

-- Archive old alerts (keep 30 days)
DELETE FROM price_alerts
WHERE status = 'triggered' 
  AND triggered_at < NOW() - INTERVAL '30 days';

-- Keep price history forever (it's valuable for analysis)
-- But you might aggregate older data:
-- 1 year old: keep only daily candles
-- 5 years old: keep only weekly candles
```

---

## Next Steps for Implementation

### 1. Setup PostgreSQL in Docker
```yaml
# deploy/compose/docker-compose.yml
services:
  postgres:
    image: postgres:16
    environment:
      POSTGRES_DB: trading_platform
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: secure_password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

### 2. Create Migration Files
Create `services/pricing-api/migrations/001_create_schema.sql` with the complete schema.

### 3. Build Data Ingestion Service
- Connect to CoinGecko API (crypto)
- Connect to Yahoo Finance API (stocks)
- Connect to Forex API
- Update `price_quotes` every 60 seconds
- Store historical candles in `price_history`

### 4. Build gRPC Endpoints
```protobuf
service PricingService {
  rpc GetCurrentPrice(GetPriceRequest) returns (PriceResponse);
  rpc GetPriceHistory(HistoryRequest) returns (HistoryResponse);
  rpc GetOrderBook(OrderBookRequest) returns (OrderBookResponse);
  rpc GetWatchlist(WatchlistRequest) returns (WatchlistResponse);
  rpc CreatePriceAlert(AlertRequest) returns (AlertResponse);
}
```

### 5. Build Frontend Components
- `PriceQuoteWidget.vue` - Display current prices
- `PriceChart.vue` - TradingView-style charts
- `OrderBookWidget.vue` - Market depth display
- `WatchlistPanel.vue` - User's favorite assets
- `PriceAlertManager.vue` - Alert creation/management

### 6. Setup Background Jobs
- Price updater (runs every 60 seconds)
- Alert checker (runs every 30 seconds)
- Data cleaner (runs daily at midnight)

---

**This document contains everything you need to reference when building the trading platform database! Keep it handy for future development. 📚**

**Created:** November 22, 2025  
**Version:** 1.0  
**Last Updated:** November 22, 2025
