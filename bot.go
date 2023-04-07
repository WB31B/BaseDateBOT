package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const tgbotapiKey = ""

var mainMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("1"),
		tgbotapi.NewKeyboardButton("2"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/root"),
	),
)

func main() {
	var (
		bot        *tgbotapi.BotAPI
		err        error
		updChannel tgbotapi.UpdatesChannel
		update     tgbotapi.Update
		updConfig  tgbotapi.UpdateConfig
	)
	bot, err = tgbotapi.NewBotAPI(tgbotapiKey)
	if err != nil {
		panic(err.Error())
	}

	updConfig.Timeout = 60
	updConfig.Limit = 1
	updConfig.Offset = 0

	updChannel = bot.GetUpdatesChan(updConfig)

	for {
		update = <-updChannel

		if update.Message != nil {
			if update.Message.IsCommand() {
				cmdText := update.Message.Command()
				if cmdText == "root" {
					msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi -> "+update.Message.From.FirstName)
					bot.Send(msgConfig)
				} else if cmdText == "menu" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Main Menu")
					msg.ReplyMarkup = mainMenu
					bot.Send(msg)
				}
			} else {
				fmt.Printf("[Message]: %s | [Name]: %s\n", update.Message.Text, update.Message.From.FirstName)

				msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				bot.Send(msgConfig)
			}
		}
	}

	bot.StopReceivingUpdates()
}
