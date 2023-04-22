package bot

import (
	"TGbot/config"
	"TGbot/database"
	"TGbot/database/action"
	"TGbot/errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func StartBot() {
	var (
		bot        *tgbotapi.BotAPI
		updChannel tgbotapi.UpdatesChannel
		update     tgbotapi.Update
		updConfig  tgbotapi.UpdateConfig
		users      []int64
	)

	db, err := database.Connect()
	errors.CheckError(err)

	defer db.Close()

	botKey, err := config.GetKey("")
	errors.CheckError(err)

	bot, err = tgbotapi.NewBotAPI(botKey)
	errors.CheckError(err)

	updConfig.Timeout = 60
	updConfig.Limit = 1
	updConfig.Offset = 0

	updChannel = bot.GetUpdatesChan(updConfig)

	for {
		update = <-updChannel

		// var user_id int64
		deleteUser := fmt.Sprintf(`delete from users where user_id = $1`)
		addNewUser := fmt.Sprintf(`insert into "users"("user_id") values(%v)`, update.Message.From.ID)
		// row := db.QueryRow("select * from users where user_id = $1", update.Message.Chat.ID)

		// err = row.Scan(&user_id)
		// errors.CheckError(err)

		// if user_id == update.Message.Chat.ID {
		// 	_, err := db.Exec(deleteUser, update.Message.Chat.ID)
		// 	errors.CheckError(err)
		// } else {
		// 	_, err := db.Exec(addNewUser)
		// 	errors.CheckError(err)
		// }

		if update.Message != nil {

			user, err := GerUsers(users)
			errors.CheckError(err)

			if user == update.Message.Chat.ID {
				if update.Message.IsCommand() {
					if update.Message.Command() == "newuser" {
						botMsg := fmt.Sprintf("Hi: %v", user)
						msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, botMsg)
						bot.Send(msgConfig)
					} else if update.Message.Command() == "delU" {
						_, err := db.Exec(deleteUser, update.Message.Chat.ID)
						errors.CheckError(err)

						users, err = action.DeleteUser(users, update.Message.Chat.ID)
						errors.CheckError(err)
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
				errors.CheckError(err)

				_, err := db.Exec(addNewUser)
				errors.CheckError(err)
			}
		}
	}

	bot.StopReceivingUpdates()
}

func GerUsers(users []int64) (int64, error) {
	for _, user := range users {
		return user, nil
	}

	return 0, nil
}
