// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"

type deviceHandle struct {
	libusbDeviceHandle *C.libusb_device_handle
}

func (devHandle *deviceHandle) Close() error {
	C.libusb_close(devHandle.libusbDeviceHandle)
	return nil
}
