package main

import (
	"encoding/json"
	"fmt"

	"github/com/fischyn/omfetch/sys"
)

func main() {
	memUsageStat, memErr := sys.MemUsage()

	if memErr != nil {
		fmt.Printf("err: %v\n", memErr)
	}

	cpuInfo, cpuErr := sys.CpuInfo()

	if cpuErr != nil {
		fmt.Printf("err: %v\n", cpuErr)
	}

	memJsonData, memMarshalErr := json.Marshal(memUsageStat)
	cpuJsonData, cpuMarshalErr := json.Marshal(cpuInfo)

	if memMarshalErr != nil {
		fmt.Printf("json marshal error: %v\n", memMarshalErr)
		return
	}

	if cpuMarshalErr != nil {
		fmt.Printf("json marshal error: %v\n", cpuMarshalErr)
	}

	fmt.Println(string(memJsonData))
	fmt.Println(string(cpuJsonData))
}
