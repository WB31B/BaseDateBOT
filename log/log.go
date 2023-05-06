package log

import (
	"TGbot/errors"
	"fmt"
	"os"
	"time"
)

var timeNow = time.Now()

func StartBot(userId int64) {
	logMSG := fmt.Sprintf("%v - [Add new user]: %v\n",
		timeNow.Format("15:04:05"), userId)
	createLogFile(logMSG)
}

func StartCommand(userId int64) {
	logMSG := fmt.Sprintf("%v - [Start command]: %v\n",
		timeNow.Format("15:04:05"), userId)
	createLogFile(logMSG)
}

func OutputUsersCommand(rootUser int64) {
	logMSG := fmt.Sprintf("%v - [Output all users]: ROOT :%v\n",
		timeNow.Format("15:04:05"), rootUser)
	createLogFile(logMSG)
}

func StopBotCommand(rootUser int64) {
	logMSG := fmt.Sprintf("%v - [Stop bot]: ROOT :%v\n",
		timeNow.Format("15:04:05"), rootUser)
	createLogFile(logMSG)
}

func IncorrectCommand(userId int64) {
	logMSG := fmt.Sprintf("%v - [Incorrect command]: %v\n",
		timeNow.Format("15:04:05"), userId)
	createLogFile(logMSG)
}

func ShowWeather(userId int64, weather string) {
	logMSG := fmt.Sprintf("%v - [Show weather]:%v: %v\n",
		timeNow.Format("15:04:05"), weather, userId)
	createLogFile(logMSG)
}

func createLogFile(msg string) {
	path := "logAction.data"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0600)
	errors.CheckError(err)

	defer file.Close()

	file.WriteString(msg)
}
