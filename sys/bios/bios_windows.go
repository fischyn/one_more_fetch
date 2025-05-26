//go:build windows

package bios

import (
	"github.com/fischyn/omfetch/sys/windows/registry"
	"golang.org/x/sys/windows"
)

func GetBIOS(bIOS *BIOSResult, opt BIOSOptions) error {
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

	if opt.ShowSystemProduct {
		systemProduct, err := registry.ReadRegSZ(key, `SystemProductName`)
		if err != nil {
			return err
		}
		bIOS.SystemProduct = windows.UTF16ToString(systemProduct)
	}

	if opt.ShowSystemManufacturer {
		systemManufacturer, err := registry.ReadRegSZ(key, `SystemManufacturer`)
		if err != nil {
			return err
		}
		bIOS.SystemManufacturer = windows.UTF16ToString(systemManufacturer)
	}

	if opt.ShowSystemFamily {
		systemFamily, err := registry.ReadRegSZ(key, `SystemFamily`)
		if err != nil {
			return err
		}
		bIOS.SystemFamily = windows.UTF16ToString(systemFamily)
	}

	if opt.ShowSystemVersion {
		systemVersion, err := registry.ReadRegSZ(key, `SystemVersion`)
		if err != nil {
			return err
		}
		bIOS.SystemVersion = windows.UTF16ToString(systemVersion)
	}

	if opt.ShowBiosVendor {
		biosVendor, err := registry.ReadRegSZ(key, `BIOSVendor`)
		if err != nil {
			return err
		}
		bIOS.BiosVendor = windows.UTF16ToString(biosVendor)
	}

	if opt.ShowBiosVersion {
		biosVersion, err := registry.ReadRegSZ(key, `BIOSVersion`)
		if err != nil {
			return err
		}
		bIOS.BiosVersion = windows.UTF16ToString(biosVersion)
	}

	if opt.ShowBiosReleaseDate {
		biosReleaseDate, err := registry.ReadRegSZ(key, `BIOSReleaseDate`)
		if err != nil {
			return err
		}
		bIOS.BiosReleaseDate = windows.UTF16ToString(biosReleaseDate)
	}

	if opt.ShowBaseBoardManufacturer {
		baseBoardManufacturer, err := registry.ReadRegSZ(key, `BaseBoardManufacturer`)
		if err != nil {
			return err
		}
		bIOS.BaseBoardManufacturer = windows.UTF16ToString(baseBoardManufacturer)
	}

	if opt.ShowBaseBoardProduct {
		baseBoardProduct, err := registry.ReadRegSZ(key, `BaseBoardProduct`)
		if err != nil {
			return err
		}
		bIOS.BaseBoardProduct = windows.UTF16ToString(baseBoardProduct)
	}

	if opt.ShowBaseBoardVersion {
		baseBoardVersion, err := registry.ReadRegSZ(key, `BaseBoardVersion`)
		if err != nil {
			return err
		}
		bIOS.BaseBoardVersion = windows.UTF16ToString(baseBoardVersion)
	}

	return nil
}
