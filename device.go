// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"

type devices struct {
	Devices **C.libusb_device
}

type device struct {
	device *C.libusb_device
}

func (ctx *context) GetDeviceList() (*devices, error) {
	var list **C.libusb_device
	const unrefDevices = 1
	numDevicesFound := C.libusb_get_device_list(ctx.context, &list)
	if numDevicesFound < 0 {
		return nil, libusbError(numDevicesFound)
	}
	defer C.libusb_free_device_list(list, unrefDevices)
	devices := &devices{
		Devices: list,
	}
	return devices, nil
}
