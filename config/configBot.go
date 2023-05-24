package config

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type UserInfo struct {
	User_id    int64
	User_name  string
	User_tgid  string
	Start_time string
}

var (
	User_id   int64
	User_name string
	User_tgid string
)

var (
	Bot        *tgbotapi.BotAPI
	UpdChannel tgbotapi.UpdatesChannel
	Update     tgbotapi.Update
	UpdConfig  tgbotapi.UpdateConfig
	User       []UserInfo
)
