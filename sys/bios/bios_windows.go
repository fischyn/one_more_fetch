//go:build windows

package bios

import (
	"context"
	"unsafe"

	"golang.org/x/sys/windows"
)

func GetBiosInfo(_ context.Context) (product, manufacturer, family, version string, err error) {

	var handler windows.Handle

	err = windows.RegOpenKeyEx(
		windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(`HARDWARE\DESCRIPTION\System\BIOS`),
		0,
		windows.KEY_READ|windows.KEY_WOW64_64KEY,
		&handler,
	)

	defer windows.RegCloseKey(handler)

	if err != nil {
		return
	}

	var bufLen uint32
	var valType uint32

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

	regBuf := make([]uint16, bufLen/2+1)

	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`SystemProductName`),
		nil,
		&valType,
		(*byte)(unsafe.Pointer(&regBuf[0])),
		&bufLen,
	)
	if err != nil {
		return
	}

	product = windows.UTF16ToString(regBuf)

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

	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`SystemManufacturer`),
		nil,
		&valType,
		(*byte)(unsafe.Pointer(&regBuf[0])),
		&bufLen,
	)
	if err != nil {
		return
	}

	manufacturer = windows.UTF16ToString(regBuf)

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

	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`SystemFamily`),
		nil,
		&valType,
		(*byte)(unsafe.Pointer(&regBuf[0])),
		&bufLen,
	)
	if err != nil {
		return
	}

	family = windows.UTF16ToString(regBuf)

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

	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`BIOSVersion`),
		nil,
		&valType,
		(*byte)(unsafe.Pointer(&regBuf[0])),
		&bufLen,
	)
	if err != nil {
		return
	}

	version = windows.UTF16ToString(regBuf)

	return product, manufacturer, family, version, nil
}
