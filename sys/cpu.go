package sys

import (
	"github.com/shirou/gopsutil/v4/cpu"
)

type CPUInfo struct {
	VendorID  string  `json:"vendorId"`
	Family    string  `json:"family"`
	Model     string  `json:"model"`
	Cores     int32   `json:"cores"`
	ModelName string  `json:"modelName"`
	Mhz       float64 `json:"mhz"`
}

func CpuInfo() ([]CPUInfo, error) {
	dst, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	var res []CPUInfo

	for _, c := range dst {
		cpu := CPUInfo{
			VendorID:  c.VendorID,
			Family:    c.Family,
			Model:     c.Model,
			Cores:     c.Cores,
			ModelName: c.ModelName,
			Mhz:       c.Mhz,
		}
		res = append(res, cpu)
	}

	return res, nil
}
