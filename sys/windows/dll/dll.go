//go:build windows

package dll

import "golang.org/x/sys/windows"

var Kernel32 = windows.NewLazySystemDLL("Kernel32.dll")
