package user

import (
	"os/user"
)

var (
	getUser = user.Current
)

func GetUserInfo() (*user.User, error) {
	return getUser()
}
