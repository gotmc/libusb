// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"
import "fmt"

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

// GetPortNumber returns the port number for the given USB device.
func (dev *Device) GetPortNumber() (uint, error) {
	portNumber, err := C.libusb_get_port_number(dev.libusbDevice)
	if err != nil {
		return 0, fmt.Errorf("Port number is unavailable for device %v", dev)
	}
	return uint(portNumber), nil
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
		ActiveConfiguration: nil,
	}
	return nil
}

// ResetDevice performs a USB port reset to reinitialize a device.
//
// Per libusb: "The system will attempt to restore the previous configuration
// and alternate settings after the reset has completed. If the reset fails,
// the descriptors change, or the previous state cannot be restored, the device
// will appear to be disconnected and reconnected. This means that the device
// handle is no longer valid (you should close it) and rediscover the device. A
// return code of LIBUSB_ERROR_NOT_FOUND indicates when this is the case.  This
// is a blocking function which usually incurs a noticeable delay.
func (dev *Device) ResetDevice() error {
	err := C.libusb_reset_device(dev.libusbDeviceHandle)
	if err != 0 {
		return ErrorCode(err)
	}
	return nil
}

func (dev *Device) GetActiveConfigDescriptor() error {
	var config *C.struct_libusb_config_descriptor
	err := C.libusb_get_active_config_descriptor(dev.libusbDevice, &config)
	defer C.libusb_free_config_descriptor(config)
	if err != 0 {
		return ErrorCode(err)
	}
	dev.ActiveConfiguration = &ConfigurationDescriptor{
		Length:               uint8(config.bLength),
		DescriptorType:       descriptorType(config.bDescriptorType),
		TotalLength:          uint16(config.wTotalLength),
		NumInterfaces:        uint8(config.bNumInterfaces),
		ConfigurationValue:   uint8(config.bConfigurationValue),
		ConfigurationIndex:   uint8(config.iConfiguration),
		Attributes:           uint8(config.bmAttributes),
		MaxPowerMilliAmperes: 2 * uint(config.MaxPower), // Convert from 2 mA to just mA
		Interfaces:           nil,
	}
	return nil
}
