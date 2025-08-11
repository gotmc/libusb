// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"

// TransferDirection represents the direction of a control transfer.
type TransferDirection byte

// Constants to set the transfer direction.
const (
	HostToDevice TransferDirection = 0x00
	DeviceToHost TransferDirection = 0x80
)

var transferDirections = map[TransferDirection]string{
	HostToDevice: "Host-to-device",
	DeviceToHost: "Device-to-host",
}

// String implements the Stringer interface for TransferDirection.
func (dir TransferDirection) String() string {
	return transferDirections[dir]
}

// RequestType represents the type of control request.
type RequestType byte

// Constants representing the libusb request types.
const (
	Standard RequestType = C.LIBUSB_REQUEST_TYPE_STANDARD
	Class    RequestType = C.LIBUSB_REQUEST_TYPE_CLASS
	Vendor   RequestType = C.LIBUSB_REQUEST_TYPE_VENDOR
	Reserved RequestType = C.LIBUSB_REQUEST_TYPE_RESERVED
)

var requestTypes = map[RequestType]string{
	Standard: "Standard",
	Class:    "Class",
	Vendor:   "Vendor",
	Reserved: "Reserved",
}

func (rt RequestType) String() string {
	return requestTypes[rt]
}

// RequestRecipient represents the recipient of a control request.
type RequestRecipient byte

// Constants representing the libusb recipient types.
const (
	DeviceRecipient    RequestRecipient = C.LIBUSB_RECIPIENT_DEVICE
	InterfaceRecipient RequestRecipient = C.LIBUSB_RECIPIENT_INTERFACE
	EndpointRecipient  RequestRecipient = C.LIBUSB_RECIPIENT_ENDPOINT
	OtherRecipient     RequestRecipient = C.LIBUSB_RECIPIENT_OTHER
)

var requestRecipients = map[RequestRecipient]string{
	DeviceRecipient:    "Device",
	InterfaceRecipient: "Interface",
	EndpointRecipient:  "Endpoint",
	OtherRecipient:     "Other",
}

func (r RequestRecipient) String() string {
	return requestRecipients[r]
}

// BitmapRequestType creates a bmRequestType byte from individual components.
// Bits 0:4 determine recipient, see libusb_request_recipient.
// Bits 5:6 determine type, see libusb_request_type.
// Bit 7 determines data transfer direction, see libusb_endpoint_direction.
func BitmapRequestType(dir TransferDirection, reqType RequestType,
	recipient RequestRecipient) byte {
	return byte(dir) | byte(reqType) | byte(recipient)
}
