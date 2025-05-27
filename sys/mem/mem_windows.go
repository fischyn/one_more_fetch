//go:build windows

package mem

import (
	"unsafe"

	"github.com/fischyn/wfetch/sys/dll"

	"golang.org/x/sys/windows"
)

// https://learn.microsoft.com/en-us/windows/win32/api/sysinfoapi/nf-sysinfoapi-globalmemorystatusex
var globalMemoryStatusExProc = dll.Kernel32.NewProc("GlobalMemoryStatusEx")

// https://learn.microsoft.com/en-us/windows/win32/api/sysinfoapi/ns-sysinfoapi-memorystatusex
type memoryStatusEx struct {
	dwLength                uint32 // The size of the structure, in bytes. You must set this member before calling
	dwMemoryLoad            uint32
	ullTotalPhys            uint64
	ullAvailPhys            uint64
	ullTotalPageFile        uint64
	ullAvailPageFile        uint64
	ullTotalVirtual         uint64
	ullAvailVirtual         uint64
	ullAvailExtendedVirtual uint64
}

func GetMemory(m *MemoryResult) error {

	var memStatEx memoryStatusEx
	memStatEx.dwLength = uint32(unsafe.Sizeof(memStatEx))
	mem, _, _ := globalMemoryStatusExProc.Call(uintptr(unsafe.Pointer(&memStatEx)))

	if mem == 0 {
		return windows.GetLastError()
	}

	m.Total = memStatEx.ullTotalPhys
	m.Avail = memStatEx.ullAvailPhys
	m.UsedPercent = float64(memStatEx.dwMemoryLoad)

	m.Used = m.Total - m.Avail

	return nil
}
