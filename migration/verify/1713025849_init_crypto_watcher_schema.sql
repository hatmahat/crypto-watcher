-- Verify crypto-watcher:1713025849_init_crypto_watcher_schema on pg

BEGIN;

SELECT 
    user_id, 
    uuid, 
    username, 
    email, 
    phone_number, 
    created_at, 
    updated_at 
FROM 
    users;

SELECT 
    rate_id, 
    currency_pair, 
    rate, 
    created_at 
FROM 
    currency_rates;

SELECT 
    price_id, 
    price_usd, 
    created_at 
FROM 
    bitcoin_prices;

SELECT 
    preference_id, 
    user_id, 
    threshold_percentage, 
    observation_period, 
    created_at, 
    updated_at 
FROM 
    user_preferences;

SELECT 
    notification_id, 
    notification_type, 
    user_id, 
    price_id, 
    parameters, 
    created_at 
FROM 
    notifications;


ROLLBACK;
