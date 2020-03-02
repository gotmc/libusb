// Copyright (c) 2015-2020 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"

// Endpoint doesn't seem to model anything. Did I replace this with
// EndpointDescriptor?
type Endpoint struct {
	// FIXME(mdr): Is this needed/used? Can this safely be deleted?
}

type endpointAddress byte
type endpointAttributes byte

// EndpointDescriptor models the descriptor for a given endpoint.
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

// EndpointDescriptors contains the available endpoint descriptors.
type EndpointDescriptors []*EndpointDescriptor

// Direction returns the endpointDirection.
func (end *EndpointDescriptor) Direction() EndpointDirection {
	// FIXME(mdr): Is this funciton needed? What purpose does it serve? If I'm
	// keeping it, I should not return an unexported type.
	return end.EndpointAddress.direction()
}

// Number returns the endpoint number in bits 0..3 in the endpoint
// address.
func (end *EndpointDescriptor) Number() byte {
	return end.EndpointAddress.endpointNumber()
}

// TransferType returns the transfer type for an endpoint.
func (end *EndpointDescriptor) TransferType() TransferType {
	// FIXME(mdr): Is this funciton needed? What purpose does it serve? If I'm
	// keeping it, I should not return an unexported type.
	return end.Attributes.transferType()
}

func (address endpointAddress) direction() EndpointDirection {
	// Bit 7 of the endpointAddress determines the direction
	const directionMask = 0x80
	const directionBit = 7
	return EndpointDirection(address&directionMask) >> directionBit
}

func (address endpointAddress) endpointNumber() byte {
	// Bits 0..3 determine the endpoint number
	const endpointNumberMask = 0x0F
	return byte(address & endpointNumberMask)
}

func (attributes endpointAttributes) transferType() TransferType {
	// Bits 0..1 of the bmAttributes determines the transfer type
	const transferTypeMask = 0x03
	return TransferType(attributes & transferTypeMask)
}
