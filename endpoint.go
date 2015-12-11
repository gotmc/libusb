// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"

type Endpoint struct {
}

type EndpointDescriptor struct {
	Length          int
	DescriptorType  descriptorType
	EndpointAddress uint8
	Attributes      uint8
	MaxPacketSize   uint16
	Interval        uint8
	Refresh         uint8
	SynchAddress    uint8
}

type EndpointDescriptors []*EndpointDescriptor
