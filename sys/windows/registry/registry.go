package registry

import (
	"encoding/binary"
	"fmt"
	"unsafe"

	"golang.org/x/sys/windows"
)

func ReadRegSZ(key windows.Handle, name string) ([]uint16, error) {
	namePtr, err := windows.UTF16PtrFromString(name)

	if err != nil {
		return nil, err
	}

	var bufLen uint32
	var valType uint32

	err = windows.RegQueryValueEx(
		key,
		namePtr,
		nil,
		&valType,
		nil,
		&bufLen,
	)

	if err != nil {
		return nil, err
	}

	if valType != windows.REG_SZ {
		return nil, fmt.Errorf("unexpected registry value type for %s: %d", name, valType)
	}

	if bufLen < 2 {
		return nil, fmt.Errorf("buffer too small for REG_SZ")
	}

	buf := make([]uint16, bufLen)

	err = windows.RegQueryValueEx(
		key,
		namePtr,
		nil,
		&valType,
		(*byte)(unsafe.Pointer(&buf[0])),
		&bufLen,
	)

	if err != nil {
		return nil, err
	}

	return buf, nil
}

func ReadRegDWORD(key windows.Handle, name string) (uint32, error) {
	namePtr, err := windows.UTF16PtrFromString(name)

	if err != nil {
		return 0, err
	}

	var bufLen uint32
	var valType uint32

	err = windows.RegQueryValueEx(
		key,
		namePtr,
		nil,
		&valType,
		nil,
		&bufLen,
	)

	if err != nil {
		return 0, err
	}

	if valType != windows.REG_DWORD {
		return 0, fmt.Errorf("unexpected registry value type for %s: %d", name, valType)
	}

	if bufLen < 4 {
		return 0, fmt.Errorf("buffer too small for REG_DWORD")
	}

	buf := make([]byte, bufLen)

	err = windows.RegQueryValueEx(
		key,
		namePtr,
		nil,
		&valType,
		&buf[0],
		&bufLen,
	)

	if err != nil {
		return 0, err
	}

	val := binary.LittleEndian.Uint32(buf[:4])

	return val, nil
}
