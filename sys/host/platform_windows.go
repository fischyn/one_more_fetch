//go:build windows

package host

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"unsafe"

	"golang.org/x/sys/windows"
)

func GetPlatformInfo(_ context.Context) (platform string, family string, version string, displayVersion string, err error) {
	versionInfo := windows.RtlGetVersion()

	var handler windows.Handle

	err = windows.RegOpenKeyEx(
		windows.HKEY_LOCAL_MACHINE,
		windows.StringToUTF16Ptr(`SOFTWARE\Microsoft\Windows NT\CurrentVersion`),
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

	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`ProductName`),
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
		windows.StringToUTF16Ptr(`ProductName`),
		nil,
		&valType,
		(*byte)(unsafe.Pointer(&regBuf[0])),
		&bufLen,
	)
	if err != nil {
		return
	}
	platform = windows.UTF16ToString(regBuf)

	if strings.Contains(platform, "Windows 10") {
		err = windows.RegQueryValueEx(
			handler,
			windows.StringToUTF16Ptr(`CurrentBuildNumber`),
			nil,
			&valType,
			nil,
			&bufLen,
		)
		if err == nil {
			regBuf = make([]uint16, bufLen/2+1)
			err = windows.RegQueryValueEx(
				handler,
				windows.StringToUTF16Ptr(`CurrentBuildNumber`),
				nil,
				&valType,
				(*byte)(unsafe.Pointer(&regBuf[0])),
				&bufLen,
			)
			if err == nil {
				buildNumberStr := windows.UTF16ToString(regBuf)
				if buildNumber, err := strconv.ParseInt(buildNumberStr, 10, 32); err == nil && buildNumber >= 22000 {
					platform = strings.Replace(platform, "Windows 10", "Windows 11", 1)
				}
			}
		}
	}
	if !strings.HasPrefix(platform, "Microsoft") {
		platform = "Microsoft " + platform
	}

	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`CSDVersion`),
		nil,
		&valType,
		nil,
		&bufLen,
	)

	if err == nil {
		regBuf = make([]uint16, bufLen/2+1)
		err = windows.RegQueryValueEx(
			handler,
			windows.StringToUTF16Ptr(`CSDVersion`),
			nil,
			&valType,
			(*byte)(unsafe.Pointer(&regBuf[0])),
			&bufLen,
		)
		if err == nil {
			platform += " " + windows.UTF16ToString(regBuf)
		}
	}

	var UBR uint32
	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`UBR`),
		nil,
		&valType,
		nil,
		&bufLen,
	)
	if err == nil {
		regBuf := make([]byte, 4)
		err = windows.RegQueryValueEx(
			handler,
			windows.StringToUTF16Ptr(`UBR`),
			nil,
			&valType,
			(*byte)(unsafe.Pointer(&regBuf[0])),
			&bufLen,
		)

		copy((*[4]byte)(unsafe.Pointer(&UBR))[:], regBuf)
	}

	err = windows.RegQueryValueEx(
		handler,
		windows.StringToUTF16Ptr(`DisplayVersion`),
		nil,
		&valType,
		nil,
		&bufLen,
	)
	if err == nil {
		regBuf := make([]uint16, bufLen/2+1)
		err = windows.RegQueryValueEx(
			handler,
			windows.StringToUTF16Ptr(`DisplayVersion`),
			nil,
			&valType,
			(*byte)(unsafe.Pointer(&regBuf[0])),
			&bufLen,
		)

		displayVersion = windows.UTF16ToString(regBuf)
	}

	version = fmt.Sprintf("%d.%d.%d.%d Build %d.%d",
		versionInfo.MajorVersion,
		versionInfo.MinorVersion,
		versionInfo.BuildNumber,
		UBR,
		versionInfo.BuildNumber,
		UBR,
	)

	switch versionInfo.ProductType {
	case 1:
		family = "Standalone Workstation"
	case 2:
		family = "Server (Domain Controller)"
	case 3:
		family = "Server"
	}

	return platform, family, version, displayVersion, nil
}
