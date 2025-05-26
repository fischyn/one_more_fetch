package cpu

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
