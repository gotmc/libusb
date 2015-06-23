// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"
import "unsafe"

// DeviceHandle represents the libusb device handle.
type DeviceHandle struct {
	libusbDeviceHandle *C.libusb_device_handle
}

// Close closes the device handle.
func (dev *DeviceHandle) Close() error {
	C.libusb_close(dev.libusbDeviceHandle)
	return nil
}

func (dev *DeviceHandle) GetStringDescriptor(
	descIndex uint8,
	langID uint16,
) (string, error) {
	var data *C.uchar
	length := 512
	usberr := C.libusb_get_string_descriptor(
		dev.libusbDeviceHandle,
		C.uint8_t(descIndex),
		C.uint16_t(langID),
		data,
		C.int(length),
	)
	if usberr < 0 {
		return "", ErrorCode(usberr)
	}
	return "Yes!!!", nil
}

func (dev *DeviceHandle) GetStringDescriptorASCII(descIndex uint8) (string, error) {
	length := 256
	data := make([]byte, length)
	usberr := C.libusb_get_string_descriptor_ascii(
		dev.libusbDeviceHandle,
		C.uint8_t(descIndex),
		// Unsafe pointer -> http://stackoverflow.com/a/16376039/95592
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.int(length),
	)
	if usberr < 0 {
		return "", ErrorCode(usberr)
	}
	return string(data), nil
}
