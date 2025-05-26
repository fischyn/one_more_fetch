package gpu

import (
	"fmt"
)

type GPUResult struct {
	Name                 string      `json:"name"`
	Vendor               string      `json:"vendor"`
	VideoMemory          VideoMemory `json:"videoMemory"`
	MaxD3D12FeatureLevel string      `json:"maxD3D12FeatureLevel"`
	MaxD3D11FeatureLevel string      `json:"maxD3D11FeatureLevel"`
	DriverVersion        string      `json:"driverVersion"`
}

type VideoMemory struct {
	Bytes  uint64 `json:"bytes"`
	MBytes uint32 `json:"mBytes"`
	GBytes uint16 `json:"gBytes"`
}

type Resolution struct {
	Horizontal int `json:"horizontal"`
	Vertical   int `json:"vertical"`
}

func DefineGPUVendor(vendorID uint32) string {
	switch vendorID {
	case 0x106b:
		return "Apple"
	case 0x1002, 0x1022:
		return "AMD"
	case 0x8086, 0x8087, 0x03e7:
		return "Intel"
	case 0x0955, 0x10de, 0x12d2:
		return "NVIDIA"
	case 0x1ed5:
		return "Moore Threads"
	case 0x5143:
		return "Qualcomm"
	case 0x14c3:
		return "MediaTek"
	case 0x15ad:
		return "VMware"
	case 0x1af4:
		return "Red Hat"
	case 0x1ab8:
		return "Parallels"
	case 0x1414:
		return "Microsoft"
	case 0x108e:
		return "Oracle"
	default:
		return ""
	}
}

func DecodeD3D12FeatureLevel(featureLevel uint32) (string, error) {
	if featureLevel == 0 {
		return "", fmt.Errorf("unknown Direct3D feature level")
	}

	major := 12
	minor := (featureLevel & 0x0F00) >> 8

	return fmt.Sprintf("Direct3D %d.%d", major, minor), nil
}

func DecodeD3FeatureLevel(featureLevel uint32) (string, error) {
	if featureLevel == 0 {
		return "", fmt.Errorf("unknown Direct3D feature level")
	}

	major := (featureLevel & 0xF000) >> 12 // upper nibble (4 bits)
	minor := (featureLevel & 0x0F00) >> 8  // next nibble (4 bits)

	switch major {
	case 0xC: // 12 decimal
		return fmt.Sprintf("Direct3D 12.%d", minor), nil
	case 0xB: // 11 decimal
		return fmt.Sprintf("Direct3D 11.%d", minor), nil
	case 0xA: // 10 decimal
		return fmt.Sprintf("Direct3D 10.%d", minor), nil
	default:
		return "", fmt.Errorf("unsupported Direct3D feature level: 0x%X", featureLevel)
	}
}

func DecodeDriverVersion(version uint64) string {
	major := uint16((version >> 48) & 0xFFFF)
	minor := uint16((version >> 32) & 0xFFFF)
	patch := uint16((version >> 16) & 0xFFFF)
	build := uint16(version & 0xFFFF)

	return fmt.Sprintf("%d.%d.%d.%d", major, minor, patch, build)
}
