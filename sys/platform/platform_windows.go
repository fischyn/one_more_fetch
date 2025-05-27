//go:build windows

package platform

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fischyn/wfetch/sys/registry"
	"golang.org/x/sys/windows"
)

func GetPlatform(pLatform *PlatformResult) error {
	var key windows.Handle

	err := windows.RegOpenKeyEx(
		windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(`SOFTWARE\Microsoft\Windows NT\CurrentVersion`),
		0,
		windows.KEY_READ|windows.KEY_WOW64_64KEY,
		&key,
	)

	if err != nil {
		return err
	}

	defer windows.RegCloseKey(key)

	productName, err := registry.ReadRegSZ(key, `ProductName`)

	if err != nil {
		return err
	}

	pLatform.Platform = windows.UTF16ToString(productName)

	if strings.Contains(pLatform.Platform, "Windows 10") {
		buildNumber, err := registry.ReadRegSZ(key, `CurrentBuildNumber`)

		if err == nil {

			if buildNumber, err := strconv.ParseInt(
				windows.UTF16ToString(buildNumber),
				10,
				32,
			); err == nil && buildNumber >= 22000 {
				pLatform.Platform = strings.Replace(pLatform.Platform, "Windows 10", "Windows 11", 1)
			}
		}

		if !strings.HasPrefix(pLatform.Platform, "Microsoft") {
			pLatform.Platform = "Microsoft " + pLatform.Platform
		}

		csdVersion, err := registry.ReadRegDWORD(key, `CSDVersion`)

		if err == nil {
			pLatform.Platform += " " + strconv.FormatUint(uint64(csdVersion), 10)
		}
	}

	displayVersion, err := registry.ReadRegSZ(key, `DisplayVersion`)

	if err != nil {
		return err
	}

	pLatform.DisplayVersion = windows.UTF16ToString(displayVersion)

	versionInfo := windows.RtlGetVersion()

	UBR, err := registry.ReadRegDWORD(key, `UBR`)

	if err != nil {
		return err
	}

	pLatform.Version = fmt.Sprintf("%d.%d.%d.%d Build %d.%d",
		versionInfo.MajorVersion,
		versionInfo.MinorVersion,
		versionInfo.BuildNumber,
		UBR,
		versionInfo.BuildNumber,
		UBR,
	)

	switch versionInfo.ProductType {
	case 1:
		pLatform.Family = "Standalone Workstation"
	case 2:
		pLatform.Family = "Server (Domain Controller)"
	case 3:
		pLatform.Family = "Server"
	}

	return nil
}
