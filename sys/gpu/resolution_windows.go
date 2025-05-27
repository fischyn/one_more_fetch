//go:build windows

package gpu

import (
	"github.com/fischyn/wfetch/sys/dll"
	"golang.org/x/sys/windows"
)

const (
	SM_CXSCREEN int = 0
	SM_CYSCREEN int = 1
)

var (
	getSystemMetricsProc = dll.User32.NewProc("GetSystemMetrics")
	setProcessDPIAware   = dll.User32.NewProc("SetProcessDPIAware")
)

// This prevents Windows from scaling DPI for the application,
// allowing functions like GetSystemMetrics to return the
// true physical screen resolution instead of scaled values.
func init() {
	setProcessDPIAware.Call()
}

func GetResolution(resn *Resolution) error {

	horizontal, _, _ := getSystemMetricsProc.Call(uintptr(SM_CXSCREEN))
	vertical, _, _ := getSystemMetricsProc.Call(uintptr(SM_CYSCREEN))

	if horizontal == 0 || vertical == 0 {
		return windows.GetLastError()
	}

	resn.Horizontal = int(horizontal)
	resn.Vertical = int(vertical)

	return nil
}
