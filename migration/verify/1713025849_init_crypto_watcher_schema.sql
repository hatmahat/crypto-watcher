-- Verify crypto-watcher:1713025849_init_crypto_watcher_schema on pg

BEGIN;

SELECT
    id,
    uuid,
    created_at,
    updated_at,
    username,
    email,
    phone_number,
    telegram_chat_id
FROM
    users;

SELECT
    id,
    created_at,
    currency_pair,
    rate
FROM
    currency_rates;

SELECT
    id,
    created_at,
    asset_type,
    asset_code,
    price_usd
FROM
    asset_prices;

SELECT
    id,
    created_at,
    updated_at,
    user_id,
    preference_type,
    asset_type,
    asset_code,
    threshold_percentage,
    observation_period,
    report_time,
    is_active
FROM
    user_preferences;

SELECT
    id,
    created_at,
    user_id,
    preference_id,
    status,
    parameters
FROM
    notifications;

ROLLBACK;
