package bios

type BIOSResult struct {
	SystemProduct         string `json:"systemProduct"`
	SystemManufacturer    string `json:"systemManufacturer"`
	SystemFamily          string `json:"systemFamily"`
	SystemVersion         string `json:"systemVersion"`
	BiosVendor            string `json:"biosVendor"`
	BiosVersion           string `json:"biosVersion"`
	BiosReleaseDate       string `json:"biosReleaseDate"`
	BaseBoardManufacturer string `json:"baseBoardManufacturer"`
	BaseBoardProduct      string `json:"baseBoardProduct"`
	BaseBoardVersion      string `json:"baseBoardVersion"`
}

type BIOSOptions struct {
	ShowSystemProduct         bool `json:"systemProduct"`
	ShowSystemManufacturer    bool `json:"systemManufacturer"`
	ShowSystemFamily          bool `json:"systemFamily"`
	ShowSystemVersion         bool `json:"systemVersion"`
	ShowBiosVendor            bool `json:"biosVendor"`
	ShowBiosVersion           bool `json:"biosVersion"`
	ShowBiosReleaseDate       bool `json:"biosReleaseDate"`
	ShowBaseBoardManufacturer bool `json:"baseBoardManufacturer"`
	ShowBaseBoardProduct      bool `json:"baseBoardProduct"`
	ShowBaseBoardVersion      bool `json:"baseBoardVersion"`
}
