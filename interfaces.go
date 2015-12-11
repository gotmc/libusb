// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.
package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"

type SupportedInterface struct {
	InterfaceDescriptors
	NumAltSettings int
}

type SupportedInterfaces []*SupportedInterface

type InterfaceDescriptor struct {
	Length            int
	DescriptorType    descriptorType
	InterfaceNumber   int
	AlternateSetting  int
	NumEndpoints      int
	InterfaceClass    uint8
	InterfaceSubClass uint8
	InterfaceProtocol uint8
	InterfaceIndex    int
	EndpointDescriptors
}

type InterfaceDescriptors []*InterfaceDescriptor
