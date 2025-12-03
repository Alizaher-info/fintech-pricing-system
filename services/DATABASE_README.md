# Database Management

This directory contains PostgreSQL database setup and migrations for the Trading Platform.

## Quick Start

### 1. Start PostgreSQL with Docker Compose

```bash
# From project root
cd deploy/compose
docker-compose up -d postgres
```

### 2. Install golang-migrate Tool

**Windows (PowerShell):**
```powershell
# Using Chocolatey
choco install golang-migrate

# Or download binary from:
# https://github.com/golang-migrate/migrate/releases
```

**Linux/Mac:**
```bash
# Using curl
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate
```

### 3. Run Migrations

```bash
cd services
migrate -path migrations -database "postgres://app:app@localhost:5432/trading_platform?sslmode=disable" up
```

## Migration Commands

### Create New Migration

```bash
cd services
migrate create -ext sql -dir migrations -seq create_new_table
```

### Apply Migrations (Up)

```bash
migrate -path migrations -database "$DATABASE_URL" up
```

### Rollback Last Migration (Down)

```bash
migrate -path migrations -database "$DATABASE_URL" down 1
```

### Check Migration Status

```bash
migrate -path migrations -database "$DATABASE_URL" version
```

### Force Migration Version (if stuck)

```bash
migrate -path migrations -database "$DATABASE_URL" force 1
```

## Database Connection

### Local Connection (Outside Docker)

```
Host: localhost
Port: 5432
Database: trading_platform
Username: app
Password: app
```

### Docker Connection String

```
postgres://app:app@postgres:5432/trading_platform?sslmode=disable
```

## Tables Structure

1. **assets** - Master table for all tradable assets (BTC, AAPL, EUR/USD)
2. **price_quotes** - Real-time current prices with bid/ask spreads
3. **price_history** - Historical OHLC candle data for charting
4. **order_book** - Market depth (buy/sell orders)
5. **watchlist** - User's favorite assets
6. **price_alerts** - Price notifications (above/below triggers)

## Useful SQL Commands

### Connect to PostgreSQL

```bash
# Using Docker
docker exec -it trading-postgres psql -U app -d trading_platform

# Using psql locally
psql -h localhost -p 5432 -U app -d trading_platform
```

### Verify Tables

```sql
\dt              -- List all tables
\d assets        -- Describe assets table
\d+ price_quotes -- Detailed info about price_quotes
```

### Sample Queries

```sql
-- Check sample assets
SELECT * FROM assets;

-- Count records in each table
SELECT 'assets' as table_name, COUNT(*) FROM assets
UNION ALL
SELECT 'price_quotes', COUNT(*) FROM price_quotes
UNION ALL
SELECT 'price_history', COUNT(*) FROM price_history;
```

## Troubleshooting

### Reset Database Completely

```bash
# Stop and remove PostgreSQL container
docker-compose down postgres
docker volume rm compose_postgres_data

# Start fresh
docker-compose up -d postgres

# Run migrations again
migrate -path migrations -database "$DATABASE_URL" up
```

### Migration Error "Dirty Database"

```bash
# Force to specific version
migrate -path migrations -database "$DATABASE_URL" force 1

# Then retry
migrate -path migrations -database "$DATABASE_URL" up
```

## Environment Files

- `.env` - Local development (localhost connection)
- `.env.docker` - Docker environment (postgres hostname)

## Next Steps

After setting up the database:

1. Update Go service to connect to PostgreSQL
2. Implement database repository layer
3. Create API endpoints for finance operations
4. Set up background jobs for price updates
