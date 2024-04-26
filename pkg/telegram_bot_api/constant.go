package telegram_bot_api

const (
	// Parse Modes
	HTML = "HTML"

	// Templates
	bitcoin_price_alert_template = `
	<b>🚨 Bitcoin Price Alert</b>
	
	Trigger: <b>Bitcoin</b> has <b>increased by %s%%</b> in the last 24 hours.
			
	<b>Current Price:</b> 
	- USD: <b>$%s</b>
	- IDR: <b>Rp%s</b>
			
	<b>Comparison:</b> 
	- The price in <b>USD</b> is <b>up $%s</b> from yesterday.
	- The price in <b>IDR</b> is <b>up Rp%s</b> from yesterday.

	<b>Date & Time:</b> %s
	`
)
