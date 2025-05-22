package main

import (
	"context"
	"fmt"
	"time"

	"github.com/fischyn/omfetch/img/ascii"
	"github.com/fischyn/omfetch/sys/bios"
	"github.com/fischyn/omfetch/sys/cpu"
	"github.com/fischyn/omfetch/sys/gpu"
	"github.com/fischyn/omfetch/sys/host"
	"github.com/fischyn/omfetch/sys/mem"
	"github.com/fischyn/omfetch/sys/platform"
	"github.com/fischyn/omfetch/sys/user"
)

var (
	blue  = "\033[34m"
	green = "\033[32m"
	reset = "\033[0m"
)

func printInfo(label string, value string) {
	fmt.Printf("%s%-12s%s %s\n", green, label+":", reset, value)
}

func formatMB(bytes uint64) string {
	return fmt.Sprintf("%d MB", bytes/1024/1024)
}

func formatCPU(model string, cores int32, mhz float64) string {
	return fmt.Sprintf("%s (%d cores) %.2f MHz", model, cores, mhz)
}

func formatResolution(horizontal uint32, vertical uint32) string {
	return fmt.Sprintf("%dx%d", horizontal, vertical)
}

func printHostname() {
	hostname, err := host.GetHostname()
	if err == nil {
		printInfo("Hostname", hostname)
	}
}

func printUserInfo() {
	user, err := user.GetUserInfo()
	if err == nil {
		printInfo("User", user.Username)
	}
}

func printPlatformInfo() {
	platform, family, version, displayVersion, err := platform.GetPlatformInfo()
	if err == nil {
		printInfo("OS", fmt.Sprintf("%s %s %s %s", platform, family, version, displayVersion))
	}
}

func printCPUInfo(ctx context.Context) {
	cpuInfo, err := cpu.GetCPUInfo(ctx)
	if err == nil && len(cpuInfo) > 0 {
		c := cpuInfo[0]
		printInfo("CPU", formatCPU(c.ModelName, c.Cores, c.Mhz))
	}
}

func printMemoryInfo() {
	memInfo, err := mem.GetMemoryInfo()
	if err == nil {
		printInfo("RAM", formatMB(memInfo.Total))
		printInfo("RAM(Used)", formatMB(memInfo.Used))
	}
}

func printGPUInfo(ctx context.Context) {
	gpuInfo, err := gpu.GetGPUInfo(ctx)
	if err == nil && len(gpuInfo) > 0 {
		gpu1 := gpuInfo[0] //skip additional gpu such as Intel Integrated... for now

		printInfo("GPU", gpu1.GPUName)
		printInfo("Resolution", formatResolution(gpu1.CurrentHorizontalResolution, gpu1.CurrentVerticalResolution))
	}
}

func printCPUCoresInfo() {
	var cp cpu.CPUResult

	err := cpu.GetNCores(&cp)
	if err != nil {
		fmt.Printf("err.Error(): %v\n", err.Error())
	}
	printInfo("CPU(Cores)", fmt.Sprintf("Phys cores: %d, Log cores: %d, Active cores: %d, Package :%d", cp.CoresPhysical, cp.CoresLogical, cp.CoresActive, cp.Packages))
}

func printBiosInfo() {
	product, manufacturer, family, version, err := bios.GetBiosInfo()
	if err == nil {
		printInfo("BIOS", fmt.Sprintf("%s %s %s %s", manufacturer, family, product, version))
	}
}

func main() {
	fmt.Println(blue + ascii.Beavis + reset)

	// ctx := context.Background()

	start := time.Now()

	// startFetchUserInfo := time.Now()
	printUserInfo()
	// fmt.Println("Time since fetching user info:", time.Since(startFetchUserInfo))

	// startFetchHostname := time.Now()
	printHostname()
	// fmt.Println("Time since fetching hostname:", time.Since(startFetchHostname))

	// startFetchPlatformInfo := time.Now()
	printPlatformInfo()
	// fmt.Println("Time since fetching platform info:", time.Since(startFetchPlatformInfo))

	// startFetchCpuInfo := time.Now()
	// printCPUInfo(ctx)
	// fmt.Println("Time since fetching cpu info:", time.Since(startFetchCpuInfo))

	// startFetchCpuCoresInfo := time.Now()
	printCPUCoresInfo()
	// fmt.Println("Time since fetching cpu cores info:", time.Since(startFetchCpuCoresInfo))

	// startFetchMemoryInfo := time.Now()
	printMemoryInfo()
	// fmt.Println("Time since fetching mem info:", time.Since(startFetchMemoryInfo))

	// startFetchGpuInfo := time.Now()
	// printGPUInfo(ctx)
	// fmt.Println("Time since fetching gpu info:", time.Since(startFetchGpuInfo))

	// startFetchBiosInfo := time.Now()
	printBiosInfo()
	// fmt.Println("Time since fetching bios info:", time.Since(startFetchBiosInfo))

	fmt.Println("Total execution time:", time.Since(start))
}
