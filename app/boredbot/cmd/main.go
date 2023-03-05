package main

import (
	"encoding/json"
	"flag"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"asgard/app/boredbot/cmd"
	"asgard/common/api/bored"
	"asgard/common/api/telegram"
	"asgard/common/api/yc"
	debug "asgard/common/debug/service"
	"asgard/common/log"
)

var debugFlag = flag.Bool("debug", false, "Print debug logs")

type Command interface {
	GetDescription() string
	GetId() string
	Handle() (string, error)
}

func main() {
	flag.Parse()
	debug.Init()

	bot := telegram.CreateBot()

	commandsTable := registerCommands(bot, []Command{
		cmd.GetInformationCommand{},
		cmd.HaveFunCommand{
			BoredClient:       bored.NewBoredClient(*debugFlag),
			TransaltionClient: yc.NewClient(),
		},
	})

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if *debugFlag {
			jsonUpdate, err := json.Marshal(update)
			if err != nil {
				log.Logger().
					WithError(err).
					Error("Failed to serialize incomming update")
				continue
			}
			log.Logger().
				WithField("content", string(jsonUpdate)).
				Debug("Incomming update")
		}

		if update.Message == nil {
			log.Logger().
				Debug("Update without message was skipped")
			continue
		}

		log.Logger().
			WithField("user_id", update.Message.From.ID).
			WithField("user_name", update.Message.From.FirstName).
			Info("Incomming message")

		incommingCommand := telegram.GetCommand(&update)
		if incommingCommand == "" {
			log.Logger().
				Debug("Update without message was skipped, because it`s not a command")
			continue
		}
		log.Logger().
			WithField("command", incommingCommand).
			Info("Incomming command")

		command, known := commandsTable[incommingCommand]
		if !known {
			log.Logger().
				WithField("command", incommingCommand).
				Error("Unknown command")
			continue
		}

		result, err := command.Handle()
		if err != nil {
			log.Logger().
				WithError(err).
				WithField("command", command.GetId()).
				Error("Failed to handle command")
			continue
		}

		response := tgbotapi.NewMessage(update.Message.Chat.ID, result)
		response.ReplyToMessageID = update.Message.MessageID
		if _, err := bot.Send(response); err != nil {
			log.Logger().
				WithError(err).
				Error("Failed to send response")
			continue
		}
		log.Logger().Info("Response was sent")
	}
}

func registerCommands(bot *tgbotapi.BotAPI, commands []Command) map[string]Command {
	commandsTable := map[string]Command{}
	for _, command := range commands {
		_, alreadyExists := commandsTable[command.GetId()]
		if alreadyExists {
			log.Logger().
				WithField("command", command.GetId()).
				Panic("Command already defined")
		}
		commandsTable[command.GetId()] = command
	}

	commandRegisterRequest := []tgbotapi.BotCommand{}
	for id, command := range commandsTable {
		commandRegisterRequest = append(commandRegisterRequest, tgbotapi.BotCommand{
			Command:     id,
			Description: command.GetDescription(),
		})
	}
	if len(commandRegisterRequest) == 0 {
		log.Logger().Panic("Can`t start without commands")
	}

	if _, err := bot.Request(tgbotapi.NewSetMyCommands(commandRegisterRequest...)); err != nil {
		log.Logger().
			WithError(err).
			Panic("Failed to register commands in Telegram API")
	}
	return commandsTable
}
