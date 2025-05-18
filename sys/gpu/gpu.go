package gpu

type GPUInfo struct {
	GPUName        string `json:"name"`
	MemoryBytes    uint64 `json:"memory"`
	DriverVersion  string `json:"driverVersion"`
	VideoProcessor string `json:"videoProcessor"`
}
