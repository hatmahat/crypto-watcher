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

	CoinPriceAlertSimple struct {
		CoinName          string
		CoinCode          string
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

func (b *CoinPriceAlertSimple) Message() string {
	return fmt.Sprintf(
		coin_price_alert_simple_template,
		b.CoinCode,
		b.CoinName,
		b.CoinCode,
		b.CurrentPriceUSD,
		b.CurrentPriceIDR,
		b.FormattedDateTime,
	)
}
