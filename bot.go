package main

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const tgbotapiKey = ""

var rootMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("All Users"),
	),
)

var DB *gorm.DB

type User struct {
	gorm.Model
	ID      int    `json:"id" gorm:"primary_key"`
	Name    string `json:"name"`
	Message string `json:"user_message"`
}

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

	db, err := gorm.Open(postgres.Open("postgres://postgres:pG2r4hack@localhost:5432/postgres"), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&User{})

	DB = db

	for {
		update = <-updChannel

		if update.Message != nil {
			if update.Message.IsCommand() {
				cmdText := update.Message.Command()
				if cmdText == "start" {
					msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "Hi -> "+update.Message.From.FirstName)

					bot.Send(msgConfig)
				}
				// } else if cmdText == "menu" {
				// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Main Menu")
				// 	msg.ReplyMarkup = mainMenu
				// 	bot.Send(msg)
				// }
			} else {
				if update.Message.Text == "dRootfaceT1" {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
					msg.ReplyMarkup = rootMenu
					bot.Send(msg)
				} else {
					msgInfoUser := fmt.Sprintf("[Your name]: %s\n[Your ID]: %v\n[Your message]: %s\n",
						update.Message.From.FirstName,
						update.Message.From.ID,
						update.Message.Text,
					)

					msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, msgInfoUser)
					bot.Send(msgConfig)
				}
			}
		}
	}

	bot.StopReceivingUpdates()
}
