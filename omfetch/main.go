package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github/com/fischyn/omfetch/sys/cpu"
	"github/com/fischyn/omfetch/sys/gpu"
	"github/com/fischyn/omfetch/sys/host"
	"github/com/fischyn/omfetch/sys/mem"
)

func main() {
	// Only for testing  now
	memoryInfo, _ := mem.GetMemoryInfo(context.Background())
	memoryJSData, _ := json.Marshal(memoryInfo)

	platform, family, version, displayVersion, _ := host.GetPlatformInfo(context.Background())

	cpuInfo, _ := cpu.GetCPUInfo(context.Background())
	cpuJSData, _ := json.Marshal(cpuInfo)

	gpuInfo, _ := gpu.GetGPUInfo(context.Background())
	gpuJSData, _ := json.Marshal(gpuInfo)

	fmt.Println(string(memoryJSData))
	fmt.Printf("Platfrom: %s\n", platform)
	fmt.Printf("Family: %s\n", family)
	fmt.Printf("version: %s\n", version)
	fmt.Printf("displayVersion: %s\n", displayVersion)
	fmt.Println(string(cpuJSData))
	fmt.Println(string(gpuJSData))
}
