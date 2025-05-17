package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github/com/fischyn/omfetch/sys/host"
	"github/com/fischyn/omfetch/sys/mem"
)

func main() {
	// Only for testing  now
	memoryInfo, _ := mem.GetMemoryInfo(context.Background())
	memoryJSData, _ := json.Marshal(memoryInfo)

	platform, family, version, displayVersion, _ := host.GetPlatformInfo(context.Background())

	fmt.Println(string(memoryJSData))
	fmt.Printf("Platfrom: %s\n", platform)
	fmt.Printf("Family: %s\n", family)
	fmt.Printf("version: %s\n", version)
	fmt.Printf("displayVersion: %s\n", displayVersion)
}
