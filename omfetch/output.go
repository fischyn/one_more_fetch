package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/user"
	"sync"
	"time"

	"github.com/fischyn/omfetch/sys/bios"
	"github.com/fischyn/omfetch/sys/cpu"
	"github.com/fischyn/omfetch/sys/gpu"
	"github.com/fischyn/omfetch/sys/host"
	"github.com/fischyn/omfetch/sys/mem"
	"github.com/fischyn/omfetch/sys/platform"
)

const DELTA = 9

var (
	biosOpt = bios.BIOSOptions{
		ShowSystemProduct:         true,
		ShowSystemManufacturer:    true,
		ShowSystemFamily:          true,
		ShowSystemVersion:         true,
		ShowBiosVendor:            true,
		ShowBiosVersion:           true,
		ShowBiosReleaseDate:       true,
		ShowBaseBoardManufacturer: true,
		ShowBaseBoardProduct:      true,
		ShowBaseBoardVersion:      true,
	}

	platformOpt = platform.PlatformOptions{
		ShowPlatform:       true,
		ShowFamily:         true,
		ShowVersion:        true,
		ShowDisplayVersion: true,
	}
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

func PrintInfo() {
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

	run(&wg, errCh, func() error {
		return gpu.GetResolution(&resolution)
	})

	run(&wg, errCh, func() error {
		return bios.GetBIOS(&biosResult, biosOpt)
	})

	run(&wg, errCh, func() error {
		if err := cpu.GetNCores(&cpuResult); err != nil {
			return err
		}
		return cpu.GetRegistryData(&cpuResult)
	})

	run(&wg, errCh, func() error {
		return gpu.GetGPUsInfo(&gpus)
	})

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

	run(&wg, errCh, func() error {
		return mem.GetMemory(&memoryResult)
	})

	run(&wg, errCh, func() error {
		return platform.GetPlatform(&platformResult, platformOpt)
	})

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

	fmt.Println("Total fetching time:", time.Since(start))

	for err := range errCh {
		log.Println("Error:", err)
	}

	printJSON("BIOS Info", biosResult)
	printJSON("CPU Info", cpuResult)
	printJSON("GPU Info", gpus)
	printJSON("Resolution", resolution)
	printJSON("Platform Info", platformResult)
	printJSON("Memory Info", memoryResult)
	fmt.Println("Hostname:", hostname)
	fmt.Println("Username:", userResult.Username)
}

func printJSON(label string, v any) {
	if j, err := json.MarshalIndent(v, "", "  "); err == nil {
		fmt.Println(label+":", string(j))
	} else {
		fmt.Println("Failed to marshal", label+":", err)
	}
}
