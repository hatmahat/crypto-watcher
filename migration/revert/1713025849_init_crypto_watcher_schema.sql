-- Revert crypto-watcher:1713025849_init_crypto_watcher_schema from pg

BEGIN;

DROP INDEX IF EXISTS idx_notifications_parameters;
DROP INDEX IF EXISTS idx_notifications_price_id;
DROP INDEX IF EXISTS idx_notifications_user_id;
DROP INDEX IF EXISTS idx_user_preferences_user_id;
DROP INDEX IF EXISTS idx_bitcoin_prices_created_at;
DROP INDEX IF EXISTS idx_currency_rates_pair_time;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_uuid;

DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS user_preferences;
DROP TABLE IF EXISTS bitcoin_prices;
DROP TABLE IF EXISTS currency_rates;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "uuid-ossp";

COMMIT;
