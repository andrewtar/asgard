package telegram

import (
	"flag"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"asgard/common/log"
)

var telegramApiEndpointFlag = flag.String("telegram-api-endpoint", tgbotapi.APIEndpoint, "Telegram API endpoint")
var debugFlag = flag.Bool("telegram-bot-debug", false, "Pring debug info")

const (
	deprecatedTokenFlagEnv = "BOT_TOKEN"
	tokenFlagEnv           = "TELEGRAM_BOT_TOKEN"
)

func CreateBot() *tgbotapi.BotAPI {
	token := os.Getenv(deprecatedTokenFlagEnv)
	if token != "" {
		log.Logger().
			WithField("env", deprecatedTokenFlagEnv).
			Warn("Telegram bot token is found in deprecatd evironment variable")
	} else if token = os.Getenv(tokenFlagEnv); token != "" {
		log.Logger().
			WithField("env", tokenFlagEnv).
			Debug("Telegram bot token is found")
	} else {
		log.Logger().
			WithField("deprecated_env", deprecatedTokenFlagEnv).
			WithField("env", tokenFlagEnv).
			Panic("Telegram bot token is missed and exprected in one of the env")
	}

	if *telegramApiEndpointFlag == "" {
		log.Logger().Panic("Telegram API must not be emmpty")
	}

	bot, err := tgbotapi.NewBotAPIWithAPIEndpoint(token, *telegramApiEndpointFlag)
	if err != nil {
		log.Logger().
			WithError(err).
			Panic("Failed to init Telegram Bot client")
	}
	bot.Debug = *debugFlag
	log.Logger().
		WithField("name", bot.Self.UserName).
		Printf("Bot was logged")
	return bot
}

func GetCommand(update *tgbotapi.Update) string {
	isCommand := update.Message != nil &&
		len(update.Message.Entities) == 1 &&
		update.Message.Entities[0].Type == "bot_command"
	if !isCommand {
		return ""
	}
	return update.Message.Text[1:]
}
