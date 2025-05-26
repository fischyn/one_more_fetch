package gpu

import (
	"unsafe"

	"github.com/fischyn/omfetch/sys/windows/dll"

	// "github.com/fischyn/omfetch/sys/windows/registry"
	"golang.org/x/sys/windows"
)

var (
	procSetupDiGetClassDevsW             = dll.SetupApi.NewProc("SetupDiGetClassDevsW")
	procSetupDiEnumDeviceInfo            = dll.SetupApi.NewProc("SetupDiEnumDeviceInfo")
	procSetupDiGetDeviceRegistryProperty = dll.SetupApi.NewProc("SetupDiGetDeviceRegistryPropertyW")
	procSetupDiOpenDevRegKey             = dll.SetupApi.NewProc("SetupDiOpenDevRegKey")
)

const (
	DICS_FLAG_GLOBAL = 0x00000001
	DIGCF_PRESENT    = 0x00000002
	DIGCF_ALLCLASSES = 0x00000004

	SPDRP_DEVICEDESC = 0x00000000

	DIREG_DEV = 0x00000001
	KEY_READ  = 0x20019

	ERROR_NO_MORE_ITEMS = 259
)

var GUID_DEVCLASS_DISPLAY = windows.GUID{
	Data1: 0x4d36e968,
	Data2: 0xe325,
	Data3: 0x11ce,
	Data4: [8]byte{0xbf, 0xc1, 0x08, 0x00, 0x2b, 0xe1, 0x03, 0x18},
}

type SP_DEVINFO_DATA struct {
	cbSize    uint32
	ClassGuid windows.GUID
	DevInst   uint32
	Reserved  uintptr
}

func GetGPUsInfo(gpus *[]GPUResult) error {
	devInfoSet, err := getDisplayDeviceInfoSet()
	if err != nil {
		return err
	}
	defer windows.SetupDiDestroyDeviceInfoList(windows.DevInfo(devInfoSet))

	devices, err := setupDiEnumDeviceInfo(devInfoSet)
	if err != nil {
		return err
	}

	for _, dev := range devices {
		name, _ := setupDiGetDeviceRegistryProperty(devInfoSet, &dev, SPDRP_DEVICEDESC)

		regKey, err := setupDiOpenDevRegKey(devInfoSet, &dev)
		if err != nil {
			continue
		}

		defer windows.RegCloseKey(regKey)

		// version, _ := registry.ReadRegSZ(regKey, `VideoID`)

		//TODO read from "SOFTWARE\Microsoft\DirectX\" regisry
		//TODO read from "SYSTEM\CurrentControlSet\Control\Class\"

		info := GPUResult{
			Name: name,
		}
		*gpus = append(*gpus, info)
	}

	return nil
}

func getDisplayDeviceInfoSet() (windows.Handle, error) {

	handle, _, err := procSetupDiGetClassDevsW.Call(
		uintptr(unsafe.Pointer(&GUID_DEVCLASS_DISPLAY)),
		0, // Enumerator
		0, // hwndParent
		uintptr(DIGCF_PRESENT),
	)

	if handle == 0 || handle == ^uintptr(0) {
		return 0, err
	}

	return windows.Handle(handle), nil
}

func setupDiEnumDeviceInfo(devInfoSet windows.Handle) ([]SP_DEVINFO_DATA, error) {
	var devices []SP_DEVINFO_DATA
	var index uint32 = 0

	for {
		var devInfo SP_DEVINFO_DATA
		devInfo.cbSize = uint32(unsafe.Sizeof(devInfo))

		ret, _, err := procSetupDiEnumDeviceInfo.Call(
			uintptr(devInfoSet),
			uintptr(index),
			uintptr(unsafe.Pointer(&devInfo)),
		)
		if ret == 0 {
			if err == windows.Errno(ERROR_NO_MORE_ITEMS) {
				break
			}
			return nil, err
		}

		devices = append(devices, devInfo)
		index++
	}

	return devices, nil
}

func setupDiGetDeviceRegistryProperty(devInfoSet windows.Handle, devInfo *SP_DEVINFO_DATA, prop uint32) (string, error) {
	var dataType uint32
	var requiredSize uint32

	ret, _, err := procSetupDiGetDeviceRegistryProperty.Call(
		uintptr(devInfoSet),
		uintptr(unsafe.Pointer(devInfo)),
		uintptr(prop),
		uintptr(unsafe.Pointer(&dataType)),
		0,
		0,
		uintptr(unsafe.Pointer(&requiredSize)),
	)
	if ret == 0 && requiredSize == 0 {
		return "", err
	}

	buf := make([]uint16, requiredSize/2)

	ret, _, err = procSetupDiGetDeviceRegistryProperty.Call(
		uintptr(devInfoSet),
		uintptr(unsafe.Pointer(devInfo)),
		uintptr(prop),
		uintptr(unsafe.Pointer(&dataType)),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(requiredSize),
		uintptr(unsafe.Pointer(&requiredSize)),
	)

	if ret == 0 {
		return "", err
	}

	return windows.UTF16ToString(buf), nil
}

func setupDiOpenDevRegKey(devInfoSet windows.Handle, devInfo *SP_DEVINFO_DATA) (windows.Handle, error) {
	regKey, _, err := procSetupDiOpenDevRegKey.Call(
		uintptr(devInfoSet),
		uintptr(unsafe.Pointer(devInfo)),
		uintptr(DICS_FLAG_GLOBAL),
		0,
		uintptr(DIREG_DEV),
		uintptr(KEY_READ),
	)
	if regKey == 0 || regKey == ^uintptr(0) {
		return 0, err
	}
	return windows.Handle(regKey), nil
}
