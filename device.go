// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"

// TODO(mdr): Do I need to be handling the reference counts in cgo?

// Device represents a USB device including the opaque libusb_device struct.
type Device struct {
	libusbDevice *C.libusb_device
	*DeviceDescriptor
	*DeviceHandle
}

// GetBusNumber returns the bus number for the USB device.
func (dev *Device) GetBusNumber() (uint, error) {
	busNumber, err := C.libusb_get_bus_number(dev.libusbDevice)
	if err != nil {
		return 0, err
	}
	return uint(busNumber), nil
}

// GetDeviceAddress returns the address for the USB device.
func (dev *Device) GetDeviceAddress() (uint, error) {
	deviceAddress, err := C.libusb_get_device_address(dev.libusbDevice)
	if err != nil {
		return 0, err
	}
	return uint(deviceAddress), nil
}

// GetDeviceSpeed returns the speed for the USB device.
func (dev *Device) GetDeviceSpeed() (speed, error) {
	deviceSpeed, err := C.libusb_get_device_speed(dev.libusbDevice)
	if err != nil {
		return 0, err
	}
	return speed(deviceSpeed), nil
}

// Open opens a USB device and obtains a device handle, which is necessary for
// any I/O operations.
func (dev *Device) Open() error {
	var handle **C.libusb_device_handle
	err := C.libusb_open(dev.libusbDevice, handle)
	if err != 0 {
		return ErrorCode(err)
	}
	dev.DeviceHandle = &DeviceHandle{
		libusbDeviceHandle: *handle,
	}

	return nil
}

// GetDeviceDescriptor returns the USB device descriptor for the given USB
// device.
func (dev *Device) GetDeviceDescriptor() error {
	var desc C.struct_libusb_device_descriptor
	err := C.libusb_get_device_descriptor(dev.libusbDevice, &desc)
	if err != 0 {
		return ErrorCode(err)
	}
	dev.DeviceDescriptor = &DeviceDescriptor{
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
	return nil
}
