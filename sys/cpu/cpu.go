package cpu

type CPUInfo struct {
	CPU        int32    `json:"cpu"`
	VendorID   string   `json:"vendorId"`
	Family     string   `json:"family"`
	Model      string   `json:"model"`
	Stepping   int32    `json:"stepping"`
	PhysicalID string   `json:"physicalId"`
	CoreID     string   `json:"coreId"`
	Cores      int32    `json:"cores"`
	ModelName  string   `json:"modelName"`
	Mhz        float64  `json:"mhz"`
	CacheSize  int32    `json:"cacheSize"`
	Flags      []string `json:"flags"`
	Microcode  string   `json:"microcode"`
}

type CPUResult struct {
	ProcessorName string // 2
	Vendor        string // 2
	Identifier    string // 2
	Mhz           uint32 // 4
	CoresPhysical uint16 // 2
	CoresLogical  uint16 // 2
	CoresActive   uint16 // 2
	Packages      uint16 // 2
}

// 18 bytes
