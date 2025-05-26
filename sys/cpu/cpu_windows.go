//go:build windows

package cpu

import (
	"unsafe"

	"github.com/fischyn/omfetch/sys/windows/dll"
	"github.com/fischyn/omfetch/sys/windows/registry"

	"golang.org/x/sys/windows"
)

// https://learn.microsoft.com/ru-ru/windows/win32/api/winnt/ne-winnt-logical_processor_relationship
type LOGICAL_PROCESSOR_RELATIONSHIP = uint16

const (
	RelationProcessorCore    LOGICAL_PROCESSOR_RELATIONSHIP = 0
	RelationProcessorPackage LOGICAL_PROCESSOR_RELATIONSHIP = 3
	RelationGroup            LOGICAL_PROCESSOR_RELATIONSHIP = 4
	RelationAll              LOGICAL_PROCESSOR_RELATIONSHIP = 0xffff
)

type SYSTEM_LOGICAL_PROCESSOR_INFORMATION_EX struct {
	Relationship uint16
	Size         uint32
	Data         [1]byte // super-duper hack?
}

type GROUP_RELATIONSHIP struct {
	MaximumGroupCount uint16
	ActiveGroupCount  uint16
	Reserved          [20]byte
}

type PROCESSOR_GROUP_INFO struct {
	MaximumProcessorCount uint8
	ActiveProcessorCount  uint8
	Reserved              [38]byte
	ActiveProcessorMask   uint64
}

var (
	procGetLogicalProcessorInformationEx = dll.Kernel32.NewProc("GetLogicalProcessorInformationEx")
)

// GetNCores queries Windows API to count physical cores, logical cores, active cores, and CPU packages.
//
// It uses GetLogicalProcessorInformationEx with RelationAll to get all processor relationship information.
// https://learn.microsoft.com/ru-ru/windows/win32/api/sysinfoapi/nf-sysinfoapi-getlogicalprocessorinformationex
// The function reads the returned buffer, parses each SYSTEM_LOGICAL_PROCESSOR_INFORMATION_EX structure,
// and accumulates counts based on the relationship type.
//
// For RelationGroup, it carefully parses the GROUP_RELATIONSHIP struct, then iterates over
// the subsequent PROCESSOR_GROUP_INFO array to count logical and active cores.
//
// Returns an error if the Windows API calls fail.
func GetNCores(cpu *CPUResult) error {
	var bufLen uint32

	// First call with NULL buffer to get required buffer size
	_, _, err := procGetLogicalProcessorInformationEx.Call(
		uintptr(RelationAll),
		uintptr(0),
		uintptr(unsafe.Pointer(&bufLen)),
	)

	if bufLen == 0 {
		return err
	}

	buffer := make([]byte, bufLen)

	// Second call to get the actual data
	ret, _, err := procGetLogicalProcessorInformationEx.Call(
		uintptr(RelationAll),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&bufLen)),
	)

	if ret == 0 {
		return err
	}

	offset := uintptr(0)

	for offset < uintptr(bufLen) {
		ptr := unsafe.Pointer(&buffer[offset])
		info := (*SYSTEM_LOGICAL_PROCESSOR_INFORMATION_EX)(ptr)

		if info.Relationship == RelationGroup {
			// ASCII diagram illustrating the memory layout:
			//
			// |---------------------------| <- uintptr(unsafe.Pointer(group))
			// | GROUP_RELATIONSHIP struct |
			// |---------------------------| <- sizeof(GROUP_RELATIONSHIP)
			// | PROCESSOR_GROUP_INFO[0]   | <- groupInfoBase
			// | PROCESSOR_GROUP_INFO[1]   |
			// | ...                       |
			//
			// So, PROCESSOR_GROUP_INFO array starts immediately after GROUP_RELATIONSHIP in memory.
			//

			group := (*GROUP_RELATIONSHIP)(unsafe.Pointer(&info.Data[0]))
			groupCount := int(group.ActiveGroupCount)

			groupInfoBase := uintptr(unsafe.Pointer(group)) + unsafe.Sizeof(*group)

			for i := 0; i < groupCount; i++ {
				groupInfo := (*PROCESSOR_GROUP_INFO)(unsafe.Pointer(groupInfoBase + uintptr(i)*unsafe.Sizeof(PROCESSOR_GROUP_INFO{})))

				cpu.CoresActive += uint16(groupInfo.ActiveProcessorCount)
				cpu.CoresLogical += uint16(groupInfo.MaximumProcessorCount)
			}

		} else if info.Relationship == RelationProcessorCore {
			// Each core counts as one physical core
			cpu.CoresPhysical++
		} else if info.Relationship == RelationProcessorPackage {
			// Each package corresponds to one physical CPU socket
			cpu.Packages++
		}

		offset += uintptr(info.Size)
	}
	return nil
}

func GetRegistryData(cpu *CPUResult) error {
	var key windows.Handle

	sKey, err := windows.UTF16PtrFromString(`HARDWARE\DESCRIPTION\System\CentralProcessor\0`)

	if err != nil {
		return err
	}

	err = windows.RegOpenKeyEx(
		windows.HKEY_LOCAL_MACHINE,
		sKey,
		0,
		windows.KEY_READ|windows.KEY_WOW64_64KEY,
		&key,
	)

	if err != nil {
		return err
	}

	defer windows.RegCloseKey(key)

	processName, err := registry.ReadRegSZ(key, `ProcessorNameString`)

	if err != nil {
		return err
	}

	cpu.ProcessorName = windows.UTF16ToString(processName)

	vendor, err := registry.ReadRegSZ(key, `VendorIdentifier`)

	if err != nil {
		return err
	}

	cpu.Vendor = windows.UTF16ToString(vendor)

	identifier, err := registry.ReadRegSZ(key, `Identifier`)

	if err != nil {
		return err
	}

	cpu.Identifier = windows.UTF16ToString(identifier)

	mhz, err := registry.ReadRegDWORD(key, `~MHz`)

	if err != nil {
		return err
	}

	cpu.Mhz = mhz

	return nil
}
