package main

import (
	"context"
	"fmt"

	"github/com/fischyn/omfetch/img/ascii"
	"github/com/fischyn/omfetch/sys/bios"
	"github/com/fischyn/omfetch/sys/cpu"
	"github/com/fischyn/omfetch/sys/gpu"
	"github/com/fischyn/omfetch/sys/host"
	"github/com/fischyn/omfetch/sys/mem"
)

var (
	blue  = "\033[34m"
	green = "\033[32m"
	reset = "\033[0m"
)

func PrintInfo(label string, value string) {
	fmt.Printf("%s%-12s%s %s\n", green, label+":", reset, value)
}

func formatMB(bytes uint64) string {
	return fmt.Sprintf("%d MB", bytes/1024/1024)
}

func formatCPU(model string, cores int32, mhz float64) string {
	return fmt.Sprintf("%s (%d cores) %.2f MHz", model, cores, mhz)
}

func printPlatformInfo(ctx context.Context) {
	platform, family, version, displayVersion, err := host.GetPlatformInfo(ctx)
	if err == nil {
		PrintInfo("OS", fmt.Sprintf("%s %s %s %s", platform, family, version, displayVersion))
	}
}

func printCPUInfo(ctx context.Context) {
	cpuInfo, err := cpu.GetCPUInfo(ctx)
	if err == nil && len(cpuInfo) > 0 {
		c := cpuInfo[0]
		PrintInfo("CPU", formatCPU(c.ModelName, c.Cores, c.Mhz))
	}
}

func printMemoryInfo(ctx context.Context) {
	memInfo, err := mem.GetMemoryInfo(ctx)
	if err == nil {
		PrintInfo("RAM", formatMB(memInfo.Total))
		PrintInfo("RAM(Used)", formatMB(memInfo.Used))
	}
}

func printGPUInfo(ctx context.Context) {
	gpuInfo, err := gpu.GetGPUInfo(ctx)
	if err == nil && len(gpuInfo) > 0 {
		PrintInfo("GPU", gpuInfo[0].GPUName)
	}
}

func printBiosInfo(ctx context.Context) {
	product, manufacturer, family, version, err := bios.GetBiosInfo(ctx)
	if err == nil {
		PrintInfo("BIOS", fmt.Sprintf("%s %s %s %s", manufacturer, family, product, version))
	}
}

func main() {
	fmt.Println(blue + ascii.Blinky + reset)

	ctx := context.Background()

	printPlatformInfo(ctx)
	printCPUInfo(ctx)
	printMemoryInfo(ctx)
	printGPUInfo(ctx)
	printBiosInfo(ctx)
}
