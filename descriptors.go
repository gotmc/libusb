// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"
import "fmt"

type classCode byte
type bcd uint16

// String implements the Stringer interface for bcd.
func (b bcd) String() string {
	return fmt.Sprintf("%#04x (%2.2f)",
		uint16(b),
		b.AsDecimal(),
	)
}

// AsDecimal converts the BCD value with a format 0xJJMN into a decimal JJ.MN
// where JJ is the major version number, M is the minor version, and N is the
// sub-minor version number.
func (b bcd) AsDecimal() float64 {
	return bcdToDecimal(uint16(b))
}

const (
	perInterface       classCode = C.LIBUSB_CLASS_PER_INTERFACE
	audio              classCode = C.LIBUSB_CLASS_AUDIO
	comm               classCode = C.LIBUSB_CLASS_COMM
	hid                classCode = C.LIBUSB_CLASS_HID
	physical           classCode = C.LIBUSB_CLASS_PHYSICAL
	printer            classCode = C.LIBUSB_CLASS_PRINTER
	ptp                classCode = C.LIBUSB_CLASS_PTP
	image              classCode = C.LIBUSB_CLASS_IMAGE
	massStorage        classCode = C.LIBUSB_CLASS_MASS_STORAGE
	hub                classCode = C.LIBUSB_CLASS_HUB
	data               classCode = C.LIBUSB_CLASS_DATA
	smartCard          classCode = C.LIBUSB_CLASS_SMART_CARD
	contentSecurity    classCode = C.LIBUSB_CLASS_CONTENT_SECURITY
	video              classCode = C.LIBUSB_CLASS_VIDEO
	personalHealthcare classCode = C.LIBUSB_CLASS_PERSONAL_HEALTHCARE
	diagnosticDevice   classCode = C.LIBUSB_CLASS_DIAGNOSTIC_DEVICE
	wireless           classCode = C.LIBUSB_CLASS_WIRELESS
	application        classCode = C.LIBUSB_CLASS_APPLICATION
	vendorSpec         classCode = C.LIBUSB_CLASS_VENDOR_SPEC
)

var classCodes = map[classCode]string{
	perInterface:       "Each interface specifies its own class information and all interfaces operate independently.",
	audio:              "Audio class.",
	comm:               "Communications class.",
	hid:                "Human Interface Device class.",
	physical:           "Physical.",
	printer:            "Printer class.",
	image:              "Image class.",
	massStorage:        "Mass storage class.",
	hub:                "Hub class.",
	data:               "Data class.",
	smartCard:          "Smart Card.",
	contentSecurity:    "Content Security.",
	video:              "Video.",
	personalHealthcare: "Personal Healthcare.",
	diagnosticDevice:   "Diagnostic Device.",
	wireless:           "Wireless class.",
	application:        "Application class.",
	vendorSpec:         "Class is vendor-specific.",
}

// String implements the Stringer interface for classCode.
func (classCode classCode) String() string {
	return classCodes[classCode]
}

type descriptorType byte

const (
	descDevice            descriptorType = C.LIBUSB_DT_DEVICE
	descConfig            descriptorType = C.LIBUSB_DT_CONFIG
	descString            descriptorType = C.LIBUSB_DT_STRING
	descInterface         descriptorType = C.LIBUSB_DT_INTERFACE
	descEndpoint          descriptorType = C.LIBUSB_DT_ENDPOINT
	descBos               descriptorType = C.LIBUSB_DT_BOS
	descDeviceCapability  descriptorType = C.LIBUSB_DT_DEVICE_CAPABILITY
	descHid               descriptorType = C.LIBUSB_DT_HID
	descReport            descriptorType = C.LIBUSB_DT_REPORT
	descPhysical          descriptorType = C.LIBUSB_DT_PHYSICAL
	descHub               descriptorType = C.LIBUSB_DT_HUB
	descSuperspeedHub     descriptorType = C.LIBUSB_DT_SUPERSPEED_HUB
	descEndpointCompanion descriptorType = C.LIBUSB_DT_SS_ENDPOINT_COMPANION
)

var descriptorTypes = map[descriptorType]string{
	descDevice:            "Device descriptor.",
	descConfig:            "Configuration descriptor.",
	descString:            "String descriptor.",
	descInterface:         "Interface descriptor.",
	descEndpoint:          "Endpoint descriptor.",
	descBos:               "BOS descriptor.",
	descDeviceCapability:  "Device Capability descriptor.",
	descHid:               "HID descriptor.",
	descReport:            "HID report descriptor.",
	descPhysical:          "Physical descriptor.",
	descHub:               "Hub descriptor.",
	descSuperspeedHub:     "SuperSpeed Hub descriptor.",
	descEndpointCompanion: "SuperSpeed Endpoint Companion descriptor.",
}

func (descriptorType descriptorType) String() string {
	return descriptorTypes[descriptorType]
}

// EndpointDirection provides the type for an in or out endpoint.
type EndpointDirection byte

const (
	// Per USB 2.0 spec bit 7 of the endpoint address defines the direction,
	// where 0 = OUT and 1 = IN. The libusb C.LIBUSB_ENDPOINT_IN enumeration is
	// 128 instead of 1. Therefore, I'm not using C.LIBUSB_ENDPOINT_IN (128).
	endpointOut   EndpointDirection = C.LIBUSB_ENDPOINT_OUT
	endpointIn    EndpointDirection = 1
	directionMask endpointAddress   = 0x80
	directionBit                    = 7
)

var endpointDirections = map[EndpointDirection]string{
	endpointOut: "Out: host-to-device.",
	endpointIn:  "In: device-to-host.",
}

// String implements the Stringer interface for endpointDirection.
func (endpointDirection EndpointDirection) String() string {
	return endpointDirections[endpointDirection]
}

// TransferType provides which type of transfer.
type TransferType int

// Endpoint transfer type http://bit.ly/enum_libusb_transfer_type
const (
	ControlTransfer     TransferType = C.LIBUSB_TRANSFER_TYPE_CONTROL
	IsochronousTransfer TransferType = C.LIBUSB_TRANSFER_TYPE_ISOCHRONOUS
	BulkTransfer        TransferType = C.LIBUSB_TRANSFER_TYPE_BULK
	InterruptTransfer   TransferType = C.LIBUSB_TRANSFER_TYPE_INTERRUPT
)

var transferTypes = map[TransferType]string{
	ControlTransfer:     "Control endpoint.",
	IsochronousTransfer: "Isochronous endpoint.",
	BulkTransfer:        "Bulk endpoint.",
	InterruptTransfer:   "Interrupt endpoint.",
}

func (transferType TransferType) String() string {
	return transferTypes[transferType]
}

// TODO(mdr): May want to replace uint8 with a type specific for indexes.

type synchronizationType byte

// Synchronization type for isochronous endpoints. "Values for bits 2:3 of the
// bmAttributes field in libusb_endpoint_descriptor"
// http://bit.ly/enum_libusb_iso_sync_type
const (
	IsoSyncTypeNone     synchronizationType = C.LIBUSB_ISO_SYNC_TYPE_NONE
	IsoSyncTypeAsync    synchronizationType = C.LIBUSB_ISO_SYNC_TYPE_ASYNC
	IsoSyncTypeAdaptive synchronizationType = C.LIBUSB_ISO_SYNC_TYPE_ADAPTIVE
	IsoSynceTypeSync    synchronizationType = C.LIBUSB_ISO_SYNC_TYPE_SYNC
)
