# Crypto Watcher - Backend  

Crypto Watcher is a dynamic monitoring & alerting system designed to track your cryptocurrency watch list in real-time. It provides detailed reports on cryptocurrency performance and sends timely notifications via Telegram to keep users updated on significant changes and trends.

## Features
- **Monitor Cryptocurrencies**: Continuously track the price and performance of various cryptocurrencies.
- **Report Cryptocurrencies**: Generate detailed reports on cryptocurrency metrics, providing insights into historical performance and current status.
- **Send Notifications via Telegram**: Leverage Telegram to send instant alerts about critical price changes based on user preferences.

## Architecture
![system design diagram](documentation/system-design.png)
[Cost Estimate Summary](https://cloud.google.com/products/calculator/estimate-preview/CiRhM2E1MzBjMi0wODI4LTQ1MjEtOTU1NC03Y2QzNjE2ZmRjOTUQAQ%3D%3D)

## Sequence Diagram
```mermaid
sequenceDiagram
    participant CRYPTO_WATCHER_DB
    participant CRYPTO_WATCHER
    participant CURRENCY_CONVERTER_API
    participant COINGECKO_API
    participant TELEGRAM_BOT_API

    Note over CRYPTO_WATCHER: runs every 1 minute

    CRYPTO_WATCHER->>CRYPTO_WATCHER_DB: Get distict asset code of crypto
    CRYPTO_WATCHER_DB-->>CRYPTO_WATCHER: Return asset codes

    CRYPTO_WATCHER->>CRYPTO_WATCHER_DB: Get today's date currency rate
    alt if current date currency rate exist
    CRYPTO_WATCHER_DB-->>CRYPTO_WATCHER: Return today's date currency rate
    else if today currency rate does not exist
    CRYPTO_WATCHER_DB-->>CRYPTO_WATCHER: Return no rows
    CRYPTO_WATCHER->>CURRENCY_CONVERTER_API: Get currency rate
    CURRENCY_CONVERTER_API-->>CRYPTO_WATCHER: Return currency rate
    CRYPTO_WATCHER->>CRYPTO_WATCHER_DB: Insert today's currency rate
    CRYPTO_WATCHER_DB-->>CRYPTO_WATCHER: Acknowledge insertion
    end

    CRYPTO_WATCHER->>COINGECKO_API: Get current coin prices based on asset codes
    COINGECKO_API-->>CRYPTO_WATCHER: Return current coin proces

    CRYPTO_WATCHER->>CRYPTO_WATCHER_DB: Insert current coin prices
    CRYPTO_WATCHER_DB-->>CRYPTO_WATCHER: Acknowledge insertion

    par for daily report
    CRYPTO_WATCHER->>CRYPTO_WATCHER_DB: Get users where report_time = time.Now().Format("15:04") and preference_type = 'daily_report'
    CRYPTO_WATCHER_DB-->>CRYPTO_WATCHER: Return user detail with current report time
    CRYPTO_WATCHER->>TELEGRAM_BOT_API: Send telegram message with chat id
        alt if response success
        TELEGRAM_BOT_API-->>CRYPTO_WATCHER: Success response
        CRYPTO_WATCHER->>CRYPTO_WATCHER_DB: Insert notification SENT
        else if response failed
        TELEGRAM_BOT_API-->>CRYPTO_WATCHER: Failed response
        CRYPTO_WATCHER->>CRYPTO_WATCHER_DB: Insert notification FAILED
        end
    end
```  

## Database ERD
![erd diagram](documentation/ERD.png)  
[DBDiagram](https://dbdiagram.io/d/661ab47403593b6b61e97fb8)