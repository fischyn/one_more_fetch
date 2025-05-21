//go:build windows

package bios

import (
	"unsafe"

	"golang.org/x/sys/windows"
)

func GetBiosInfo() (product, manufacturer, family, version string, err error) {
	var handler windows.Handle

	err = windows.RegOpenKeyEx(
		windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(`HARDWARE\DESCRIPTION\System\BIOS`),
		0,
		windows.KEY_READ|windows.KEY_WOW64_64KEY,
		&handler,
	)
	if err != nil {
		return
	}
	defer windows.RegCloseKey(handler)

	var bufLen uint32
	var valType uint32

	// SystemProductName
	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`SystemProductName`),
		nil,
		&valType,
		nil,
		&bufLen,
	)
	if err != nil {
		return
	}

	productBuf := make([]uint16, bufLen/2+1)
	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`SystemProductName`),
		nil,
		&valType,
		(*byte)(unsafe.Pointer(&productBuf[0])),
		&bufLen,
	)
	if err != nil {
		return
	}
	product = windows.UTF16ToString(productBuf)

	// SystemManufacturer
	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`SystemManufacturer`),
		nil,
		&valType,
		nil,
		&bufLen,
	)
	if err != nil {
		return
	}

	manufacturerBuf := make([]uint16, bufLen/2+1)
	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`SystemManufacturer`),
		nil,
		&valType,
		(*byte)(unsafe.Pointer(&manufacturerBuf[0])),
		&bufLen,
	)
	if err != nil {
		return
	}
	manufacturer = windows.UTF16ToString(manufacturerBuf)

	// SystemFamily
	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`SystemFamily`),
		nil,
		&valType,
		nil,
		&bufLen,
	)
	if err != nil {
		return
	}

	familyBuf := make([]uint16, bufLen/2+1)
	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`SystemFamily`),
		nil,
		&valType,
		(*byte)(unsafe.Pointer(&familyBuf[0])),
		&bufLen,
	)
	if err != nil {
		return
	}
	family = windows.UTF16ToString(familyBuf)

	// BIOSVersion
	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`BIOSVersion`),
		nil,
		&valType,
		nil,
		&bufLen,
	)
	if err != nil {
		return
	}

	versionBuf := make([]uint16, bufLen/2+1)
	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`BIOSVersion`),
		nil,
		&valType,
		(*byte)(unsafe.Pointer(&versionBuf[0])),
		&bufLen,
	)
	if err != nil {
		return
	}
	version = windows.UTF16ToString(versionBuf)

	return product, manufacturer, family, version, nil
}
