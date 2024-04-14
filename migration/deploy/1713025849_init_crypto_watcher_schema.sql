-- Deploy crypto-watcher:1713025849_init_crypto_watcher_schema to pg

BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4(),
    username VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    phone_number VARCHAR(50) UNIQUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE currency_rates (
    id SERIAL PRIMARY KEY,
    currency_pair VARCHAR(7) DEFAULT 'USD_IDR',
    rate DECIMAL(18, 4),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE asset_prices (
    id SERIAL PRIMARY KEY,
    asset_type VARCHAR(50), -- 'crypto' or 'stock'
    asset_code VARCHAR(50), -- such as 'BTC', 'ETH' for cryptos or 'AAPL', 'GOOGL' for stocks
    price_usd DECIMAL(18, 4),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_preferences (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    asset_type VARCHAR(50), -- 'crypto' or 'stock'
    asset_code VARCHAR(50), -- such as 'BTC', 'ETH' for cryptos or 'AAPL', 'GOOGL' for stocks
    threshold_percentage DECIMAL(5, 2),
    observation_period INTEGER DEFAULT 15, -- This is in minutes
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE notifications (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    preference_id INTEGER REFERENCES user_preferences(id), -- Reference to user preferences
    message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_uuid ON users USING btree (uuid);
CREATE INDEX IF NOT EXISTS idx_users_username ON users USING btree (username);
CREATE INDEX IF NOT EXISTS idx_currency_rates_pair_time ON currency_rates USING btree (currency_pair, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_asset_prices_created_at ON asset_prices USING btree (asset_type, asset_code, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_preference_id ON notifications USING btree (preference_id);

COMMIT;
