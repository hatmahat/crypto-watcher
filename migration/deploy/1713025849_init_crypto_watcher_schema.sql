-- Deploy crypto-watcher:1713025849_init_crypto_watcher_schema to pg

BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uuid UUID NOT NULL DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    username VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone_number VARCHAR(50) UNIQUE,
    telegram_chat_id BIGINT DEFAULT NULL
);

CREATE TABLE currency_rates (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    currency_pair VARCHAR(7) NOT NULL DEFAULT 'USD_IDR' ,
    rate DECIMAL(18, 4) NOT NULL
);

CREATE TABLE asset_prices (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    asset_type VARCHAR(50) NOT NULL, -- 'CRYPTO' or 'STOCK'
    asset_code VARCHAR(50) NOT NULL, -- such as 'BTC', 'ETH' for cryptos or 'AAPL', 'GOOGL' for stocks
    price_usd DECIMAL(18, 4) NOT NULL
);

CREATE TABLE user_preferences (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    preference_type VARCHAR(50) NOT NULL, -- 'daily_report', 'weekly_report', 'price_alert', etc
    operator VARCHAR(20), -- '>,<', '+,-'. '+%,-%'
    asset_type VARCHAR(50) NOT NULL, -- 'CRYPTO' or 'STOCK'
    asset_code VARCHAR(50) NOT NULL, -- such as 'BTC', 'ETH' for cryptos or 'AAPL', 'GOOGL' for stocks
    price_checkpoint DECIMAL(18, 4),
    threshold_percentage DECIMAL(5, 2),
    observation_period INTEGER, -- This is in minutes
    report_time TIME WITHOUT TIME ZONE DEFAULT NULL,
    is_active boolean NOT NULL DEFAULT false
);

CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id INTEGER REFERENCES users(id) NOT NULL,
    preference_id INTEGER REFERENCES user_preferences(id) NOT NULL, -- Reference to user preferences
    status VARCHAR(50) NOT NULL, -- Notification status
    metadata JSONB DEFAULT '{}'
);

CREATE INDEX IF NOT EXISTS idx_users_uuid ON users USING btree (uuid);
CREATE INDEX IF NOT EXISTS idx_users_username ON users USING btree (username);
CREATE INDEX IF NOT EXISTS idx_currency_rates_pair_time ON currency_rates USING btree (currency_pair, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_asset_prices_created_at ON asset_prices USING btree (asset_type, asset_code, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_preference_id ON notifications USING btree (preference_id);

COMMIT;
