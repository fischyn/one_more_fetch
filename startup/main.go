package main

import (
	"encoding/json"
	"fmt"

	"github/com/fischyn/omfetch/sys"
)

func main() {
	memUsageStat, err := sys.MemUsage()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	jsonData, err := json.Marshal(memUsageStat)
	if err != nil {
		fmt.Printf("json marshal error: %v\n", err)
		return
	}

	fmt.Println(string(jsonData))
}
