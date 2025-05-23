package cache

import (
	"os"
)

const ReadOnly = 0444

func SaveCPUCache(filename string, data []byte) error {
	return os.WriteFile(filename, data, ReadOnly)
}

func LoadCPUCache(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
