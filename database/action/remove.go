package action

import (
	"TGbot/log"
)

func DeleteUser(users []int64, userId int64) ([]int64, error) {
	for index, user := range users {
		if user == userId {
			log.OutputDeleteUser(userId)
			users = append(users[:index], users[index+1:]...)
		}
	}

	return users, nil
}
