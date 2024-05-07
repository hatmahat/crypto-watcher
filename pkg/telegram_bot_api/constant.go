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

	<b>%s</b>
	`

	coin_price_alert_simple_template = `
	<b>🚨 %s (%s) Price Alert</b>
	
	Trigger: <b>%s</b> daily update
			
	<b>Current Price:</b> 
	- USD: <b>$%s</b>
	- IDR: <b>Rp%s</b>

	<b>%s</b>
	`

	coin_price_alert_super_simple_template = `
	<b>🚨 %s Daily Alert</b>
	<b>%s</b> = %s USD
	<b>%s</b> = %s IDR
	<b>%s</b>
	`
)
