package log

import "fmt"

func OutputDeleteUser(userId int64) {
	fmt.Printf("[Delete user] id: %v\n", userId)
}

func OutputAddUser(userId int64) {
	fmt.Printf("[Add user] id: %v\n", userId)
}
