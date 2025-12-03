# Fintech Pricing System - Development Plan

**Project:** Trading Platform with Real-Time Price Updates  
**Start Date:** December 3, 2025  
**Architecture:** Frontend (Vue.js) ↔ Symfony (PHP) ↔ Go Service  
**Technologies:** WebSocket, gRPC, Kafka, PostgreSQL, MySQL, Redis

---

## 🎯 Project Goals

Build a complete trading platform with:
- ✅ Real-time price updates (WebSocket)
- ✅ User order management
- ✅ Price alerts & watchlist
- ✅ Historical data & charts
- ✅ Scalable microservices architecture

---

## 🏗️ System Architecture

```
┌─────────────────────────────────────────────────────────┐
│  Frontend (Vue.js)                                       │
│  - Real-time charts                                     │
│  - Order management UI                                  │
└────────────────────┬────────────────────────────────────┘
                     │
                     │ WebSocket (real-time)
                     │ REST API (requests)
                     ↓
┌─────────────────────────────────────────────────────────┐
│  Symfony Backend (PHP)                                   │
│  ├─ Authentication (JWT + MySQL)                        │
│  ├─ API Controllers (REST endpoints)                    │
│  ├─ gRPC Client (calls Go service)                      │
│  ├─ Kafka Consumer (receives events)                    │
│  └─ WebSocket Server (broadcasts updates)              │
└────────────────────┬───────────────┬────────────────────┘
                     │               │
                     │ gRPC          │ Kafka
                     │ (sync)        │ (async)
                     ↓               ↓
┌─────────────────────────────────────────────────────────┐
│  Go Pricing Service                                      │
│  ├─ Background Worker (fetch prices every 60s)         │
│  ├─ gRPC Server (serve data to Symfony)                │
│  ├─ Kafka Producer (publish price events)              │
│  ├─ PostgreSQL Operations                              │
│  └─ Internal: Goroutines + Channels                    │
└────────────────────┬───────────────┬────────────────────┘
                     │               │
                     ↓               ↓
              PostgreSQL          Kafka Cluster
              (Trading Data)      (Event Stream)
```

---

## 📊 Technology Stack

| Layer | Technology | Purpose |
|-------|-----------|---------|
| **Frontend** | Vue.js 3 + Vite | User interface, charts |
| **Backend API** | Symfony 7 (PHP 8.3) | REST API, authentication |
| **Pricing Service** | Go 1.21+ | High-performance data processing |
| **Communication** | gRPC + Kafka | Sync requests + async events |
| **Real-Time** | WebSocket (Ratchet) | Live updates to frontend |
| **Databases** | PostgreSQL + MySQL | Trading data + user data |
| **Cache** | Redis | Session storage, message broker |
| **Orchestration** | Docker Compose | Development environment |

---

## 📅 Development Timeline

**Total Duration:** 6-7 weeks  
**Approach:** Incremental development with testing after each phase

---

# Phase 1: Foundation (Database & Infrastructure) ✅

**Duration:** 3-4 days  
**Status:** ✅ COMPLETED (December 3, 2025)

## ✅ Step 1.1: PostgreSQL Database Setup
**Status:** DONE ✅

**Completed Tasks:**
- [x] Created migration files (000001_create_trading_schema.up.sql)
- [x] Set up PostgreSQL in Docker Compose
- [x] Created 6 tables:
  - assets (8 sample records)
  - price_quotes
  - price_history
  - order_book (with user_id)
  - watchlist (with user_id)
  - price_alerts (with user_id)
- [x] Applied migrations successfully
- [x] Verified with TablePlus connection

**Database Connection:**
- Host: localhost
- Port: 5432
- Database: trading_platform
- User: app
- Password: app

**Deliverables:**
- ✅ PostgreSQL running in Docker
- ✅ All tables created with proper indexes
- ✅ Sample data inserted
- ✅ Connection tested via TablePlus

---

# Phase 2: Go Price Fetcher Service (Background Worker)

**Duration:** 1 week (5-7 days)  
**Status:** 🔄 NEXT PHASE  
**Goal:** Automatically fetch cryptocurrency prices and save to PostgreSQL

---

## Step 2.1: Project Structure Setup

**Duration:** 1 day  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Update `services/go.mod` with dependencies:
  - github.com/lib/pq (PostgreSQL driver)
  - github.com/joho/godotenv (environment config)
  - github.com/superoo7/go-gecko/v3 (CoinGecko API)
  - github.com/robfig/cron/v3 (scheduled tasks)
- [ ] Create folder structure:
  ```
  services/pricing-api/
  ├─ cmd/
  │  └─ server/
  │     └─ main.go                    # Entry point
  ├─ internal/
  │  ├─ config/
  │  │  └─ config.go                 # Configuration loader
  │  ├─ database/
  │  │  └─ postgres.go               # Database connection
  │  ├─ fetcher/
  │  │  ├─ coingecko.go              # CoinGecko API client
  │  │  └─ price_fetcher.go          # Price fetching logic
  │  ├─ models/
  │  │  └─ price.go                  # Data structures
  │  └─ repository/
  │     └─ price_repository.go       # Database operations
  └─ pkg/
     └─ logger/
        └─ logger.go                  # Logging utility
  ```
- [ ] Create `.env` file for configuration
- [ ] Implement basic config loader
- [ ] Set up logging

### Testing Checklist:
- [ ] `go build` compiles successfully
- [ ] Can load environment variables from `.env`
- [ ] Can connect to PostgreSQL (connection test)
- [ ] Can read sample assets from database
- [ ] Logs are formatted correctly

### Expected Output:
```
2025-12-03 10:00:00 [INFO] Starting Pricing API Service
2025-12-03 10:00:00 [INFO] Loaded configuration from .env
2025-12-03 10:00:01 [INFO] Connected to PostgreSQL: trading_platform
2025-12-03 10:00:01 [INFO] Found 8 assets in database
```

---

## Step 2.2: CoinGecko API Integration

**Duration:** 2 days  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Install CoinGecko Go client library
- [ ] Create CoinGecko API wrapper
- [ ] Implement rate limiting (50 calls/minute free tier)
- [ ] Fetch single asset price (BTC test)
- [ ] Fetch multiple assets in batch
- [ ] Map asset symbols to CoinGecko IDs:
  - BTC → bitcoin
  - ETH → ethereum
  - SOL → solana
- [ ] Handle API errors and retries
- [ ] Parse response data (price, volume, market cap, etc.)

### Testing Checklist:
- [ ] Successfully fetch Bitcoin current price
- [ ] Successfully fetch all 8 assets
- [ ] API rate limiting works (doesn't exceed 50/min)
- [ ] Handles API failures gracefully (timeout, 429 error)
- [ ] Logs API requests and responses
- [ ] Data structure contains all required fields:
  - current_price
  - bid_price / ask_price (from order book)
  - volume_24h
  - change_24h
  - high_24h / low_24h
  - market_cap

### Expected Output:
```
2025-12-03 10:00:00 [INFO] Fetching prices from CoinGecko...
2025-12-03 10:00:01 [INFO] BTC: $95,000.00 | 24h: +2.5% | Vol: $45B
2025-12-03 10:00:01 [INFO] ETH: $3,200.00 | 24h: +1.8% | Vol: $15B
2025-12-03 10:00:01 [INFO] SOL: $100.00 | 24h: +5.2% | Vol: $2B
2025-12-03 10:00:02 [INFO] Successfully fetched 8 assets
```

---

## Step 2.3: Database Repository Layer

**Duration:** 1 day  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Create `PriceRepository` interface
- [ ] Implement `UpsertPrice` function (INSERT or UPDATE)
- [ ] Implement high/low 24h calculation logic
- [ ] Handle database transactions
- [ ] Add error handling and logging
- [ ] Create helper functions:
  - GetAssetBySymbol
  - GetAllActiveAssets
  - GetLatestPrice

### Database Operation:
```sql
-- Upsert logic (INSERT or UPDATE)
INSERT INTO price_quotes (
    asset_id, current_price, bid_price, ask_price,
    volume_24h, change_24h, high_24h, low_24h,
    market_cap, source, last_updated
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW())
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
    last_updated = NOW();
```

### Testing Checklist:
- [ ] Insert new price quote (first time)
- [ ] Update existing price quote (second time)
- [ ] high_24h updates correctly when price increases
- [ ] low_24h updates correctly when price decreases
- [ ] No duplicate records created
- [ ] Verify in TablePlus: price_quotes table updated
- [ ] Transaction rollback works on error

### Expected Database State:
```
price_quotes table after test:
┌────┬──────────┬───────────┬───────────┬───────────┬─────────────┐
│ id │ asset_id │ current   │ high_24h  │ low_24h   │ last_update │
├────┼──────────┼───────────┼───────────┼───────────┼─────────────┤
│ 1  │ 1 (BTC)  │ $95,000   │ $95,500   │ $91,800   │ 10:00:02    │
│ 2  │ 2 (ETH)  │ $3,200    │ $3,250    │ $3,100    │ 10:00:02    │
└────┴──────────┴───────────┴───────────┴───────────┴─────────────┘
```

---

## Step 2.4: Background Worker with Goroutines

**Duration:** 2 days  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Create main fetcher service
- [ ] Implement goroutines for concurrent fetching
- [ ] Set up channels for communication:
  - fetchChannel (API → Processor)
  - dbChannel (Processor → Database)
- [ ] Add ticker for 60-second interval
- [ ] Implement graceful shutdown (SIGTERM/SIGINT)
- [ ] Add health check endpoint
- [ ] Docker container configuration

### Goroutine Architecture:
```go
// Goroutine 1: Fetch prices
go func() {
    ticker := time.NewTicker(60 * time.Second)
    for range ticker.C {
        prices := fetchFromCoinGecko()
        fetchChannel <- prices
    }
}()

// Goroutine 2: Process and distribute
go func() {
    for prices := range fetchChannel {
        processedPrices := validateAndEnrich(prices)
        dbChannel <- processedPrices
    }
}()

// Goroutine 3: Save to database
go func() {
    for prices := range dbChannel {
        saveToDB(prices)
    }
}()
```

### Testing Checklist:
- [ ] Service starts successfully
- [ ] Fetches prices immediately on startup
- [ ] Continues fetching every 60 seconds
- [ ] Updates PostgreSQL correctly
- [ ] Logs show activity clearly
- [ ] Check TablePlus: price_quotes table updates every 60 seconds
- [ ] Ctrl+C shuts down gracefully (no hanging processes)
- [ ] Service restarts automatically (Docker restart policy)
- [ ] Memory usage stable (no memory leaks)
- [ ] CPU usage reasonable (< 10%)

### Expected Logs:
```
2025-12-03 10:00:00 [INFO] Starting Price Fetcher Service
2025-12-03 10:00:00 [INFO] Goroutine 1: Fetcher started
2025-12-03 10:00:00 [INFO] Goroutine 2: Processor started
2025-12-03 10:00:00 [INFO] Goroutine 3: Database writer started
2025-12-03 10:00:01 [INFO] Fetching prices from CoinGecko...
2025-12-03 10:00:02 [INFO] Processed 8 assets
2025-12-03 10:00:03 [INFO] Saved 8 prices to database
2025-12-03 10:00:03 [INFO] Next fetch in 60 seconds...
2025-12-03 10:01:03 [INFO] Fetching prices from CoinGecko...
```

### Docker Compose Configuration:
```yaml
pricing-api:
  build:
    context: ../..
    dockerfile: services/pricing-api/Dockerfile
  ports:
    - "50051:50051"
  environment:
    POSTGRES_HOST: postgres
    POSTGRES_PORT: 5432
    POSTGRES_DB: trading_platform
    POSTGRES_USER: app
    POSTGRES_PASSWORD: app
  depends_on:
    postgres:
      condition: service_healthy
  restart: unless-stopped
```

### Phase 2 Deliverable:
✅ **Go service running 24/7, automatically updating cryptocurrency prices every 60 seconds**

---

# Phase 3: gRPC Service (API Layer)

**Duration:** 1 week (5-7 days)  
**Status:** ⏳ PENDING  
**Goal:** Enable Symfony backend to request data from Go service

---

## Step 3.1: Protocol Buffers Definition

**Duration:** 1 day  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Install protobuf compiler (protoc)
- [ ] Install Go protobuf plugin
- [ ] Install PHP protobuf plugin
- [ ] Create `.proto` files:
  - `protos/pricing/v1/pricing.proto` (main service)
  - `protos/pricing/v1/asset.proto` (asset messages)
  - `protos/pricing/v1/price.proto` (price messages)
- [ ] Define gRPC service methods:
  - GetCurrentPrice
  - GetPriceHistory
  - GetAssetList
  - GetOrderBook (future)
- [ ] Generate Go code
- [ ] Generate PHP code
- [ ] Create Makefile for code generation

### Proto Definition Example:
```protobuf
syntax = "proto3";

package pricing.v1;

service PricingService {
  rpc GetCurrentPrice(GetPriceRequest) returns (PriceResponse);
  rpc GetPriceHistory(HistoryRequest) returns (HistoryResponse);
  rpc GetAssetList(AssetListRequest) returns (AssetListResponse);
}

message GetPriceRequest {
  string symbol = 1;  // e.g., "BTC"
}

message PriceResponse {
  string symbol = 1;
  double current_price = 2;
  double bid_price = 3;
  double ask_price = 4;
  double volume_24h = 5;
  double change_24h = 6;
  double high_24h = 7;
  double low_24h = 8;
  double market_cap = 9;
  string last_updated = 10;
}
```

### Testing Checklist:
- [ ] Protobuf files compile without errors
- [ ] Go code generated in `gen/pricing/v1/`
- [ ] PHP code generated in `backend/src/Grpc/`
- [ ] No syntax errors in generated code
- [ ] Makefile command works: `make proto-gen`

---

## Step 3.2: gRPC Server Implementation (Go)

**Duration:** 2 days  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Implement gRPC server struct
- [ ] Create service handler methods
- [ ] Connect to PostgreSQL repository
- [ ] Add error handling
- [ ] Add request logging
- [ ] Implement rate limiting
- [ ] Add health check endpoint

### Testing Checklist:
- [ ] gRPC server starts on port 50051
- [ ] Use `grpcurl` to test endpoints:
  ```bash
  grpcurl -plaintext \
    -d '{"symbol":"BTC"}' \
    localhost:50051 \
    pricing.v1.PricingService/GetCurrentPrice
  ```
- [ ] Returns correct data
- [ ] Handles invalid symbol gracefully
- [ ] Response time < 50ms
- [ ] Check logs: requests logged properly

### Expected Response:
```json
{
  "symbol": "BTC",
  "currentPrice": 95000.00,
  "bidPrice": 94995.00,
  "askPrice": 95005.00,
  "volume24h": 45000000000.00,
  "change24h": 2.5,
  "high24h": 95500.00,
  "low24h": 91800.00,
  "marketCap": 1852500000000.00,
  "lastUpdated": "2025-12-03T10:00:00Z"
}
```

---

## Step 3.3: gRPC Client (Symfony)

**Duration:** 2 days  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Install gRPC PHP extension
- [ ] Install Symfony gRPC bundle
- [ ] Create gRPC client service
- [ ] Create factory for connection pooling
- [ ] Create test controller
- [ ] Add error handling
- [ ] Add timeout configuration

### Symfony Service Configuration:
```yaml
# config/services.yaml
services:
  App\Client\PricingGrpcClient:
    arguments:
      $grpcTarget: '%env(PRICING_GRPC_TARGET)%'
```

### Testing Checklist:
- [ ] Symfony can connect to Go service
- [ ] Create test endpoint: `GET /api/test/price/{symbol}`
- [ ] Returns data from Go service
- [ ] Check browser: correct JSON response
- [ ] Response time < 100ms
- [ ] Check logs: gRPC call logged
- [ ] Error handling works (Go service down)

### Expected Response (Symfony):
```json
{
  "symbol": "BTC",
  "price": 95000.00,
  "change_24h": 2.5,
  "volume_24h": 45000000000.00,
  "last_updated": "2025-12-03T10:00:00Z"
}
```

### Phase 3 Deliverable:
✅ **Symfony ↔ Go communication working via gRPC**

---

# Phase 4: Kafka Integration (Event Streaming)

**Duration:** 1 week (5-7 days)  
**Status:** ⏳ PENDING  
**Goal:** Real-time event streaming for price updates

---

## Step 4.1: Kafka Setup

**Duration:** 1 day  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Add Kafka to docker-compose.yml
- [ ] Add Zookeeper (required for Kafka)
- [ ] Add Kafka UI (optional, for monitoring)
- [ ] Start Kafka cluster
- [ ] Create topics: "price-updates"
- [ ] Configure retention policies
- [ ] Test produce/consume

### Docker Compose Addition:
```yaml
zookeeper:
  image: confluentinc/cp-zookeeper:latest
  environment:
    ZOOKEEPER_CLIENT_PORT: 2181

kafka:
  image: confluentinc/cp-kafka:latest
  ports:
    - "9092:9092"
  environment:
    KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
    KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://kafka:9092

kafka-ui:
  image: provectuslabs/kafka-ui:latest
  ports:
    - "8090:8080"
  environment:
    KAFKA_CLUSTERS_0_BOOTSTRAPSERVERS: kafka:9092
```

### Testing Checklist:
- [ ] Kafka container running
- [ ] Zookeeper container running
- [ ] Kafka UI accessible at http://localhost:8090
- [ ] Topic "price-updates" created
- [ ] Can produce test message via UI
- [ ] Can consume test message via UI

---

## Step 4.2: Kafka Producer (Go)

**Duration:** 2 days  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Install Kafka Go library (segmentio/kafka-go)
- [ ] Create Kafka producer service
- [ ] Publish price updates to "price-updates" topic
- [ ] Add goroutine for async publishing
- [ ] Implement retry logic
- [ ] Add message serialization (JSON)
- [ ] Error handling and logging

### Message Format:
```json
{
  "event_type": "price_update",
  "timestamp": "2025-12-03T10:00:00Z",
  "asset": {
    "id": 1,
    "symbol": "BTC",
    "name": "Bitcoin"
  },
  "price": {
    "current": 95000.00,
    "previous": 94500.00,
    "change_percent": 0.53,
    "volume_24h": 45000000000.00
  }
}
```

### Testing Checklist:
- [ ] Go publishes to Kafka successfully
- [ ] Use Kafka UI to see messages in "price-updates" topic
- [ ] Messages contain correct data structure
- [ ] Messages published every 60 seconds (after price fetch)
- [ ] No message loss (delivery confirmation)
- [ ] Check logs: Kafka publish operations logged

---

## Step 4.3: Kafka Consumer (Symfony)

**Duration:** 2 days  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Install Symfony Messenger component
- [ ] Install Kafka transport for Messenger
- [ ] Create message handler
- [ ] Configure consumer service
- [ ] Process price update messages
- [ ] Store in Redis cache
- [ ] Add error handling
- [ ] Consumer runs as background process

### Symfony Configuration:
```yaml
# config/packages/messenger.yaml
framework:
  messenger:
    transports:
      kafka:
        dsn: '%env(KAFKA_DSN)%'
        options:
          topic: 'price-updates'
          group_id: 'symfony-consumers'
```

### Testing Checklist:
- [ ] Symfony consumer starts successfully
- [ ] Receives messages from Kafka
- [ ] Check logs: messages consumed and processed
- [ ] Data parsed correctly
- [ ] Stored in Redis with TTL
- [ ] Consumer runs continuously (no crashes)
- [ ] Handles malformed messages gracefully

### Expected Logs:
```
[2025-12-03 10:00:05] [INFO] Kafka consumer started
[2025-12-03 10:00:05] [INFO] Subscribed to topic: price-updates
[2025-12-03 10:01:05] [INFO] Received price update: BTC $95,100
[2025-12-03 10:01:05] [INFO] Cached price update in Redis
```

### Phase 4 Deliverable:
✅ **Go → Kafka → Symfony pipeline working, price updates streaming in real-time**

---

# Phase 5: WebSocket Real-Time Updates

**Duration:** 1 week (5-7 days)  
**Status:** ⏳ PENDING  
**Goal:** Frontend receives live price updates without page refresh

---

## Step 5.1: WebSocket Server (Symfony)

**Duration:** 2 days  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Install Ratchet WebSocket library
- [ ] Create WebSocket server command
- [ ] Implement JWT authentication
- [ ] Connection management (track connected users)
- [ ] Heartbeat/ping-pong mechanism
- [ ] Graceful disconnection handling
- [ ] Add to docker-compose as separate service

### Testing Checklist:
- [ ] WebSocket server starts on port 8080
- [ ] Can connect from browser console:
  ```javascript
  const ws = new WebSocket('ws://localhost:8080');
  ws.onopen = () => console.log('Connected');
  ```
- [ ] JWT token validation works
- [ ] Unauthorized connections rejected
- [ ] Multiple clients can connect simultaneously
- [ ] Check logs: connections logged

---

## Step 5.2: Kafka → WebSocket Bridge

**Duration:** 2 days  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Integrate Kafka consumer with WebSocket server
- [ ] Broadcast Kafka messages to all connected clients
- [ ] Implement topic-based subscriptions (optional)
- [ ] Handle disconnections during broadcast
- [ ] Add rate limiting per client
- [ ] Message queuing for offline users (optional)

### Message Flow:
```
Go Service → Kafka → Symfony Consumer → WebSocket Server → All Clients
```

### Testing Checklist:
- [ ] Price update arrives in Kafka
- [ ] WebSocket broadcasts to all connected clients
- [ ] Multiple browser tabs receive same message
- [ ] Check browser console: messages received in real-time
- [ ] Disconnected clients don't cause errors
- [ ] New clients receive latest price on connect

---

## Step 5.3: Frontend WebSocket Client

**Duration:** 1 day  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Create WebSocket composable/service
- [ ] Implement connection logic with JWT
- [ ] Handle reconnection on disconnect
- [ ] Update chart in real-time
- [ ] Display connection status indicator
- [ ] Add error notifications

### Testing Checklist:
- [ ] Frontend connects to WebSocket automatically
- [ ] Chart updates without page refresh
- [ ] Connection status shows "Connected" / "Disconnected"
- [ ] Reconnects automatically after disconnect
- [ ] Price changes visible immediately (< 1 second delay)
- [ ] No memory leaks (check DevTools)

### Expected Behavior:
```
User opens app:
├─ Connects to WebSocket
├─ Shows "Connected" indicator (green)
├─ Chart displays current prices
└─ Prices update automatically every 60 seconds
    (smooth animation, no flickering)
```

### Phase 5 Deliverable:
✅ **Real-time price updates working! Users see live prices without refresh.**

---

# Phase 6: Order Management

**Duration:** 1 week (5-7 days)  
**Status:** ⏳ PENDING  
**Goal:** Users can create, view, and cancel orders

---

## Step 6.1: gRPC Order Endpoints (Go)

**Duration:** 2 days  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Define `order.proto` file
- [ ] Implement CreateOrder RPC method
- [ ] Implement GetUserOrders RPC method
- [ ] Implement CancelOrder RPC method
- [ ] Implement GetOrderById RPC method
- [ ] Add order validation logic
- [ ] Connect to order_book table

### Proto Definition:
```protobuf
service OrderService {
  rpc CreateOrder(CreateOrderRequest) returns (OrderResponse);
  rpc GetUserOrders(GetUserOrdersRequest) returns (UserOrdersResponse);
  rpc CancelOrder(CancelOrderRequest) returns (OrderResponse);
}

message CreateOrderRequest {
  int32 user_id = 1;
  string symbol = 2;
  string side = 3;  // "bid" or "ask"
  double price = 4;
  double quantity = 5;
}
```

### Testing Checklist:
- [ ] Can create order via grpcurl
- [ ] Order saved to PostgreSQL (check TablePlus)
- [ ] Can retrieve user's orders
- [ ] Can cancel order (status updated)
- [ ] user_id validation works
- [ ] Invalid data rejected with error

---

## Step 6.2: Symfony Order Controller

**Duration:** 2 days  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Create OrderController
- [ ] POST /api/orders/create endpoint
- [ ] GET /api/orders/mine endpoint
- [ ] DELETE /api/orders/{id} endpoint
- [ ] GET /api/orders/{id} endpoint
- [ ] Extract user_id from JWT token
- [ ] Add input validation
- [ ] Error handling

### Testing Checklist:
- [ ] Create order via Postman:
  ```json
  POST /api/orders/create
  {
    "symbol": "BTC",
    "side": "bid",
    "price": 95000,
    "quantity": 2.0
  }
  ```
- [ ] Check PostgreSQL: order created with correct user_id
- [ ] GET /api/orders/mine returns user's orders only
- [ ] Can cancel order via DELETE
- [ ] User can only cancel their own orders
- [ ] Proper error messages for invalid requests

---

## Step 6.3: Frontend Order UI

**Duration:** 1 day  
**Status:** ⏳ PENDING

### Tasks:
- [ ] Create order placement form
- [ ] Display list of user's orders
- [ ] Add cancel order button
- [ ] Real-time order status updates (via WebSocket)
- [ ] Success/error notifications
- [ ] Form validation

### Testing Checklist:
- [ ] Can place order from UI
- [ ] Order appears in list immediately
- [ ] Can cancel order via UI
- [ ] Success notification shown
- [ ] Form validation prevents invalid input
- [ ] Real-time updates when order status changes

### Phase 6 Deliverable:
✅ **Complete order management: users can place, view, and cancel orders**

---

# Phase 7: Watchlist & Price Alerts

**Duration:** 3-4 days  
**Status:** ⏳ PENDING  
**Goal:** Users can track favorite assets and get price alerts

---

## Step 7.1: Watchlist Feature

**Duration:** 2 days  
**Status:** ⏳ PENDING

### Tasks:
**Go Service:**
- [ ] Define watchlist.proto
- [ ] Implement AddToWatchlist RPC
- [ ] Implement GetWatchlist RPC
- [ ] Implement RemoveFromWatchlist RPC

**Symfony:**
- [ ] Create WatchlistController
- [ ] POST /api/watchlist/{symbol}
- [ ] GET /api/watchlist
- [ ] DELETE /api/watchlist/{symbol}

**Frontend:**
- [ ] Watchlist UI component
- [ ] Add/remove buttons
- [ ] Display with current prices

### Testing Checklist:
- [ ] Add asset to watchlist
- [ ] View watchlist with current prices
- [ ] Remove from watchlist
- [ ] Watchlist persists after logout/login

---

## Step 7.2: Price Alerts

**Duration:** 2 days  
**Status:** ⏳ PENDING

### Tasks:
**Go Service:**
- [ ] Create background worker to check alerts
- [ ] Query active alerts every 60 seconds
- [ ] Compare with current prices
- [ ] Publish alert triggered event to Kafka
- [ ] Update alert status to "triggered"

**Symfony:**
- [ ] Create AlertController
- [ ] POST /api/alerts endpoint
- [ ] GET /api/alerts endpoint
- [ ] DELETE /api/alerts/{id} endpoint
- [ ] Consume alert-triggered events from Kafka
- [ ] Broadcast via WebSocket to user

**Frontend:**
- [ ] Alert creation form
- [ ] Alert list display
- [ ] WebSocket notification handler
- [ ] Browser notification (optional)

### Testing Checklist:
- [ ] Create alert: "Notify if BTC > $100,000"
- [ ] Alert saved to database
- [ ] Alert triggers when condition met
- [ ] WebSocket notification received on frontend
- [ ] Alert status updated to "triggered"
- [ ] Browser notification shown (if enabled)

### Phase 7 Deliverable:
✅ **Complete watchlist and price alert system working**

---

# Phase 8: Testing & Optimization

**Duration:** 3-4 days  
**Status:** ⏳ PENDING  
**Goal:** Production-ready system with tests and optimizations

---

## Step 8.1: Integration Testing

**Duration:** 1 day

### Tasks:
- [ ] Write integration tests for complete user flows:
  - User registration → login → view prices
  - Place order → view order → cancel order
  - Add watchlist → create alert → receive notification
- [ ] Test error scenarios
- [ ] Load testing (simulate 100 concurrent users)
- [ ] API endpoint tests (Postman collection)

---

## Step 8.2: Performance Optimization

**Duration:** 1 day

### Tasks:
- [ ] Database query optimization (add missing indexes)
- [ ] Kafka producer/consumer tuning
- [ ] WebSocket connection pooling
- [ ] Redis caching strategy
- [ ] gRPC connection pooling
- [ ] Go service profiling (memory/CPU)

### Performance Targets:
- API response time: < 100ms (p95)
- gRPC call latency: < 50ms (p95)
- WebSocket message delivery: < 1 second
- Database queries: < 50ms (p95)

---

## Step 8.3: Monitoring & Logging

**Duration:** 1 day

### Tasks:
- [ ] Centralized logging (all services)
- [ ] Health check endpoints
- [ ] Metrics collection
- [ ] Error alerting
- [ ] Dashboard (optional: Grafana)

---

## Phase 8 Deliverable:
✅ **Production-ready system with tests, monitoring, and optimizations**

---

# 📈 Success Metrics

At the end of development, the system should meet these criteria:

### Functional Requirements:
- [x] ✅ PostgreSQL database with 6 tables
- [ ] ⏳ Go service fetches prices automatically every 60 seconds
- [ ] ⏳ Symfony backend serves REST API
- [ ] ⏳ gRPC communication working (Symfony ↔ Go)
- [ ] ⏳ Kafka event streaming working (Go → Symfony)
- [ ] ⏳ WebSocket real-time updates working (Symfony → Frontend)
- [ ] ⏳ Users can place and cancel orders
- [ ] ⏳ Watchlist feature working
- [ ] ⏳ Price alerts working with notifications

### Performance Requirements:
- [ ] API response time < 100ms
- [ ] WebSocket latency < 1 second
- [ ] Price updates every 60 seconds
- [ ] System handles 100 concurrent users
- [ ] No memory leaks
- [ ] CPU usage < 50% under normal load

### Quality Requirements:
- [ ] All features tested manually
- [ ] Integration tests passing
- [ ] No critical bugs
- [ ] Logs are clear and helpful
- [ ] Error handling works properly
- [ ] Code documented

---

# 🛠️ Development Environment

## Required Software:
- [x] ✅ Windows 11
- [x] ✅ Docker Desktop
- [x] ✅ VS Code
- [x] ✅ Git
- [ ] ⏳ Go 1.21+
- [ ] ⏳ Node.js 20+
- [ ] ⏳ PHP 8.3+
- [ ] ⏳ TablePlus (database client)

## VS Code Extensions:
- [ ] Go extension
- [ ] PHP Intelephense
- [ ] Vetur (Vue.js)
- [ ] Docker extension
- [ ] REST Client (for API testing)

---

# 📝 Documentation

As we progress, we'll maintain:

1. **API Documentation:** Endpoints, request/response examples
2. **Architecture Diagrams:** System flow, data flow
3. **Setup Instructions:** How to run locally
4. **Troubleshooting Guide:** Common issues and solutions

---

# 🚀 Next Steps

## Current Status: Phase 1 Completed ✅

**Next: Phase 2 - Go Price Fetcher Service**

**Start with:** Step 2.1 - Project Structure Setup

**Timeline:** Ready to begin immediately

---

## 📞 Communication

After each completed step:
1. ✅ Test and verify
2. ✅ Review together
3. ✅ Document any issues
4. ✅ Commit to git
5. ✅ Move to next step

---

**Last Updated:** December 3, 2025  
**Version:** 1.0  
**Status:** Phase 1 Complete, Phase 2 Ready to Start

---

**Let's build something amazing! 🚀**
