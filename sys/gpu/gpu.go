package gpu

type GPUInfo struct {
	GPUName                     string `json:"name"`
	MemoryBytes                 uint32 `json:"memory"`
	DriverVersion               string `json:"driverVersion"`
	VideoProcessor              string `json:"videoProcessor"`
	CurrentHorizontalResolution uint32 `json:"horizontalResolution"`
	CurrentVerticalResolution   uint32 `json:"verticalResolution"`
}
