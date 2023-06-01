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

const weatherTitle = "🌏 [WEATHER INFORMATION] 🌕"

func main() {
	db, err := database.Connect()
	errors.CheckError(err)

	defer db.Close()

	config.Bot, err = tgbotapi.NewBotAPI(config.BOTKEY)
	errors.CheckError(err)

	config.UpdConfig.Timeout = 60
	config.UpdConfig.Limit = 1
	config.UpdConfig.Offset = 0

	config.UpdChannel = config.Bot.GetUpdatesChan(config.UpdConfig)

	for {
		timeNow := time.Now()
		start_time := fmt.Sprintf("%d-%02d-%02dT%02d:%02d",
			timeNow.Year(), timeNow.Month(), timeNow.Day(),
			timeNow.Hour(), timeNow.Minute())

		config.Update = <-config.UpdChannel

		command := config.Update.Message.Command()

		row := db.QueryRow(config.USERDB, config.Update.Message.Chat.ID)
		err = row.Scan(&config.User_id, &config.User_name, &config.User_tgid, &start_time)
		if err != nil {
			fmt.Println("BOT START")
			log.StartBot(config.Update.Message.From.ID)
			_, err := db.Exec(config.ADDNEWUSER,
				config.Update.Message.Chat.ID,
				config.Update.Message.From.FirstName,
				config.Update.Message.From.UserName,
				start_time)
			errors.CheckError(err)

			reply := fmt.Sprintf("Hello, [%v], the developer of this bot is @WB31B The bot was created to display the weather of the region you specified. Write the city and the Bot will tell you the weather",
				config.Update.Message.From.FirstName)
			msgConfig := tgbotapi.NewMessage(config.Update.Message.Chat.ID, reply)
			config.Bot.Send(msgConfig)
			continue
		}

		if command == "stop" && config.Update.Message.From.ID == config.ROOTUSER {
			msgConfig := tgbotapi.NewMessage(config.Update.Message.From.ID, "Bot stoped!")
			config.Bot.Send(msgConfig)
			log.StopBotCommand(config.Update.Message.From.ID)
			config.Bot.StopReceivingUpdates()
			break
		} else if command == "users" && config.Update.Message.From.ID == config.ROOTUSER {
			log.OutputUsersCommand(config.ROOTUSER)
			rows, err := db.Query(config.USERSDB)
			errors.CheckError(err)

			defer rows.Close()

			for rows.Next() {
				ui := config.UserInfo{}
				err := rows.Scan(&ui.User_id, &ui.User_name, &ui.User_tgid, &ui.Start_time)
				if err != nil {
					fmt.Println(err)
					continue
				}
				config.User = append(config.User, ui)
			}

			// get document with users
			path, err := OutputUsers(config.User)
			errors.CheckError(err)
			data, _ := ioutil.ReadFile(path)
			msgFile := tgbotapi.FileBytes{Name: "usersDatabaseInfo.txt", Bytes: data}
			msgConfig := tgbotapi.NewDocument(config.Update.Message.Chat.ID, msgFile)
			config.Bot.Send(msgConfig)
		} else if command == "" {
			log.ShowWeather(config.Update.Message.From.ID, config.Update.Message.Text)
			weather, err := weather.Weather(config.Update.Message.Text)
			errors.CheckError(err)

			weatherInfo, err := weatherTemperature(weather, config.Update)
			errors.CheckError(err)

			msgConfig := tgbotapi.NewMessage(config.Update.Message.From.ID, weatherInfo)
			config.Bot.Send(msgConfig)
		} else {
			if command == "start" {
				log.StartCommand(config.Update.Message.From.ID)
				msgConfig := tgbotapi.NewMessage(config.Update.Message.From.ID,
					"Hello, the developer of this bot is @WB31B The bot was created to display the weather of the region you specified. Write the city and the Bot will tell you the weather")
				config.Bot.Send(msgConfig)
			} else {
				log.IncorrectCommand(config.Update.Message.From.ID)
				msgConfig := tgbotapi.NewMessage(config.Update.Message.From.ID, "This command is INCORRECT!")
				config.Bot.Send(msgConfig)
			}

		}
	}

}

func weatherTemperature(weather *weather.WeatherData, update tgbotapi.Update) (string, error) {
	if weather.Data.Values.Temperature < 15 {
		weatherInfo := fmt.Sprintf("%s\n\n👨‍💻 User ID: [%v]\n🌍 Country: %v\n🥶 Temperature: %v\n💧 Humidity: %v\n☁️ Cloud Cover: %v\n💨 Visibility: %v\n\n⏰ Latest update time: %v\n",
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
		weatherInfo := fmt.Sprintf("%s\n\n👨‍💻 User ID: [%v]\n🌍 Country: %v\n🥵 Temperature: %v\n💧 Humidity: %v\n☁️ Cloud Cover: %v\n💨 Visibility: %v\n\n⏰ Latest update time: %v\n",
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

func OutputUsers(user []config.UserInfo) (string, error) {
	path := "usersDatabaseInfo.txt"
	file, err := os.Create(path)
	errors.CheckError(err)

	defer file.Close()

	for index, ui := range user {
		writingInFile(file, ui, index)
	}

	return path, nil
}

func writingInFile(file *os.File, user config.UserInfo, index int) {
	userInfo := fmt.Sprintf("[%d] Username: %v | User ID: %v | User TG Name: %v\n",
		index, user.User_name, user.User_id, user.User_tgid)
	_, err := io.Copy(file, strings.NewReader(userInfo))
	errors.CheckError(err)
}
