package user

import (
	"context"
	"os/user"
)

var (
	getUser = user.Current
)

func GetUserInfo(_ context.Context) (*user.User, error) {
	return getUser()
}
