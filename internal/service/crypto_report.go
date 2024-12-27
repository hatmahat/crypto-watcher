package service

import (
	"context"
	"crypto-watcher-backend/internal/config"
	"crypto-watcher-backend/internal/constant/asset_const"
	"crypto-watcher-backend/internal/constant/notification_const"
	"crypto-watcher-backend/internal/constant/user_preference_const"
	"crypto-watcher-backend/internal/entity"
	"crypto-watcher-backend/internal/entity/helper"
	"crypto-watcher-backend/internal/repository"
	"crypto-watcher-backend/pkg/format"
	"crypto-watcher-backend/pkg/telegram_bot_api"

	"github.com/sirupsen/logrus"
)

func (cs *cryptoService) dailyCoinPriceReport(ctx context.Context, assetCode string, coinPriceUSD float64, rateUSDToIDR int) {
	const funcName = "[internal][service]dailyCoinPriceReport"

	getUserFilter := repository.GetUserFilter{
		ReportTime:     format.GetSimpleTime(),
		AssetType:      asset_const.CRYPTO,
		AssetCode:      assetCode,
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

	usdPrice := format.ThousandSepartor(int64(coinPriceUSD), ',')
	idrPrice := format.ThousandSepartor(int64(coinPriceUSD*float64(rateUSDToIDR)), '.')

	if config.DebugMode {
		logrus.Infof("Asset Code (%s):\nUSD %s\nIDR %s\n", assetCode, usdPrice, idrPrice)
	}

	//var coinName string
	// coinNameMap, err := validation.ValidateFromMapper(assetCode, asset_const.AssetCodeNameMapper)
	// if err != nil {
	// 	logrus.Errorf("%s: Coin Name [%s] Not Found", funcName, assetCode)
	// }

	// if coinNameMap != nil {
	// 	coinName = *coinNameMap
	// }

	message := telegram_bot_api.CoinPriceAlertSuperSimple{
		CoinCode:          assetCode,
		CurrentPriceUSD:   usdPrice,
		CurrentPriceIDR:   idrPrice,
		FormattedDateTime: format.GetFormattedDateTimeWithDay(),
	}
	for _, user := range users {
		if user.TelegramChatId == nil {
			logrus.Errorf("%s: User ID [%d] telegram_chat_id is null", funcName, user.Id)
			continue
		}
		// TODO: Send telegram via jobs?
		go cs.sendTelegramMessage(ctx, user, &message)
	}
}

// sendTelegramMessage sends a message to a user via Telegram and logs the notification status in the database.
//
// Parameters:
//   - ctx: The context for controlling the request lifetime.
//   - user: A UserAndUserPreference struct containing the user's ID, preference ID, and Telegram chat ID.
//   - message: A Message struct representing the message to be sent via Telegram.
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
		Metadata: &entity.Metadata{
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
