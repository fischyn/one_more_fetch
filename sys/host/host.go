package host

import (
	"os"
)

var (
	getHostname = os.Hostname
)

func GetHostname() (string, error) {
	return getHostname()
}
