package main

import (
	"context"
	"fmt"

	"github/com/fischyn/omfetch/sys/bios"
	"github/com/fischyn/omfetch/sys/cpu"
	"github/com/fischyn/omfetch/sys/gpu"
	"github/com/fischyn/omfetch/sys/host"
	"github/com/fischyn/omfetch/sys/mem"
)

const Logo = `
           ▄▄▄▄▄
        ▄█████████▄
      ▄███▀▀▀▀▀▀▀███▄
     ███▀          ▀███
    ███              ███
    ███              ███
     ███▄          ▄███
      ▀███▄▄▄▄▄▄▄███▀
         ▀▀▀▀▀▀▀▀▀
`

const Logo2 = "       .------..\n" +
	"     -          -\n" +
	"   /              \\\n" +
	" /                   \\\n" +
	"/    .--._    .---.   |\n" +
	"|  /      -__-     \\   |\n" +
	"| |                 |  |\n" +
	"||     ._   _.      ||\n" +
	"||      o   o       ||\n" +
	"||      _  |_      ||\n" +
	"C|     (o\\_/o)     |O     Uhhh, this computer\n" +
	" \\      _____      /       is like, busted or\n" +
	"   \\ ( /#####\\ ) /       something. So go away.\n" +
	"    \\  `====='  /\n" +
	"     \\  -___-  /\n" +
	"      |       |\n" +
	"      /-_____-\\" + "\n" +
	"    /           \\\n" +
	"  /               \\\n" +
	" /__|  AC / DC  |__\\\n" +
	" | ||           |\\ \\"

var (
	blue  = "\033[34m"
	green = "\033[32m"
	reset = "\033[0m"
)

func PrintInfo(label string, value string) {
	fmt.Printf("%s%-12s%s %s\n", green, label+":", reset, value)
}

func main() {
	//test only lol

	fmt.Println(blue + Logo2 + reset)

	platform, family, platformVersion, displayVersion, _ := host.GetPlatformInfo(context.Background())
	PrintInfo("OS", platform+" "+family+" "+platformVersion+" "+displayVersion)

	cpuInfo, _ := cpu.GetCPUInfo(context.Background())
	c1 := cpuInfo[0]

	PrintInfo("CPU", fmt.Sprintf("%s (%d) cores %f Mhz", c1.ModelName, c1.Cores, c1.Mhz))

	memoryInfo, _ := mem.GetMemoryInfo(context.Background())

	PrintInfo("RAM", fmt.Sprintf("%d MB", memoryInfo.Total/1024/1024))
	PrintInfo("RAM(Used)", fmt.Sprintf("%d MB", memoryInfo.Used/1024/1024))

	gpuInfo, _ := gpu.GetGPUInfo(context.Background())

	PrintInfo("GPU", gpuInfo[0].GPUName)

	product, manufacturer, family, biosVersion, _ := bios.GetBiosInfo(context.Background())

	PrintInfo("BIOS", fmt.Sprintf("%s %s %s %s", manufacturer, family, product, biosVersion))
}
