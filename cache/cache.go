package cache

import (
	"os"
)

const ReadOnlyPerm = 0444

func SaveCPUCache(filename string, data []byte) error {
	return os.WriteFile(filename, data, ReadOnlyPerm)
}

func LoadCPUCache(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}
