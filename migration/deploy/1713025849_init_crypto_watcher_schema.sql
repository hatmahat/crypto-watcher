-- Deploy crypto-watcher:1713025849_init_crypto_watcher_schema to pg

BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    uuid UUID DEFAULT uuid_generate_v4(),
    username VARCHAR(255),
    email VARCHAR(255) UNIQUE,
    phone_number VARCHAR(50) UNIQUE,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE currency_rates (
    rate_id SERIAL PRIMARY KEY,
    currency_pair VARCHAR(7) DEFAULT 'USD_IDR',
    rate DECIMAL(18, 4),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE bitcoin_prices (
    price_id SERIAL PRIMARY KEY,
    price_usd DECIMAL(18, 4),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE user_preferences (
    preference_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    threshold_percentage DECIMAL(5, 2),
    observation_period INTEGER DEFAULT 15,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE notifications (
    notification_id SERIAL PRIMARY KEY,
    notification_type VARCHAR(255),
    user_id INTEGER REFERENCES users(user_id),
    price_id INTEGER REFERENCES bitcoin_prices(price_id),
    parameters JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_users_uuid ON users USING btree (uuid);
CREATE INDEX IF NOT EXISTS idx_users_username ON users USING btree (username);
CREATE INDEX IF NOT EXISTS idx_currency_rates_pair_time ON currency_rates USING btree (currency_pair, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_bitcoin_prices_created_at ON bitcoin_prices USING btree (created_at DESC);
CREATE INDEX IF NOT EXISTS idx_user_preferences_user_id ON user_preferences USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications USING btree (user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_price_id ON notifications USING btree (price_id);
CREATE INDEX IF NOT EXISTS idx_notifications_parameters ON notifications USING gin (parameters);

COMMIT;
