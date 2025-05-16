package sys

import (
	"github.com/shirou/gopsutil/v4/mem"
)

type MemUsageStat struct {
	Total uint64 `json:"total"`
	Used  uint64 `json:"used"`
}

func MemUsage() (*MemUsageStat, error) {
	vMemStat, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &MemUsageStat{
		Total: vMemStat.Total / 1024 / 1024,
		Used:  vMemStat.Used / 1024 / 1024,
	}, nil
}
