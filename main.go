package main

import (
	"TGbot/config"
	"TGbot/database"
	"TGbot/errors"
	"TGbot/log"
	"TGbot/weather"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"

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

const weatherTitle = "🌏 [WEATHER INFORMATION] 🌕"

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

	start_time := time.Now()

	updChannel = bot.GetUpdatesChan(updConfig)

	for {
		update = <-updChannel

		command := update.Message.Command()

		row := db.QueryRow(config.UserDB, update.Message.Chat.ID)
		err = row.Scan(&User.user_id, &User.user_name, &User.user_tgid, &start_time)
		if err != nil {
			fmt.Println("1")
			log.StartBot(update.Message.From.ID)
			_, err := db.Exec(config.AddNewUser, update.Message.Chat.ID, update.Message.From.FirstName, update.Message.From.UserName, start_time.Format("15:04:05"))
			errors.CheckError(err)

			reply := fmt.Sprintf("Hello, [%v], the developer of this bot is @WB31B The bot was created to display the weather of the region you specified. Write the city and the Bot will tell you the weather", update.Message.From.FirstName)
			msgConfig := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			bot.Send(msgConfig)
			continue
		}

		if command == "stop" && update.Message.From.ID == config.ROOTUSER {
			msgConfig := tgbotapi.NewMessage(update.Message.From.ID, "Bot stoped!")
			bot.Send(msgConfig)
			log.StopBotCommand(update.Message.From.ID)
			bot.StopReceivingUpdates()
			break
		} else if command == "users" && update.Message.From.ID == config.ROOTUSER {
			log.OutputUsersCommand(config.ROOTUSER)
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
		} else if command == "" {
			log.ShowWeather(update.Message.From.ID, update.Message.Text)
			weather, err := weather.Weather(update.Message.Text)
			errors.CheckError(err)

			weatherInfo, err := weatherTemperature(weather, update)
			errors.CheckError(err)

			msgConfig := tgbotapi.NewMessage(update.Message.From.ID, weatherInfo)

			_, err = db.Exec(config.AddNewMessage, update.Message.Text, start_time.Format("15:04:05"), update.Message.From.ID)
			errors.CheckError(err)
			bot.Send(msgConfig)
		} else {
			if command == "start" {
				log.StartCommand(update.Message.From.ID)
				msgConfig := tgbotapi.NewMessage(update.Message.From.ID,
					"Hello, the developer of this bot is @WB31B The bot was created to display the weather of the region you specified. Write the city and the Bot will tell you the weather")
				bot.Send(msgConfig)
			} else {
				log.IncorrectCommand(update.Message.From.ID)
				msgConfig := tgbotapi.NewMessage(update.Message.From.ID, "This command is INCORRECT!")
				bot.Send(msgConfig)
			}

		}
	}

}

func weatherTemperature(weather *weather.WeatherData, update tgbotapi.Update) (string, error) {
	if weather.Data.Values.Temperature < 15 {
		weatherInfo := fmt.Sprintf("%s\n\n👨‍💻 User ID: [%v]\n🌍 Country: %v\n🥶 Temperature: %v\n💧 Humidity: %v\n☁️ Cloud Cover: %v\n💨 Visibility: %v\n\n⏰ Time: %v\n",
			weatherTitle,
			update.Message.From.ID,
			weather.Location.Name,
			weather.Data.Values.Temperature,
			weather.Data.Values.Humidity,
			weather.Data.Values.CloudCover,
			weather.Data.Values.Visibility,
			weather.Data.Time)
		return weatherInfo, nil
	} else {
		weatherInfo := fmt.Sprintf("%s\n\n👨‍💻 User ID: [%v]\n🌍 Country: %v\n🥵 Temperature: %v\n💧 Humidity: %v\n☁️ Cloud Cover: %v\n💨 Visibility: %v\n\n⏰ Time: %v\n",
			weatherTitle,
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
