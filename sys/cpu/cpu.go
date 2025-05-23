package cpu

//Deprecated
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
	ProcessorName string `json:"processorName"`
	Vendor        string `json:"vendor"`
	Identifier    string `json:"identifier"`
	Mhz           uint32 `json:"mhz"`
	CoresPhysical uint16 `json:"coresPhysical"`
	CoresLogical  uint16 `json:"coresLogical"`
	CoresActive   uint16 `json:"coresActive"`
	Packages      uint16 `json:"packages"`
}
