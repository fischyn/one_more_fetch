//go:build windows

package cpu

import (
	"unsafe"

	"github.com/fischyn/omfetch/sys/windows/dll"
)

// https://learn.microsoft.com/ru-ru/windows/win32/api/winnt/ne-winnt-logical_processor_relationship
type LOGICAL_PROCESSOR_RELATIONSHIP = uint16

// Constants defining types of processor relationships.
// These correspond to the Windows API GetLogicalProcessorInformationEx constants.
const (
	RelationProcessorCore    LOGICAL_PROCESSOR_RELATIONSHIP = 0
	RelationProcessorPackage LOGICAL_PROCESSOR_RELATIONSHIP = 3
	RelationGroup            LOGICAL_PROCESSOR_RELATIONSHIP = 4
	RelationAll              LOGICAL_PROCESSOR_RELATIONSHIP = 0xffff
)

// SYSTEM_LOGICAL_PROCESSOR_INFORMATION_EX represents the Windows SYSTEM_LOGICAL_PROCESSOR_INFORMATION_EX structure.
//
// This structure describes information about logical processors and their relationships.
//
// Fields:
// - Relationship: The type of relationship (core, package, group, etc.).
// - Size: Size in bytes of the entire structure including the variable-sized Data field.
// - Data: Placeholder for variable-sized data. The actual content depends on Relationship type.
type SYSTEM_LOGICAL_PROCESSOR_INFORMATION_EX struct {
	Relationship uint16
	Size         uint32
	Data         [1]byte // super-duper hack?
}

// GROUP_RELATIONSHIP represents processor group information.
//
// According to MSDN, this structure is followed immediately in memory
// by an array of PROCESSOR_GROUP_INFO structs, describing each group.
type GROUP_RELATIONSHIP struct {
	MaximumGroupCount uint16
	ActiveGroupCount  uint16
	Reserved          [20]byte
	// Note: In C, following this struct directly in memory is an array of PROCESSOR_GROUP_INFO structs.
}

// PROCESSOR_GROUP_INFO represents information about one processor group.
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
