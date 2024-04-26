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
