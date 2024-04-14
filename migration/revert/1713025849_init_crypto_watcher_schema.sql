-- Revert crypto-watcher:1713025849_init_crypto_watcher_schema from pg

BEGIN;

-- Drop the indexes first
DROP INDEX IF EXISTS idx_notifications_user_preference_id;
DROP INDEX IF EXISTS idx_user_preferences_user_id;
DROP INDEX IF EXISTS idx_asset_prices_created_at;
DROP INDEX IF EXISTS idx_currency_rates_pair_time;
DROP INDEX IF EXISTS idx_users_username;
DROP INDEX IF EXISTS idx_users_uuid;

-- Drop the tables next. Notifications must be dropped before user_preferences, and user_preferences before users, due to the foreign key constraints.
DROP TABLE IF EXISTS notifications;
DROP TABLE IF EXISTS user_preferences;
DROP TABLE IF EXISTS asset_prices;
DROP TABLE IF EXISTS currency_rates;
DROP TABLE IF EXISTS users;

-- Drop the extension if it's no longer needed and not used elsewhere.
DROP EXTENSION IF EXISTS "uuid-ossp";

COMMIT;
