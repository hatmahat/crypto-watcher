package service

import (
	"context"
	"crypto-watcher-backend/internal/constant/asset_const"
	"crypto-watcher-backend/internal/repository"
	"crypto-watcher-backend/pkg/format"
	"crypto-watcher-backend/pkg/telegram_bot_api"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (cs *cryptoService) dailyBitcoinPriceReport(ctx context.Context, bitcoinPriceUSD, rateUSDToIDR *int) {
	const funcName = "[internal][service]DailyBitcoinPriceReport"

	// TODO (improvement): not only support bitcoin, get from user preference
	getUserFilter := repository.GetUserFilter{
		ReportTime: format.GetSimpleTime(),
		AssetType:  asset_const.CRYPTO,
		AssetCode:  asset_const.BTC,
	}
	users, err := cs.userRepo.GetUserByReportTime(ctx, getUserFilter) // TODO (improvement): make it effective so it won't query every minute (chaching)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":         err.Error(),
			"user_filter": getUserFilter,
		}).Errorf("%s: Error Getting User from Report Time", funcName)
		return
	}

	usdPrice := format.ThousandSepartor(int64(*bitcoinPriceUSD), ',')
	idrPrice := format.ThousandSepartor(int64(*bitcoinPriceUSD*(*rateUSDToIDR)), '.')
	fmt.Printf("USD %s\nIDR %s\n", usdPrice, idrPrice)

	message := telegram_bot_api.BitcoinPriceAlert{
		PercentageIncrease: "3.5",
		CurrentPriceUSD:    usdPrice,
		CurrentPriceIDR:    idrPrice,
		PriceChangeUSD:     "1,400",
		PriceChangeIDR:     "20,000,000",
		FormattedDateTime:  format.GetFormattedDateTimeWithDay(),
	}
	for _, user := range users {
		if user.TelegramChatId == nil {
			logrus.Errorf("%s: User ID [%d] telegram_chat_id is null", funcName, user.Id)
			continue
		}
		go cs.sendTelegramMessage(*user.TelegramChatId, &message)
	}
}

func (cs *cryptoService) sendTelegramMessage(chatId int64, message telegram_bot_api.Message) {
	const funcName = "[internal][service]sendTelegramMessage"
	err := cs.telegramBot.SendTelegramMessageByMessageId(chatId, message.Message())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":     err.Error(),
			"message": message.Message(),
			"chat_id": chatId,
		}).Errorf("%s: Error Sending Message via Telegram", funcName)
	}
	// TODO: insert to notifications table
}
