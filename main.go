package main

import (
	"TGbot/config"
	"TGbot/database"
	"TGbot/errors"
	"TGbot/weather"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

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

	db, err := database.Connect()
	errors.CheckError(err)

	defer db.Close()

	bot, err = tgbotapi.NewBotAPI(config.BOTKEY)
	errors.CheckError(err)

	updConfig.Timeout = 60
	updConfig.Limit = 1
	updConfig.Offset = 0

	updChannel = bot.GetUpdatesChan(updConfig)

	for {
		update = <-updChannel

		for {
			if update.Message != nil {
				row := db.QueryRow(config.UserDB, update.Message.Chat.ID)
				err = row.Scan(&User.user_id, &User.user_name, &User.user_tgid)
				if err != nil {
					_, err := db.Exec(config.AddNewUser, update.Message.Chat.ID, update.Message.From.FirstName, update.Message.From.UserName)
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

					data, _ := ioutil.ReadFile("images/6.png")
					msgPhoto := tgbotapi.FileBytes{Name: "images/6.png", Bytes: data}
					msgConfig := tgbotapi.NewPhoto(update.Message.Chat.ID, msgPhoto)

					weatherInfo, err := weatherTemperature(weather, update)
					errors.CheckError(err)

					msgConfig.Caption = weatherInfo
					bot.Send(msgConfig)

				} else if update.Message.Command() == "stop" && update.Message.From.ID == config.ROOT {
					bot.StopReceivingUpdates()
				}
			} else if update.Message.Text == "users" && update.Message.From.ID == config.ROOT {
				rows, err := db.Query(config.UsersFromDB)
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

				// get document with users
				path, err := OutputUsers(user)
				errors.CheckError(err)
				data, _ := ioutil.ReadFile(path)
				msgFile := tgbotapi.FileBytes{Name: "usersDatabaseInfo.txt", Bytes: data}
				msgConfig := tgbotapi.NewDocument(update.Message.Chat.ID, msgFile)
				bot.Send(msgConfig)
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

func OutputUsers(user []UserInfo) (string, error) {
	path := "usersDatabaseInfo.txt"
	file, err := os.Create(path)
	errors.CheckError(err)

	defer file.Close()

	for index, ui := range user {
		writingInFile(file, ui, index)
	}

	return path, nil
}

func writingInFile(file *os.File, user UserInfo, index int) {
	userInfo := fmt.Sprintf("[%d] Username: %v | User ID: %v\n",
		index, user.user_name, user.user_id)
	_, err := io.Copy(file, strings.NewReader(userInfo))
	errors.CheckError(err)
}
