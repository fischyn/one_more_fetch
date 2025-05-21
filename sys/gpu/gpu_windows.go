//go:build windows

package gpu

import (
	"context"

	"github.com/fischyn/omfetch/sys/windows/wql"

	"github.com/yusufpapurcu/wmi"
)

// https://learn.microsoft.com/ru-ru/windows/win32/cimwin32prov/win32-videocontroller
type win32_VideoController struct {
	Name                        string
	AdapterRAM                  *uint32
	DriverVersion               string
	VideoProcessor              string
	CurrentHorizontalResolution *uint32
	CurrentVerticalResolution   *uint32
}

func GetGPUInfo(ctx context.Context) ([]GPUInfo, error) {
	var ret []GPUInfo
	var dst []win32_VideoController
	q := wmi.CreateQuery(&dst, "")
	if err := wql.WMIQuery(ctx, q, &dst); err != nil {
		return ret, err
	}

	for _, gpu := range dst {
		memory := uint32(0)
		if gpu.AdapterRAM != nil {
			memory = *gpu.AdapterRAM
		}

		horRes := uint32(0)
		if gpu.CurrentHorizontalResolution != nil {
			horRes = *gpu.CurrentHorizontalResolution
		}

		verRes := uint32(0)
		if gpu.CurrentVerticalResolution != nil {
			verRes = *gpu.CurrentVerticalResolution
		}

		info := GPUInfo{
			GPUName:                     gpu.Name,
			MemoryBytes:                 memory,
			DriverVersion:               gpu.DriverVersion,
			VideoProcessor:              gpu.VideoProcessor,
			CurrentHorizontalResolution: horRes,
			CurrentVerticalResolution:   verRes,
		}
		ret = append(ret, info)
	}

	return ret, nil
}
