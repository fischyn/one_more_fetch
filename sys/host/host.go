package host

import (
	"context"
	"os"
)

var (
	getHostname = os.Hostname
)

func GetHostname(_ context.Context) (string, error) {
	return getHostname()
}
