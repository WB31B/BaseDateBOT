package main

import (
	"TGbot/config"
	"TGbot/database"
	"TGbot/database/action"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

var rootMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("All Users"),
	),
)

func mainn() {
	var (
		bot        *tgbotapi.BotAPI
		err        error
		updChannel tgbotapi.UpdatesChannel
		update     tgbotapi.Update
		updConfig  tgbotapi.UpdateConfig
	)

	tgbotapiKey, err := config.GetKey()
	CheckError(err)

	bot, err = tgbotapi.NewBotAPI(tgbotapiKey)
	CheckError(err)

	updConfig.Timeout = 60
	updConfig.Limit = 1
	updConfig.Offset = 0

	updChannel = bot.GetUpdatesChan(updConfig)

	var users []int64

	db, err := database.Connect()
	CheckError(err)

	defer db.Close()

	// rows, err := db.Query(`SELECT user_id FROM users;`)
	// CheckError(err)

	for {
		update = <-updChannel

		if update.Message != nil {
			addNewUser := fmt.Sprintf(`insert into "users"("user_id") values(%v)`, update.Message.From.ID)
			deleteUser := fmt.Sprint(`delete from users where user_id = $1`)

			// for rows.Next() {
			// 	var user_id int64

			// 	err := rows.Scan(&user_id)
			// 	CheckError(err)

			// 	users = append(users, user_id)
			// }

			user, err := GerUsers(users)
			CheckError(err)

			if user == update.Message.Chat.ID {
				if update.Message.IsCommand() {
					if update.Message.Command() == "newuser" {
						botMsg := fmt.Sprintf("Hi: %v", user)
						msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, botMsg)
						bot.Send(msgConfig)
					} else if update.Message.Command() == "delU" {
						_, err := db.Exec(deleteUser, update.Message.Chat.ID)
						CheckError(err)

						users, err = action.DeleteUser(users, update.Message.Chat.ID)
						CheckError(err)
					}
				} else if update.Message.Text == "users" {
					for _, user := range users {
						botMSG := fmt.Sprintf("user ID: %v", user)
						msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, botMSG)
						bot.Send(msgConfig)
					}
				}
			} else {
				botMSG := fmt.Sprintf("HI New user -> %v", update.Message.Chat.ID)
				msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, botMSG)
				bot.Send(msgConfig)

				users, err = action.AddUser(users, update.Message.Chat.ID)
				CheckError(err)

				_, e := db.Exec(addNewUser)
				CheckError(e)
			}
		}

		// if update.Message.Text == "users" {
		// 	for _, user := range users {
		// 		botMSG := fmt.Sprintf("user ID: %v", user)
		// 		msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, botMSG)
		// 		bot.Send(msgConfig)
		// 	}
		// }
		// defer rows.Close()
	}

	// } else if cmdText == "menu" {
	// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Main Menu")
	// 	msg.ReplyMarkup = mainMenu
	// 	bot.Send(msg)
	// }

	bot.StopReceivingUpdates()
}

func CheckError(err error) {
	if err != nil {
		panic(err.Error())
	}
}

func GerUsers(users []int64) (int64, error) {
	for _, user := range users {
		return user, nil
	}

	return 0, nil
}
