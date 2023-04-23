package main

import (
	"TGbot/config"
	"TGbot/database"
	"TGbot/errors"
	"TGbot/weather"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func main() {
	var (
		bot        *tgbotapi.BotAPI
		updChannel tgbotapi.UpdatesChannel
		update     tgbotapi.Update
		updConfig  tgbotapi.UpdateConfig
		users      []int64
		user_id    int64
	)

	deleteUser := fmt.Sprintf(`delete from users where user_id = $1`)
	addNewUser := fmt.Sprintf(`insert into "users"("user_id") values($1)`)
	userDB := fmt.Sprintf(`select * from users where user_id = $1`)

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

		if update.Message != nil {

			row := db.QueryRow(userDB, update.Message.Chat.ID)
			err = row.Scan(&user_id)
			if err != nil {
				_, err := db.Exec(addNewUser, update.Message.Chat.ID)
				errors.CheckError(err)

				if update.Message.IsCommand() {
					if update.Message.Command() == "weather" {
						weather, err := weather.Weather()
						errors.CheckError(err)

						fmt.Printf("%+v\n", weather)
					} else if update.Message.Command() == "delU" {
						_, err := db.Exec(deleteUser, update.Message.Chat.ID)
						errors.CheckError(err)

						// users, err = action.DeleteUser(users, update.Message.Chat.ID)
						// errors.CheckError(err)
					}
				} else if update.Message.Text == "users" {
					for _, user := range users {
						botMSG := fmt.Sprintf("user ID: %v", user)
						msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, botMSG)
						bot.Send(msgConfig)
					}
				}
			} else {
				_, err := db.Exec(deleteUser, update.Message.Chat.ID)
				errors.CheckError(err)
			}
		}
	}

	bot.StopReceivingUpdates()
}
