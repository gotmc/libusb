// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"
import (
	"reflect"
	"unsafe"
)

type libusbSpeed int

const (
	speedUnknown libusbSpeed = C.LIBUSB_SPEED_UNKNOWN
	speedLow     libusbSpeed = C.LIBUSB_SPEED_LOW
	speedFull    libusbSpeed = C.LIBUSB_SPEED_FULL
	speedHigh    libusbSpeed = C.LIBUSB_SPEED_HIGH
	speedSuper   libusbSpeed = C.LIBUSB_SPEED_SUPER
)

var speedCodes = map[libusbSpeed]string{
	speedUnknown: "The OS doesn't report or know the device speed.",
	speedLow:     "The device is operating at low speed (1.5MBit/s)",
	speedFull:    "The device is operating at full speed (12MBit/s)",
	speedHigh:    "The device is operating at high speed (480MBit/s)",
	speedSuper:   "The device is operating at super speed (5000MBit/s)",
}

func (speed libusbSpeed) String() string {
	return speedCodes[speed]
}

type device struct {
	libusbDevice *C.libusb_device
}

func (ctx *context) GetDeviceList() ([]*device, error) {
	var devices []*device
	var list **C.libusb_device
	const unrefDevices = 1
	numDevicesFound := int(C.libusb_get_device_list(ctx.context, &list))
	if numDevicesFound < 0 {
		return nil, ErrorCode(numDevicesFound)
	}
	defer C.libusb_free_device_list(list, unrefDevices)
	var libusbDevices []*C.libusb_device
	*(*reflect.SliceHeader)(unsafe.Pointer(&libusbDevices)) = reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(list)),
		Len:  numDevicesFound,
		Cap:  numDevicesFound,
	}
	for _, thisLibusbDevice := range libusbDevices {
		thisDevice := device{
			libusbDevice: thisLibusbDevice,
		}
		devices = append(devices, &thisDevice)
	}
	return devices, nil
}

func (dev *device) GetBusNumber() (uint, error) {
	busNumber, err := C.libusb_get_bus_number(dev.libusbDevice)
	if err != nil {
		return 0, err
	}
	return uint(busNumber), nil
}

func (dev *device) GetDeviceAddress() (uint, error) {
	deviceAddress, err := C.libusb_get_device_address(dev.libusbDevice)
	if err != nil {
		return 0, err
	}
	return uint(deviceAddress), nil
}

func (dev *device) GetDeviceSpeed() (libusbSpeed, error) {
	deviceSpeed, err := C.libusb_get_device_speed(dev.libusbDevice)
	if err != nil {
		return 0, err
	}
	return libusbSpeed(deviceSpeed), nil
}
