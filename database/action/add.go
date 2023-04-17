package action

import (
	"TGbot/log"
)

func AddUser(users []int64, userId int64) ([]int64, error) {
	users = append(users, userId)
	log.OutputAddUser(userId)
	return users, nil
}
