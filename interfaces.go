// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
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

// USB Interface Class Codes
const (
	// USB Interface Class Codes - useful for finding interfaces by class
	InterfaceClassAudio              uint8 = 0x01
	InterfaceClassComm               uint8 = 0x02
	InterfaceClassHID                uint8 = 0x03
	InterfaceClassPhysical           uint8 = 0x05
	InterfaceClassImage              uint8 = 0x06
	InterfaceClassPrinter            uint8 = 0x07
	InterfaceClassMassStorage        uint8 = 0x08
	InterfaceClassHub                uint8 = 0x09
	InterfaceClassData               uint8 = 0x0A
	InterfaceClassSmartCard          uint8 = 0x0B
	InterfaceClassContentSecurity    uint8 = 0x0D
	InterfaceClassVideo              uint8 = 0x0E
	InterfaceClassPersonalHealthcare uint8 = 0x0F
	InterfaceClassAudioVideo         uint8 = 0x10
	InterfaceClassWireless           uint8 = 0xE0
	InterfaceClassApplication        uint8 = 0xFE
	InterfaceClassVendorSpec         uint8 = 0xFF
)

// GetAllInterfacesByClass searches through all supported interfaces and their interface
// descriptors to find all interfaces matching the given class code.
// This is especially useful for printer devices (InterfaceClassPrinter) and other specialized devices.
func (si SupportedInterfaces) GetAllInterfacesByClass(class uint8) InterfaceDescriptors {
	var matchingInterfaces InterfaceDescriptors
	for _, supportedIface := range si {
		for _, ifaceDesc := range supportedIface.InterfaceDescriptors {
			if ifaceDesc.InterfaceClass == class {
				matchingInterfaces = append(matchingInterfaces, ifaceDesc)
			}
		}
	}
	return matchingInterfaces
}

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

// GetInterfacesByClass returns all interface descriptors that match the given USB class code.
// This is particularly useful for devices that report bDeviceClass as LIBUSB_CLASS_PER_INTERFACE
// (0), where individual interfaces have their own class codes.
func (ifaces InterfaceDescriptors) GetInterfacesByClass(class uint8) InterfaceDescriptors {
	var matchingInterfaces InterfaceDescriptors
	for _, iface := range ifaces {
		if iface.InterfaceClass == class {
			matchingInterfaces = append(matchingInterfaces, iface)
		}
	}
	return matchingInterfaces
}
