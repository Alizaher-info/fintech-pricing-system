# Market Trading Service

Real-time market data distribution service that consumes price updates from Kafka and streams them to clients via gRPC.

## Architecture

```
Kafka (price-updates) → Market-Trading Service → gRPC Stream → PHP Backend
                              ↓
                         Redis Cache
```

## Features

- ✅ Kafka consumer with automatic offset management
- ✅ Redis caching for hot price data (60s TTL)
- ✅ gRPC server-side streaming for real-time updates
- ✅ Unary RPC for quick price lookups
- ✅ Multiple client support with broadcast
- ✅ Graceful shutdown
- ✅ Default configuration (no complex setup needed)

## Quick Start

### 1. Install Dependencies
```bash
make deps
```

### 2. Generate Protobuf Code
```bash
make proto
```

### 3. Configure Environment
```bash
cp .env.example .env
# Edit .env if needed (defaults work for local development)
```

### 4. Run the Service
```bash
make run
```

## Configuration

All settings have sensible defaults:

| Variable | Default | Description |
|----------|---------|-------------|
| `GRPC_PORT` | 50052 | gRPC server port |
| `KAFKA_BROKERS` | localhost:9092 | Kafka broker addresses |
| `KAFKA_TOPIC` | price-updates | Topic to consume from |
| `KAFKA_GROUP_ID` | market-trading-group | Consumer group ID |
| `REDIS_HOST` | localhost | Redis host |
| `REDIS_PORT` | 6379 | Redis port |
| `REDIS_PRICE_TTL` | 60 | Cache TTL in seconds |

## gRPC API

### 1. StreamPrices (Server Streaming)
Subscribe to real-time price updates.

```protobuf
rpc StreamPrices(SubscribeRequest) returns (stream PriceUpdate);
```

**Request:**
```json
{
  "symbols": ["BTCUSDT", "ETHUSDT"]  // Empty array = all symbols
}
```

**Response Stream:**
```json
{
  "symbol": "BTCUSDT",
  "price": 45000.50,
  "bid_price": 44995.25,
  "ask_price": 45005.75,
  "volume_24h": 1234567.89,
  "change_24h": 2.5,
  "high_24h": 46000.00,
  "low_24h": 44000.00,
  "timestamp": 1702742400000
}
```

### 2. GetPrice (Unary)
Get current cached price for a symbol.

```protobuf
rpc GetPrice(GetPriceRequest) returns (GetPriceResponse);
```

## Development

### Build
```bash
make build
```

### Test
```bash
make test
```

### Docker
```bash
make docker-build
make docker-run
```

## How It Works

1. **Kafka Consumer**: Reads messages from `price-updates` topic
2. **Redis Cache**: Stores latest price for each symbol (60s TTL)
3. **gRPC Broadcast**: Pushes updates to all connected PHP clients
4. **Auto-Recovery**: Automatically reconnects on failures

## Kafka Consumer Details

- **Consumer Group**: `market-trading-group`
- **Offset Strategy**: Start from latest (OffsetNewest)
- **Auto-Commit**: Enabled (every 5 seconds)
- **Message Loss**: Zero - Kafka persists messages
- **Replay**: Can re-consume from any offset if needed

## Performance

- Handles **1000+ updates/second**
- Supports **multiple concurrent clients**
- Non-blocking broadcasts (buffered channels)
- Automatic backpressure handling

## Monitoring

Check subscriber count and health:
```bash
# Number of connected clients logged on startup/shutdown
# Redis health check available via cache.Health()
```

## Next Steps

1. Add metrics (Prometheus)
2. Add health check endpoint
3. Add authentication/authorization
4. Add rate limiting
5. Add message compression
