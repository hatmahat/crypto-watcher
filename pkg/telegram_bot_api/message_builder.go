package telegram_bot_api

import "fmt"

type (
	Message interface {
		Message() string
	}

	BitcoinPriceAlert struct {
		PercentageIncrease string
		CurrentPriceUSD    string
		CurrentPriceIDR    string
		PriceChangeUSD     string
		PriceChangeIDR     string
		FormattedDateTime  string
	}

	BitcoinPriceAlertSimple struct {
		CurrentPriceUSD   string
		CurrentPriceIDR   string
		FormattedDateTime string
	}
)

func (b *BitcoinPriceAlert) Message() string {
	return fmt.Sprintf(
		bitcoin_price_alert_template,
		b.PercentageIncrease,
		b.CurrentPriceUSD,
		b.CurrentPriceIDR,
		b.PriceChangeUSD,
		b.PriceChangeIDR,
		b.FormattedDateTime,
	)
}

func (b *BitcoinPriceAlertSimple) Message() string {
	return fmt.Sprintf(
		bitcoin_price_alert_simple_template,
		b.CurrentPriceUSD,
		b.CurrentPriceIDR,
		b.FormattedDateTime,
	)
}
