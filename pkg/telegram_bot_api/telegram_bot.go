package telegram_bot_api

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

type (
	TelegramBot interface {
		SendTelegramMessageByMessageId(messageId int64, message string) error
	}

	telegramBot struct {
		bot *tgbotapi.BotAPI
	}
)

func NewTelegramBot(apiKey string) TelegramBot {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		errMsg := "Failed to create new BotAPI instance: %s"
		logrus.Errorf(errMsg, err)
		panic(fmt.Sprintf(errMsg, err))
	}

	return &telegramBot{
		bot: bot,
	}
}

func (t *telegramBot) SendTelegramMessageByMessageId(chatId int64, message string) error {
	const funcName = "[pkg][telegram_bot_api]SendTelegramMessageByMessageId"

	msgConfig := tgbotapi.NewMessage(chatId, message)
	msgConfig.ParseMode = HTML

	_, err := t.bot.Send(msgConfig)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"err":            err.Error(),
			"message_config": msgConfig,
		}).Errorf("%s: Error Send Message via Telegram", funcName)
		return err
	}

	return nil
}
