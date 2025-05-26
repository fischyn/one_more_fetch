package platform

type PlatformResult struct {
	Platform       string `json:"platform"`
	Family         string `json:"family"`
	Version        string `json:"version"`
	DisplayVersion string `json:"displayVersion"`
}

type PlatformOptions struct {
	ShowPlatform       bool `json:"platform"`
	ShowFamily         bool `json:"family"`
	ShowVersion        bool `json:"version"`
	ShowDisplayVersion bool `json:"displayVersion"`
}
