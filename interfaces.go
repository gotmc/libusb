// Copyright (c) 2015-2017 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"

// SupportedInterface models an supported USB interface and its associated
// interface descriptors.
type SupportedInterface struct {
	InterfaceDescriptors
	NumAltSettings int
}

// SupportedInterfaces contains an array of the supported USB interfaces for a
// given USB device.
type SupportedInterfaces []*SupportedInterface

// InterfaceDescriptor "provides information about a function or feature that a
// device implements." (Source: *USB Complete* 5th edition by Jan Axelson)
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

// InterfaceDescriptors contains a slice of pointers to the available interface
// descriptors.
type InterfaceDescriptors []*InterfaceDescriptor
