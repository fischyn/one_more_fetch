package mem

type MemoryResult struct {
	Total       uint64  `json:"total"`
	Avail       uint64  `json:"avail"`
	Used        uint64  `json:"used"`
	UsedPercent float64 `json:"usedPercent"`
}
