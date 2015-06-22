// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"

// TODO(mdr): Do I need to be handling the reference counts in cgo?

type device struct {
	libusbDevice *C.libusb_device
}

func (dev *device) GetBusNumber() (uint, error) {
	busNumber, err := C.libusb_get_bus_number(dev.libusbDevice)
	if err != nil {
		return 0, err
	}
	return uint(busNumber), nil
}

func (dev *device) GetDeviceAddress() (uint, error) {
	deviceAddress, err := C.libusb_get_device_address(dev.libusbDevice)
	if err != nil {
		return 0, err
	}
	return uint(deviceAddress), nil
}

func (dev *device) GetDeviceSpeed() (speed, error) {
	deviceSpeed, err := C.libusb_get_device_speed(dev.libusbDevice)
	if err != nil {
		return 0, err
	}
	return speed(deviceSpeed), nil
}

func (dev *device) Open() (*deviceHandle, error) {
	var handle **C.libusb_device_handle
	err := C.libusb_open(dev.libusbDevice, handle)
	if err != 0 {
		return nil, ErrorCode(err)
	}
	deviceHandle := deviceHandle{
		libusbDeviceHandle: *handle,
	}

	return &deviceHandle, nil
}

func (dev *device) GetDeviceDescriptor() (*deviceDescriptor, error) {
	var desc C.struct_libusb_device_descriptor
	err := C.libusb_get_device_descriptor(dev.libusbDevice, &desc)
	if err != 0 {
		return nil, ErrorCode(err)
	}
	descriptor := deviceDescriptor{
		Length:              uint8(desc.bLength),
		DescriptorType:      descriptorType(desc.bDescriptorType),
		USBSpecification:    bcd(desc.bcdUSB),
		DeviceClass:         classCode(desc.bDeviceClass),
		DeviceSubClass:      byte(desc.bDeviceSubClass),
		DeviceProtocol:      byte(desc.bDeviceProtocol),
		MaxPacketSize0:      uint8(desc.bMaxPacketSize0),
		VendorID:            uint16(desc.idVendor),
		ProductID:           uint16(desc.idProduct),
		DeviceReleaseNumber: bcd(desc.bcdDevice),
		ManufacturerIndex:   uint8(desc.iManufacturer),
		ProductIndex:        uint8(desc.iProduct),
		SerialNumberIndex:   uint8(desc.iSerialNumber),
		NumConfigurations:   uint8(desc.bNumConfigurations),
	}
	return &descriptor, nil
}
