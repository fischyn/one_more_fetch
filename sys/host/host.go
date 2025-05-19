package host

import (
	"context"
	"os"
)

var (
	getHostname = os.Hostname
)

func GetHostname(_ context.Context) (string, error) {
	hostname, err := getHostname()
	if err != nil {
		return "", err
	}

	return hostname, nil
}
