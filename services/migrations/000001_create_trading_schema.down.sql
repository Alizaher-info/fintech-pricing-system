-- Rollback migration - Drop all tables in reverse order

DROP TABLE IF EXISTS price_alerts CASCADE;
DROP TABLE IF EXISTS watchlist CASCADE;
DROP TABLE IF EXISTS order_book CASCADE;
DROP TABLE IF EXISTS price_history CASCADE;
DROP TABLE IF EXISTS price_quotes CASCADE;
DROP TABLE IF EXISTS assets CASCADE;
