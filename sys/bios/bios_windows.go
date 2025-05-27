//go:build windows

package bios

import (
	"github.com/fischyn/wfetch/sys/registry"
	"golang.org/x/sys/windows"
)

// TODO add  opt BIOSOptions
func GetBIOS(bIOS *BIOSResult) error {
	var key windows.Handle

	err := windows.RegOpenKeyEx(
		windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(`HARDWARE\DESCRIPTION\System\BIOS`),
		0,
		windows.KEY_READ|windows.KEY_WOW64_64KEY,
		&key,
	)

	if err != nil {
		return err
	}

	defer windows.RegCloseKey(key)

	systemProduct, err := registry.ReadRegSZ(key, `SystemProductName`)
	if err != nil {
		return err
	}
	bIOS.SystemProduct = windows.UTF16ToString(systemProduct)

	systemManufacturer, err := registry.ReadRegSZ(key, `SystemManufacturer`)
	if err != nil {
		return err
	}
	bIOS.SystemManufacturer = windows.UTF16ToString(systemManufacturer)

	systemFamily, err := registry.ReadRegSZ(key, `SystemFamily`)
	if err != nil {
		return err
	}
	bIOS.SystemFamily = windows.UTF16ToString(systemFamily)

	systemVersion, err := registry.ReadRegSZ(key, `SystemVersion`)
	if err != nil {
		return err
	}
	bIOS.SystemVersion = windows.UTF16ToString(systemVersion)

	biosVendor, err := registry.ReadRegSZ(key, `BIOSVendor`)
	if err != nil {
		return err
	}
	bIOS.BiosVendor = windows.UTF16ToString(biosVendor)

	biosVersion, err := registry.ReadRegSZ(key, `BIOSVersion`)
	if err != nil {
		return err
	}
	bIOS.BiosVersion = windows.UTF16ToString(biosVersion)

	biosReleaseDate, err := registry.ReadRegSZ(key, `BIOSReleaseDate`)
	if err != nil {
		return err
	}
	bIOS.BiosReleaseDate = windows.UTF16ToString(biosReleaseDate)

	baseBoardManufacturer, err := registry.ReadRegSZ(key, `BaseBoardManufacturer`)
	if err != nil {
		return err
	}
	bIOS.BaseBoardManufacturer = windows.UTF16ToString(baseBoardManufacturer)

	baseBoardProduct, err := registry.ReadRegSZ(key, `BaseBoardProduct`)
	if err != nil {
		return err
	}
	bIOS.BaseBoardProduct = windows.UTF16ToString(baseBoardProduct)

	baseBoardVersion, err := registry.ReadRegSZ(key, `BaseBoardVersion`)
	if err != nil {
		return err
	}
	bIOS.BaseBoardVersion = windows.UTF16ToString(baseBoardVersion)

	return nil
}
