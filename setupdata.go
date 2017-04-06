// Copyright (c) 2015-2016 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"

type bmRequestType byte

type transferDirection byte

// Constants to set the transfer direction.
const (
	HostToDevice transferDirection = 0x00
	DeviceToHost transferDirection = 0x80
)

var transferDirections = map[transferDirection]string{
	HostToDevice: "Host-to-device",
	DeviceToHost: "Device-to-host",
}

// String implements the Stringer interface for endpointDirection.
func (dir transferDirection) String() string {
	return transferDirections[dir]
}

type requestType byte

// Constants representing the libusb request types.
const (
	Standard requestType = C.LIBUSB_REQUEST_TYPE_STANDARD
	Class    requestType = C.LIBUSB_REQUEST_TYPE_CLASS
	Vendor   requestType = C.LIBUSB_REQUEST_TYPE_VENDOR
	Reserved requestType = C.LIBUSB_REQUEST_TYPE_RESERVED
)

var requestTypes = map[requestType]string{
	Standard: "Standard",
	Class:    "Class",
	Vendor:   "Vendor",
	Reserved: "Reserved",
}

func (rt requestType) String() string {
	return requestTypes[rt]
}

type requestRecipient byte

// Constants representing the libusb recipient types.
const (
	DeviceRecipient    requestRecipient = C.LIBUSB_RECIPIENT_DEVICE
	InterfaceRecipient requestRecipient = C.LIBUSB_RECIPIENT_INTERFACE
	EndpointRecipient  requestRecipient = C.LIBUSB_RECIPIENT_ENDPOINT
	OtherRecipient     requestRecipient = C.LIBUSB_RECIPIENT_OTHER
)

var requestRecipients = map[requestRecipient]string{
	DeviceRecipient:    "Device",
	InterfaceRecipient: "Interface",
	EndpointRecipient:  "Endpoint",
	OtherRecipient:     "Other",
}

func (r requestRecipient) String() string {
	return requestRecipients[r]
}

// BitmapRequestType returns the Request type.  Bits 0:4 determine recipient,
// see libusb_request_recipient. Bits 5:6 determine type, see
// libusb_request_type. Bit 7 determines data transfer direction, see
// libusb_endpoint_direction.
func BitmapRequestType(
	reqDirection transferDirection,
	reqType requestType,
	reqRecipient requestRecipient,
) bmRequestType {
	return bmRequestType(
		byte(reqDirection) | byte(reqType) | byte(reqRecipient))
}
