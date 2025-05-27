//go:build windows

package dll

import "golang.org/x/sys/windows"

var (
	Kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	User32   = windows.NewLazySystemDLL("user32.dll")
	SetupApi = windows.NewLazySystemDLL("setupapi.dll")
)
