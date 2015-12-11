// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"

func (dev *Device) ControlTransfer(
	requestType uint8,
	request uint8,
	value uint16,
	index uint16,
	data []byte,
	length uint16,
	timeout uint,
) (int, error) {
	return 0, nil
}
