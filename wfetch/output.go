package main

import (
	"fmt"
	"log"
	"os/user"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/fischyn/wfetch/img/ascii"
	"github.com/fischyn/wfetch/sys/bios"
	"github.com/fischyn/wfetch/sys/cpu"
	"github.com/fischyn/wfetch/sys/gpu"
	"github.com/fischyn/wfetch/sys/host"
	"github.com/fischyn/wfetch/sys/mem"
	"github.com/fischyn/wfetch/sys/platform"
)

const DELTA = 9

var (
	biosColor       = color.New(color.FgYellow).SprintFunc()
	cpuColor        = color.New(color.FgCyan).SprintFunc()
	gpuColor        = color.New(color.FgBlue).SprintFunc()
	memoryColor     = color.New(color.FgMagenta).SprintFunc()
	resolutionColor = color.New(color.FgGreen).SprintFunc()
	platformColor   = color.New(color.FgHiBlue).SprintFunc()
	hostColor       = color.New(color.FgWhite).SprintFunc()
	userColor       = color.New(color.FgRed).SprintFunc()
	asciiColor      = color.New(color.FgGreen).SprintFunc()
)

func run(wg *sync.WaitGroup, errCh chan<- error, f func() error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := f(); err != nil {
			errCh <- err
		}
	}()
}

func Fetch() {
	var (
		biosResult     bios.BIOSResult
		cpuResult      cpu.CPUResult
		resolution     gpu.Resolution
		platformResult platform.PlatformResult
		memoryResult   mem.MemoryResult
		hostname       string
		userResult     user.User
		gpus           []gpu.GPUResult
	)

	var (
		muHost sync.Mutex
		muUser sync.Mutex
	)

	var wg sync.WaitGroup
	errCh := make(chan error, DELTA)

	start := time.Now()

	run(&wg, errCh, func() error { return gpu.GetResolution(&resolution) })
	run(&wg, errCh, func() error { return bios.GetBIOS(&biosResult) })
	run(&wg, errCh, func() error {
		if err := cpu.GetNCores(&cpuResult); err != nil {
			return err
		}
		return cpu.GetRegistryData(&cpuResult)
	})
	run(&wg, errCh, func() error { return gpu.GetGPUsInfo(&gpus) })
	run(&wg, errCh, func() error {
		h, err := host.GetHostname()
		if err != nil {
			return err
		}
		muHost.Lock()
		hostname = h
		muHost.Unlock()
		return nil
	})
	run(&wg, errCh, func() error { return mem.GetMemory(&memoryResult) })
	run(&wg, errCh, func() error { return platform.GetPlatform(&platformResult) })
	run(&wg, errCh, func() error {
		u, err := host.GetUser()
		if err != nil {
			return err
		}
		muUser.Lock()
		userResult = *u
		muUser.Unlock()
		return nil
	})

	wg.Wait()
	close(errCh)

	for err := range errCh {
		log.Println("Error:", err)
	}

	rightLines := []string{
		hostColor("Hostname:"),
		fmt.Sprintf("  %s", hostname),
		"",
		userColor("Username:"),
		fmt.Sprintf("  %s", userResult.Username),
		"",
		cpuColor("CPU Info:"),
		fmt.Sprintf("  Name: %s", cpuResult.ProcessorName),
		fmt.Sprintf("  Vendor: %s", cpuResult.Vendor),
		fmt.Sprintf("  Identifier: %s", cpuResult.Identifier),
		fmt.Sprintf("  Speed: %d MHz", cpuResult.Mhz),
		fmt.Sprintf("  Cores (Physical/Logical/Active): %d/%d/%d", cpuResult.CoresPhysical, cpuResult.CoresLogical, cpuResult.CoresActive),
		fmt.Sprintf("  Packages: %d", cpuResult.Packages),
		"",
		biosColor("BIOS Info:"),
		fmt.Sprintf("  Vendor: %s", biosResult.BiosVendor),
		fmt.Sprintf("  Version: %s", biosResult.BiosVersion),
		fmt.Sprintf("  Release Date: %s", biosResult.BiosReleaseDate),
		fmt.Sprintf("  System Product: %s", biosResult.SystemProduct),
		fmt.Sprintf("  Manufacturer: %s", biosResult.SystemManufacturer),
		"",
		gpuColor("GPU Info:"),
	}

	for i, gpu := range gpus {
		rightLines = append(rightLines,
			fmt.Sprintf("  GPU #%d Name: %s", i+1, gpu.Name),
			fmt.Sprintf("  Vendor: %s", gpu.Vendor),
			fmt.Sprintf("  Video Memory: %d MB", gpu.VideoMemory.MBytes),
			fmt.Sprintf("  Driver Version: %s", gpu.DriverVersion),
			fmt.Sprintf("  Max D3D12 Feature Level: %s", gpu.MaxD3D12FeatureLevel),
			fmt.Sprintf("  Max D3D11 Feature Level: %s", gpu.MaxD3D11FeatureLevel),
			"",
		)
	}

	rightLines = append(rightLines,
		resolutionColor("Resolution:"),
		fmt.Sprintf("  %dx%d", resolution.Horizontal, resolution.Vertical),
		"",
		memoryColor("Memory Info:"),
		fmt.Sprintf("  Total: %d MB", memoryResult.Total/1024/1024),
		fmt.Sprintf("  Available: %d MB", memoryResult.Avail/1024/1024),
		fmt.Sprintf("  Used: %d MB", memoryResult.Used/1024/1024),
		fmt.Sprintf("  Used Percent: %.2f%%", memoryResult.UsedPercent),
		"",
		platformColor("Platform Info:"),
		fmt.Sprintf("  Platform: %s", platformResult.Platform),
		fmt.Sprintf("  Family: %s", platformResult.Family),
		fmt.Sprintf("  Version: %s", platformResult.Version),
		fmt.Sprintf("  Display Version: %s", platformResult.DisplayVersion),
	)

	leftLinesRaw := strings.Split(strings.TrimRight(ascii.Allien, "\n"), "\n")

	maxLeftWidth := 0
	for _, line := range leftLinesRaw {
		if len(line) > maxLeftWidth {
			maxLeftWidth = len(line)
		}
	}

	maxLines := len(leftLinesRaw)
	if len(rightLines) > maxLines {
		maxLines = len(rightLines)
	}

	for i := 0; i < maxLines; i++ {
		leftRaw := ""
		right := ""

		if i < len(leftLinesRaw) {
			leftRaw = leftLinesRaw[i]
		}
		if i < len(rightLines) {
			right = rightLines[i]
		}

		padding := strings.Repeat(" ", maxLeftWidth-len(leftRaw))
		leftColored := asciiColor(leftRaw + padding)

		fmt.Printf("%s │ %s\n", leftColored, right)
	}

	fmt.Println("Total fetching time:", time.Since(start))
}
