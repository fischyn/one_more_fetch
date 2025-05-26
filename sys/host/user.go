package host

import (
	"os/user"
)

var (
	getUser = user.Current
)

func GetUser() (*user.User, error) {
	return getUser()
}
