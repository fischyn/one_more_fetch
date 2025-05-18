//go:build windows

package gpu

import (
	"context"

	"github/com/fischyn/omfetch/sys/windows/wql"

	"github.com/yusufpapurcu/wmi"
)

type win32_VideoController struct {
	Name           string
	AdapterRAM     *uint64
	DriverVersion  string
	VideoProcessor string
}

func GetGPUInfo(ctx context.Context) ([]GPUInfo, error) {
	var ret []GPUInfo
	var dst []win32_VideoController
	q := wmi.CreateQuery(&dst, "")
	if err := wql.WMIQuery(ctx, q, &dst); err != nil {
		return ret, err
	}

	for _, gpu := range dst {
		memory := uint64(0)
		if gpu.AdapterRAM != nil {
			memory = *gpu.AdapterRAM
		}

		info := GPUInfo{
			GPUName:        gpu.Name,
			MemoryBytes:    memory,
			DriverVersion:  gpu.DriverVersion,
			VideoProcessor: gpu.VideoProcessor,
		}
		ret = append(ret, info)
	}

	return ret, nil
}
