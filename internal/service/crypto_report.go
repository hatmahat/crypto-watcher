package service

import (
	"context"
	"crypto-watcher-backend/internal/constant/asset_const"
	"crypto-watcher-backend/internal/constant/notification_const"
	"crypto-watcher-backend/internal/constant/user_preference_const"
	"crypto-watcher-backend/internal/entity"
	"crypto-watcher-backend/internal/entity/helper"
	"crypto-watcher-backend/internal/repository"
	"crypto-watcher-backend/pkg/format"
	"crypto-watcher-backend/pkg/telegram_bot_api"
	"fmt"

	"github.com/sirupsen/logrus"
)

func (cs *cryptoService) dailyBitcoinPriceReport(ctx context.Context, bitcoinPriceUSD, rateUSDToIDR *int) {
	const funcName = "[internal][service]DailyBitcoinPriceReport"

	// TODO (improvement): not only support bitcoin, get it from user preference
	getUserFilter := repository.GetUserFilter{
		ReportTime:     format.GetSimpleTime(),
		AssetType:      asset_const.CRYPTO,
		AssetCode:      asset_const.BTC,
		PreferenceType: user_preference_const.DailyReport,
	}
	users, err := cs.userRepo.GetUserAndUserPreferenceByReportTime(ctx, getUserFilter) // TODO (improvement): make it effective so it won't query every minute (chaching)
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

	message := telegram_bot_api.BitcoinPriceAlertSimple{
		CurrentPriceUSD:   usdPrice,
		CurrentPriceIDR:   idrPrice,
		FormattedDateTime: format.GetFormattedDateTimeWithDay(),
	}
	for _, user := range users {
		if user.TelegramChatId == nil {
			logrus.Errorf("%s: User ID [%d] telegram_chat_id is null", funcName, user.Id)
			continue
		}
		go cs.sendTelegramMessage(ctx, user, &message)
	}
}

func (cs *cryptoService) sendTelegramMessage(ctx context.Context, user helper.UserAndUserPreference, message telegram_bot_api.Message) {
	const funcName = "[internal][service]sendTelegramMessage"

	msg := message.Message()
	pvdr := notification_const.TelegramBotAPI

	err := cs.telegramBot.SendTelegramMessageByMessageId(*user.TelegramChatId, msg)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":     err.Error(),
			"message": msg,
			"chat_id": user.TelegramChatId,
		}).Errorf("%s: Error Sending Message via Telegram", funcName)
	}

	status := notification_const.SENT
	if err != nil {
		status = notification_const.FAILED
	}

	notif := entity.Notification{
		UserId:       user.Id,
		PreferenceId: user.PreferenceId,
		Status:       status,
		Parameters: entity.Parameters{
			Message:  &msg,
			Provider: &pvdr,
		},
	}
	err = cs.notifRepo.InsertNotification(ctx, notif)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":   err.Error(),
			"user":  user,
			"notif": notif,
		}).Errorf("%s: Error Inserting Notification to DB", funcName)
	}
}
