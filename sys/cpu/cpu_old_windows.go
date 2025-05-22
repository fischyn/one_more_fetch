//go:build windows

package cpu

import (
	"context"
	"strconv"

	"github.com/fischyn/omfetch/sys/windows/wql"

	"github.com/yusufpapurcu/wmi"
)

type win32_Processor struct {
	Family                    uint16
	Manufacturer              string
	Name                      string
	NumberOfLogicalProcessors uint32
	NumberOfCores             uint32
	ProcessorID               *string
	Stepping                  *string
	MaxClockSpeed             uint32
}

// Deprecated
func GetCPUInfo(ctx context.Context) ([]CPUInfo, error) {
	var ret []CPUInfo
	var dst []win32_Processor
	q := wmi.CreateQuery(&dst, "")
	if err := wql.WMIQuery(ctx, q, &dst); err != nil {
		return ret, err
	}

	var procID string
	for i, l := range dst {
		procID = ""
		if l.ProcessorID != nil {
			procID = *l.ProcessorID
		}

		cpu := CPUInfo{
			CPU:        int32(i),
			Family:     strconv.FormatUint(uint64(l.Family), 10),
			VendorID:   l.Manufacturer,
			ModelName:  l.Name,
			Cores:      int32(l.NumberOfLogicalProcessors),
			PhysicalID: procID,
			Mhz:        float64(l.MaxClockSpeed),
			Flags:      []string{},
		}
		ret = append(ret, cpu)
	}

	return ret, nil
}
