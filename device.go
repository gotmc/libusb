// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

// Device represents a USB device including the opaque libusb_device struct.
type Device struct {
	libusbDevice        *C.libusb_device
	ActiveConfiguration *ConfigDescriptor
}

// deviceFinalizer is called by the garbage collector to clean up
// unreferenced Device objects that weren't explicitly closed.
func deviceFinalizer(dev *Device) {
	if dev.libusbDevice != nil {
		C.libusb_unref_device(dev.libusbDevice)
		dev.libusbDevice = nil
	}
}

// newDevice creates a new Device with proper finalizer setup.
func newDevice(libusbDevice *C.libusb_device) *Device {
	dev := &Device{
		libusbDevice: libusbDevice,
	}
	runtime.SetFinalizer(dev, deviceFinalizer)
	return dev
}

// Close decrements the reference count of the device. If the decrement
// operation causes the reference count to reach zero, the device shall be
// destroyed.
func (dev *Device) Close() {
	if dev.libusbDevice != nil {
		C.libusb_unref_device(dev.libusbDevice)
		dev.libusbDevice = nil
		// Clear finalizer since we've explicitly closed the device
		runtime.SetFinalizer(dev, nil)
	}
}

// Descriptor represents a USB device descriptor as a Go struct.
type Descriptor struct {
	Length              uint8
	DescriptorType      descriptorType
	USBSpecification    bcd
	DeviceClass         classCode
	DeviceSubClass      byte
	DeviceProtocol      byte
	MaxPacketSize0      uint8
	VendorID            uint16
	ProductID           uint16
	DeviceReleaseNumber bcd
	ManufacturerIndex   uint8
	ProductIndex        uint8
	SerialNumberIndex   uint8
	NumConfigurations   uint8
}

// BusNumber gets "the number of the bus that a device is connected to."
// (Source: libusb docs)
func (dev *Device) BusNumber() (int, error) {
	if dev == nil || dev.libusbDevice == nil {
		return 0, ErrorCode(errorInvalidParam)
	}
	busNumber, err := C.libusb_get_bus_number(dev.libusbDevice)
	if err != nil {
		return 0, err
	}
	return int(busNumber), nil
}

// PortNumber gets "the number of the port that a device is connected to.
// Unless the OS does something funky, or you are hot-plugging USB extension
// cards, the port number returned by this call is usually guaranteed to be
// uniquely tied to a physical port, meaning that different devices plugged on
// the same physical port should return the same port number.  But outside of
// this, there is no guarantee that the port number returned by this call will
// remain the same, or even match the order in which ports have been numbered
// by the HUB/HCD manufacturer." (Source: libusb docs)
func (dev *Device) PortNumber() (int, error) {
	if dev == nil || dev.libusbDevice == nil {
		return 0, ErrorCode(errorInvalidParam)
	}
	portNumber, err := C.libusb_get_port_number(dev.libusbDevice)
	if err != nil {
		return 0, fmt.Errorf("port number is unavailable for device %v", dev)
	}
	return int(portNumber), nil
}

// MaxPacketSize is a "convenience function to retrieve the wMaxPacketSize
// value for a particular endpoint in the active device configuration. This
// function was originally intended to be of assistance when setting up
// isochronous transfers, but a design mistake resulted in this function
// instead. It simply returns the wMaxPacketSize value without considering its
// contents. If you're dealing with isochronous transfers, you probably want
// libusb_get_max_iso_packet_size() instead." (Source: libusb docs)
func (dev *Device) MaxPacketSize(ep endpointAddress) (int, error) {
	if dev == nil || dev.libusbDevice == nil {
		return 0, ErrorCode(errorInvalidParam)
	}
	maxPacketSize, err := C.libusb_get_max_packet_size(dev.libusbDevice, C.uchar(ep))
	if err != nil {
		return 0, fmt.Errorf("wMaxPacketSize is unavailable for device %v", dev)
	}
	return int(maxPacketSize), nil
}

// DeviceAddress gets "the address of the device on the bus it is connected
// to." (Source: libusb docs)
func (dev *Device) DeviceAddress() (int, error) {
	if dev == nil || dev.libusbDevice == nil {
		return 0, ErrorCode(errorInvalidParam)
	}
	deviceAddress, err := C.libusb_get_device_address(dev.libusbDevice)
	if err != nil {
		return 0, err
	}
	return int(deviceAddress), nil
}

// Speed gets "the negotiated connection speed for a device." (Source:
// libusb docs)
func (dev *Device) Speed() (SpeedType, error) {
	if dev == nil || dev.libusbDevice == nil {
		return 0, ErrorCode(errorInvalidParam)
	}
	deviceSpeed, err := C.libusb_get_device_speed(dev.libusbDevice)
	if err != nil {
		return 0, err
	}
	return SpeedType(deviceSpeed), nil
}

// Open will "open a device and obtain a device handle. A handle allows you to
// perform I/O on the device in question. Internally, this function adds a
// reference to the device and makes it available to you through
// libusb_get_device(). This reference is removed during libusb_close()." This
// is a non-blocking function; no requests are sent over the bus. (Source:
// libusb docs)
func (dev *Device) Open() (*DeviceHandle, error) {
	if dev == nil || dev.libusbDevice == nil {
		return nil, ErrorCode(errorInvalidParam)
	}
	var handle *C.libusb_device_handle
	err := C.libusb_open(dev.libusbDevice, &handle)
	if err != 0 {
		return nil, ErrorCode(err)
	}
	deviceHandle := newDeviceHandle(handle)
	return deviceHandle, nil
}

// DeviceDescriptor implements the libusb_get_device_descriptor function to
// update the DeviceDescriptor struct embedded in the Device.  DeviceDescriptor
// gets "the USB device descriptor for a given device. This is a non-blocking
// function; the device descriptor is cached in memory. Note since
// libusb-1.0.16, LIBUSB_API_VERSION >= 0x01000102, this function always
// succeeds." (Source: libusb docs)
func (dev *Device) DeviceDescriptor() (*Descriptor, error) {
	if dev == nil || dev.libusbDevice == nil {
		return nil, ErrorCode(errorInvalidParam)
	}
	var desc C.struct_libusb_device_descriptor
	err := C.libusb_get_device_descriptor(dev.libusbDevice, &desc)
	if err != 0 {
		return nil, ErrorCode(err)
	}
	deviceDescriptor := Descriptor{
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
	return &deviceDescriptor, nil
}

// parseConfigDescriptor converts a C libusb_config_descriptor into a Go
// ConfigDescriptor, fully populating SupportedInterfaces with all interface
// and endpoint descriptors.
func parseConfigDescriptor(
	config *C.struct_libusb_config_descriptor,
) *ConfigDescriptor {
	cd := &ConfigDescriptor{
		Length:               int(config.bLength),
		DescriptorType:       descriptorType(config.bDescriptorType),
		TotalLength:          uint16(config.wTotalLength),
		NumInterfaces:        int(config.bNumInterfaces),
		ConfigurationValue:   uint8(config.bConfigurationValue),
		ConfigurationIndex:   uint8(config.iConfiguration),
		Attributes:           uint8(config.bmAttributes),
		MaxPowerMilliAmperes: 2 * uint(config.MaxPower),
	}
	numInterfaces := cd.NumInterfaces
	if numInterfaces == 0 {
		return cd
	}
	libusbInterfaces := unsafe.Slice(config._interface, numInterfaces)
	supportedInterfaces := make(SupportedInterfaces, 0, numInterfaces)
	for _, libusbInterface := range libusbInterfaces {
		numAlt := int(libusbInterface.num_altsetting)
		supportedInterface := SupportedInterface{
			NumAltSettings: numAlt,
		}
		if numAlt == 0 {
			supportedInterfaces = append(supportedInterfaces, &supportedInterface)
			continue
		}
		libusbInterfaceDescriptors := unsafe.Slice(
			libusbInterface.altsetting, numAlt,
		)
		interfaceDescriptors := make(InterfaceDescriptors, 0, numAlt)
		for _, lid := range libusbInterfaceDescriptors {
			ifaceDesc := InterfaceDescriptor{
				Length:            int(lid.bLength),
				DescriptorType:    descriptorType(lid.bDescriptorType),
				InterfaceNumber:   int(lid.bInterfaceNumber),
				AlternateSetting:  int(lid.bAlternateSetting),
				NumEndpoints:      int(lid.bNumEndpoints),
				InterfaceClass:    uint8(lid.bInterfaceClass),
				InterfaceSubClass: uint8(lid.bInterfaceSubClass),
				InterfaceProtocol: uint8(lid.bInterfaceProtocol),
				InterfaceIndex:    int(lid.iInterface),
			}
			numEP := int(lid.bNumEndpoints)
			if numEP > 0 {
				libusbEndpoints := unsafe.Slice(lid.endpoint, numEP)
				epDescs := make(EndpointDescriptors, 0, numEP)
				for _, lep := range libusbEndpoints {
					epDescs = append(epDescs, &EndpointDescriptor{
						Length:          int(lep.bLength),
						DescriptorType:  descriptorType(lep.bDescriptorType),
						EndpointAddress: endpointAddress(lep.bEndpointAddress),
						Attributes:      endpointAttributes(lep.bmAttributes),
						MaxPacketSize:   uint16(lep.wMaxPacketSize),
						Interval:        uint8(lep.bInterval),
					})
				}
				ifaceDesc.EndpointDescriptors = epDescs
			}
			interfaceDescriptors = append(interfaceDescriptors, &ifaceDesc)
		}
		supportedInterface.InterfaceDescriptors = interfaceDescriptors
		supportedInterfaces = append(supportedInterfaces, &supportedInterface)
	}
	cd.SupportedInterfaces = supportedInterfaces
	return cd
}

// ActiveConfigDescriptor "gets the USB configuration descriptor for the
// currently active configuration. This is a non-blocking function which does
// not involve any requests being sent to the device." (Source: libusb docs)
func (dev *Device) ActiveConfigDescriptor() (*ConfigDescriptor, error) {
	if dev == nil || dev.libusbDevice == nil {
		return nil, ErrorCode(errorInvalidParam)
	}
	var config *C.struct_libusb_config_descriptor
	err := C.libusb_get_active_config_descriptor(dev.libusbDevice, &config)
	if err != 0 {
		return nil, ErrorCode(err)
	}
	defer C.libusb_free_config_descriptor(config)
	return parseConfigDescriptor(config), nil
}

// ConfigDescriptor "gets a USB configuration descriptor based on its index.
// This is a non-blocking function which does not involve any requests being
// sent to the device." (Source: libusb docs)
func (dev *Device) ConfigDescriptor(configIndex int) (*ConfigDescriptor, error) {
	if dev == nil || dev.libusbDevice == nil {
		return nil, ErrorCode(errorInvalidParam)
	}
	var cConfig *C.struct_libusb_config_descriptor
	err := C.libusb_get_config_descriptor(
		dev.libusbDevice, C.uint8_t(configIndex), &cConfig,
	)
	if err != 0 {
		return nil, ErrorCode(err)
	}
	defer C.libusb_free_config_descriptor(cConfig)
	return parseConfigDescriptor(cConfig), nil
}

// ConfigDescriptorByValue gets "a USB configuration descriptor with a
// specific bConfigurationValue. This is a non-blocking function which does not
// involve any requests being sent to the device. (Source: libusb docs)
func (dev *Device) ConfigDescriptorByValue(configValue int) (*ConfigDescriptor, error) {
	if dev == nil || dev.libusbDevice == nil {
		return nil, ErrorCode(errorInvalidParam)
	}
	var cConfig *C.struct_libusb_config_descriptor
	err := C.libusb_get_config_descriptor_by_value(
		dev.libusbDevice, C.uint8_t(configValue), &cConfig,
	)
	if err != 0 {
		return nil, ErrorCode(err)
	}
	defer C.libusb_free_config_descriptor(cConfig)
	return parseConfigDescriptor(cConfig), nil
}

// FindInterfacesByClass finds all interfaces that match the given USB class code.
// This is particularly useful for devices where the device class is reported as
// LIBUSB_CLASS_PER_INTERFACE (0), where individual interfaces have their own class codes.
//
// For example, to find all printer class interfaces:
//
//	printerInterfaces, err := device.FindInterfacesByClass(libusb.InterfaceClassPrinter)
//
// Example of finding and printing information about a printer device:
//
//	package main
//
//	import (
//		"fmt"
//		"github.com/gotmc/libusb/v2"
//	)
//
//	func main() {
//		ctx, _ := libusb.NewContext()
//		devices, _ := ctx.DeviceList()
//
//		for _, device := range devices {
//			desc, _ := device.DeviceDescriptor()
//			// Check if this is a per-interface class device
//			if desc.DeviceClass == 0 {
//				// Find printer interfaces (class 7)
//				printerIfaces, _ := device.FindInterfacesByClass(libusb.InterfaceClassPrinter)
//				if len(printerIfaces) > 0 {
//					fmt.Printf("Found printer device: VID=0x%04x, PID=0x%04x\n", desc.VendorID, desc.ProductID)
//				}
//			}
//		}
//	}
func (dev *Device) FindInterfacesByClass(
	class uint8,
) (InterfaceDescriptors, error) {
	if dev == nil || dev.libusbDevice == nil {
		return nil, ErrorCode(errorInvalidParam)
	}
	config, err := dev.ActiveConfigDescriptor()
	if err != nil {
		return nil, err
	}
	return config.GetAllInterfacesByClass(class), nil
}
