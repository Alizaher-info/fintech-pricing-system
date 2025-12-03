# Database Migration Update Summary

**Date:** December 3, 2025  
**Migration:** 000001_create_trading_schema.up.sql  
**Changes:** Added user order management capabilities

---

## What Changed

### 1. ORDER_BOOK Table - Updated ✅

**Added columns:**
- `user_id INTEGER NOT NULL` - Tracks which user placed the order (from MySQL)
- `quantity_filled DECIMAL(30, 8) DEFAULT 0` - Tracks how much of the order has been filled
- `status VARCHAR(20) DEFAULT 'open'` - Order status: open, filled, cancelled, partial
- `created_at TIMESTAMP DEFAULT NOW()` - When order was placed
- `updated_at TIMESTAMP DEFAULT NOW()` - Last modification time

**Added constraints:**
- `CHECK (status IN ('open', 'filled', 'cancelled', 'partial'))` - Valid status values

**Added indexes:**
- `idx_orderbook_user` on `(user_id, status)` - Fast user order queries
- `idx_orderbook_status` on `(status, created_at)` - Fast status filtering

**Before:**
```sql
CREATE TABLE order_book (
    id SERIAL PRIMARY KEY,
    asset_id INTEGER NOT NULL,
    side VARCHAR(10) NOT NULL,
    price DECIMAL(20, 8) NOT NULL,
    quantity DECIMAL(30, 8) NOT NULL,
    total DECIMAL(30, 8) NOT NULL,
    timestamp TIMESTAMP DEFAULT NOW()
);
```

**After:**
```sql
CREATE TABLE order_book (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,              -- NEW
    asset_id INTEGER NOT NULL,
    side VARCHAR(10) NOT NULL,
    price DECIMAL(20, 8) NOT NULL,
    quantity DECIMAL(30, 8) NOT NULL,
    quantity_filled DECIMAL(30, 8) DEFAULT 0,  -- NEW
    total DECIMAL(30, 8) NOT NULL,
    status VARCHAR(20) DEFAULT 'open',     -- NEW
    created_at TIMESTAMP DEFAULT NOW(),    -- NEW (renamed from timestamp)
    updated_at TIMESTAMP DEFAULT NOW()     -- NEW
);
```

---

### 2. WATCHLIST Table - Clarified ✅

**Added documentation:**
- Clarified that `user_id` references MySQL users table
- No foreign key constraint (managed at application layer)

---

### 3. PRICE_ALERTS Table - Clarified ✅

**Added documentation:**
- Clarified that `user_id` references MySQL users table
- No foreign key constraint (managed at application layer)

---

## Key Design Decisions

### No Users Table in PostgreSQL ✅

**Reason:** User data lives in MySQL (Symfony backend). PostgreSQL only stores `user_id` as an integer reference.

**Benefits:**
- ✅ No data duplication
- ✅ Single source of truth for user data (MySQL)
- ✅ Clear separation of concerns
- ✅ No sync issues between databases

### User ID Flow

```
┌─────────────┐
│   MySQL     │  User logs in, gets JWT with user_id
│  (Symfony)  │
└──────┬──────┘
       │
       │ JWT token contains: { user_id: 123, email: "..." }
       │
       ↓
┌─────────────┐
│   Symfony   │  Validates JWT, extracts user_id
│  Controller │
└──────┬──────┘
       │
       │ gRPC call: CreateOrder(user_id: 123, ...)
       │
       ↓
┌─────────────┐
│  Go Service │  Receives user_id from gRPC
└──────┬──────┘
       │
       │ INSERT INTO order_book (user_id, ...) VALUES (123, ...)
       │
       ↓
┌─────────────┐
│ PostgreSQL  │  Stores user_id as INTEGER (no FK)
└─────────────┘
```

---

## New Features Enabled

### 1. User Order Management
```sql
-- Get user's open orders
SELECT * FROM order_book 
WHERE user_id = 123 AND status = 'open';

-- Cancel user's order
UPDATE order_book 
SET status = 'cancelled', updated_at = NOW()
WHERE id = 456 AND user_id = 123;
```

### 2. Order History
```sql
-- Get user's order history
SELECT * FROM order_book 
WHERE user_id = 123 
ORDER BY created_at DESC 
LIMIT 50;
```

### 3. Partial Fills
```sql
-- Track partially filled orders
SELECT * FROM order_book 
WHERE user_id = 123 
  AND status = 'partial'
  AND quantity_filled < quantity;
```

### 4. Order Lifecycle Tracking
```sql
-- Order states:
-- 'open' → 'partial' → 'filled'
-- 'open' → 'cancelled'

UPDATE order_book
SET 
    quantity_filled = 1.5,
    status = CASE 
        WHEN 1.5 >= quantity THEN 'filled'
        ELSE 'partial'
    END,
    updated_at = NOW()
WHERE id = 789;
```

---

## Migration Instructions

### Apply Migration
```powershell
cd services
make migrate-up
```

### Verify Tables
```powershell
make db-shell

# In PostgreSQL shell:
\d order_book
SELECT * FROM order_book LIMIT 5;
```

### Rollback (if needed)
```powershell
make migrate-down
```

---

## Testing Queries

### Insert Test Order
```sql
INSERT INTO order_book (user_id, asset_id, side, price, quantity, total, status)
VALUES (123, 1, 'bid', 95000.00, 2.0, 190000.00, 'open');
```

### Query User Orders
```sql
SELECT 
    o.id,
    o.user_id,
    a.symbol,
    o.side,
    o.price,
    o.quantity,
    o.quantity_filled,
    o.status,
    o.created_at
FROM order_book o
JOIN assets a ON o.asset_id = a.id
WHERE o.user_id = 123
ORDER BY o.created_at DESC;
```

### Update Order Status
```sql
UPDATE order_book
SET 
    quantity_filled = 2.0,
    status = 'filled',
    updated_at = NOW()
WHERE id = 1 AND user_id = 123;
```

---

## What's Next

1. **Go Service Updates**
   - Implement CreateOrder gRPC endpoint
   - Implement GetUserOrders gRPC endpoint
   - Implement CancelOrder gRPC endpoint
   - Add order matching logic

2. **Symfony Integration**
   - Create OrderController to interface with Go service
   - Pass authenticated user_id to Go via gRPC
   - Return order data to frontend

3. **Frontend**
   - Create "My Orders" page
   - Add order placement form
   - Display order history
   - Real-time order status updates

---

**Migration Status:** ✅ Ready to Apply  
**Database Impact:** Schema change (adds columns to order_book)  
**Data Loss:** None (backward compatible with existing data)
