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

type endpointAddress byte
type endpointAttributes byte

type EndpointDescriptor struct {
	Length          int
	DescriptorType  descriptorType
	EndpointAddress endpointAddress
	Attributes      endpointAttributes
	MaxPacketSize   uint16
	Interval        uint8
	Refresh         uint8
	SynchAddress    uint8
}

type EndpointDescriptors []*EndpointDescriptor

// Direction returns the endpointDirection.
func (end *EndpointDescriptor) Direction() endpointDirection {
	return end.EndpointAddress.direction()
}

// Number returns the endpoint number in bits 0..3 in the endpoint
// address.
func (end *EndpointDescriptor) Number() byte {
	return end.EndpointAddress.endpointNumber()
}
func (end *EndpointDescriptor) TransferType() transferType {
	return end.Attributes.transferType()
}

func (address endpointAddress) direction() endpointDirection {
	// Bit 7 of the endpointAddress determines the direction
	const directionMask = 0x80
	const directionBit = 7
	return endpointDirection(address&directionMask) >> directionBit
}

func (address endpointAddress) endpointNumber() byte {
	// Bits 0..3 determine the endpoint number
	const endpointNumberMask = 0x0F
	return byte(address & endpointNumberMask)
}

func (attributes endpointAttributes) transferType() transferType {
	// Bits 0..1 of the bmAttributes determines the transfer type
	const transferTypeMask = 0x03
	return transferType(attributes & transferTypeMask)
}
