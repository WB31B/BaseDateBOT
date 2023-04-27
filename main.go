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

var User struct {
	user_id   int64
	user_name string
}

func main() {
	var (
		bot        *tgbotapi.BotAPI
		updChannel tgbotapi.UpdatesChannel
		update     tgbotapi.Update
		updConfig  tgbotapi.UpdateConfig
		users      []int64
	)

	// deleteUser := fmt.Sprintf(`delete from users where user_id = $1`)
	addNewUser := fmt.Sprintf(`insert into "users"("user_id", "user_name") values($1, $2)`)
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

		for {
			if update.Message != nil {
				row := db.QueryRow(userDB, update.Message.Chat.ID)
				err = row.Scan(&User.user_id, &User.user_name)
				if err != nil {
					_, err := db.Exec(addNewUser, update.Message.Chat.ID, update.Message.From.FirstName)
					errors.CheckError(err)
					break
				}
				break
			}
		}

		if update.Message != nil {
			if update.Message.IsCommand() {
				if update.Message.Command() == "weather" {
					weather, err := weather.Weather("london")
					errors.CheckError(err)

					weatherInfo := fmt.Sprintf("User ID: [%v]\nCountry: %v\nTemperature: %v\nHumidity: %v\nCloud Cover: %v\nVisibility: %v\n\nTime: %v\n",
						update.Message.From.ID,
						weather.Location.Name,
						weather.Data.Values.Temperature,
						weather.Data.Values.Humidity,
						weather.Data.Values.CloudCover,
						weather.Data.Values.Visibility,
						weather.Data.Time)

					msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, weatherInfo)
					bot.Send(msgConfig)

				}
			} else if update.Message.Text == "users" {
				for _, user := range users {
					botMSG := fmt.Sprintf("user ID: %v", user)
					msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, botMSG)
					bot.Send(msgConfig)
				}
			}
		}
	}

	bot.StopReceivingUpdates()
}
