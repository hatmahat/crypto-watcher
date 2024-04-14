-- Verify crypto-watcher:1713025849_init_crypto_watcher_schema on pg

BEGIN;

-- Select statement for users table
SELECT
    id,
    uuid,
    username,
    email,
    phone_number,
    created_at,
    updated_at
FROM
    users;

-- Select statement for currency_rates table
SELECT
    id,
    currency_pair,
    rate,
    created_at
FROM
    currency_rates;

-- Select statement for asset_prices table
SELECT
    id,
    asset_type,
    asset_code,
    price_usd,
    created_at
FROM
    asset_prices;

-- Select statement for user_preferences table
SELECT
    id,
    user_id,
    asset_type,
    asset_code,
    threshold_percentage,
    observation_period,
    created_at,
    updated_at
FROM
    user_preferences;

-- Select statement for notifications table
SELECT
    id,
    user_id,
    preference_id,
    message,
    created_at
FROM
    notifications;

ROLLBACK;
