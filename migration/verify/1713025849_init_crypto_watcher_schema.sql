-- Verify crypto-watcher:1713025849_init_crypto_watcher_schema on pg

BEGIN;

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

SELECT
    id,
    currency_pair,
    rate,
    created_at
FROM
    currency_rates;

SELECT
    id,
    asset_type,
    asset_code,
    price_usd,
    created_at
FROM
    asset_prices;

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

SELECT
    id,
    user_id,
    preference_id,
    message,
    created_at
FROM
    notifications;

ROLLBACK;
