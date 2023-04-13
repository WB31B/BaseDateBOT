package main

import (
	"database/sql"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

const tgbotapiKey = ""

var rootMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("All Users"),
	),
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "postgres"
	password = ""
	dbName   = "tusergbot"
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
	CheckError(err)

	updConfig.Timeout = 60
	updConfig.Limit = 1
	updConfig.Offset = 0

	updChannel = bot.GetUpdatesChan(updConfig)

	psqlsconn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	db, err := sql.Open("postgres", psqlsconn)
	CheckError(err)

	defer db.Close()

	for {
		update = <-updChannel

		if update.Message != nil {
			if update.Message.IsCommand() {
				cmdText := update.Message.Command()
				if cmdText == "db" {
					insertStmt := fmt.Sprintf(`insert into "users"("user_id", "user_name", "user_message") values(%v, '%v', '%v')`,
						update.Message.From.ID, update.Message.From.FirstName, update.Message.Text)

					_, e := db.Exec(insertStmt)
					CheckError(e)

					msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, "hello")

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

					fmt.Println(msgInfoUser)

					msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, msgInfoUser)
					bot.Send(msgConfig)
				}
			}
		}

	}

	bot.StopReceivingUpdates()
}

func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
