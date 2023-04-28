package main

import (
	"TGbot/config"
	"TGbot/database"
	"TGbot/errors"
	"TGbot/weather"
	"fmt"
	"io/ioutil"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

var User struct {
	user_id   int64
	user_name string
	user_tgid string
}

type UserInfo struct {
	user_id   int64
	user_name string
	user_tgid string
}

func main() {
	var (
		bot        *tgbotapi.BotAPI
		updChannel tgbotapi.UpdatesChannel
		update     tgbotapi.Update
		updConfig  tgbotapi.UpdateConfig
		user       []UserInfo
	)

	// deleteUser := fmt.Sprintf(`delete from users where user_id = $1`)
	addNewUser := fmt.Sprintf(`insert into "users"("user_id", "user_name", "user_tgid") values($1, $2, $3)`)
	userDB := fmt.Sprintf(`select * from users where user_id = $1`)
	usersDB := fmt.Sprintf(`select * from users`)

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
				err = row.Scan(&User.user_id, &User.user_name, &User.user_tgid)
				if err != nil {
					_, err := db.Exec(addNewUser, update.Message.Chat.ID, update.Message.From.FirstName, update.Message.From.UserName)
					errors.CheckError(err)
					break
				}
				break
			}
		}

		if update.Message != nil {
			if update.Message.IsCommand() {
				if update.Message.Command() == "weather" {
					weather, err := weather.Weather("kyiv")
					errors.CheckError(err)

					data, _ := ioutil.ReadFile("images/6.png")
					msgPhoto := tgbotapi.FileBytes{Name: "images/6.png", Bytes: data}
					msgConfig := tgbotapi.NewPhoto(update.Message.Chat.ID, msgPhoto)

					weatherInfo, err := weatherTemperature(weather, update)
					errors.CheckError(err)

					msgConfig.Caption = weatherInfo
					bot.Send(msgConfig)
				} else if update.Message.Command() == "stop" && update.Message.From.ID == 673324657 {
					bot.StopReceivingUpdates()
				}
			} else if update.Message.Text == "users" && update.Message.From.ID == 673324657 {
				rows, err := db.Query(usersDB)
				errors.CheckError(err)

				defer rows.Close()

				for rows.Next() {
					ui := UserInfo{}
					err := rows.Scan(&ui.user_id, &ui.user_name, &ui.user_tgid)
					if err != nil {
						fmt.Println(err)
						continue
					}
					user = append(user, ui)
				}

				outputUsers(user)
			}
		}
	}
}

func weatherTemperature(weather *weather.WeatherData, update tgbotapi.Update) (string, error) {
	if weather.Data.Values.Temperature < 8 {
		weatherInfo := fmt.Sprintf("👨‍💻 User ID: [%v]\n🌍 Country: %v\n🥶 Temperature: %v\n💧 Humidity: %v\n☁️ Cloud Cover: %v\n💨 Visibility: %v\n\n⏰ Time: %v\n",
			update.Message.From.ID,
			weather.Location.Name,
			weather.Data.Values.Temperature,
			weather.Data.Values.Humidity,
			weather.Data.Values.CloudCover,
			weather.Data.Values.Visibility,
			weather.Data.Time)
		return weatherInfo, nil
	} else {
		weatherInfo := fmt.Sprintf("👨‍💻 User ID: [%v]\n🌍 Country: %v\n🥵 Temperature: %v\n💧 Humidity: %v\n☁️ Cloud Cover: %v\n💨 Visibility: %v\n\n⏰ Time: %v\n",
			update.Message.From.ID,
			weather.Location.Name,
			weather.Data.Values.Temperature,
			weather.Data.Values.Humidity,
			weather.Data.Values.CloudCover,
			weather.Data.Values.Visibility,
			weather.Data.Time)
		return weatherInfo, nil
	}
}

func outputUsers(user []UserInfo) {
	for _, ui := range user {
		fmt.Println(ui.user_name)
	}
}
