// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"
import (
	"unsafe"
)

// DeviceHandle represents the libusb device handle.
type DeviceHandle struct {
	libusbDeviceHandle *C.libusb_device_handle
}

// StringDescriptor retrieves a descriptor from a device.
func (dh *DeviceHandle) StringDescriptor(
	descIndex uint8,
	langID uint16,
) (string, error) {
	if dh == nil || dh.libusbDeviceHandle == nil {
		return "", ErrorCode(errorInvalidParam)
	}

	// Allocate buffer for data
	length := 512
	cData := make([]C.uchar, length)

	usberr := C.libusb_get_string_descriptor(
		dh.libusbDeviceHandle,
		C.uint8_t(descIndex),
		C.uint16_t(langID),
		&cData[0],
		C.int(length),
	)
	if usberr < 0 {
		return "", ErrorCode(usberr)
	}

	// Convert to Go string
	data := (*C.char)(unsafe.Pointer(&cData[0]))
	return C.GoString(data), nil
}

// StringDescriptorASCII retrieve(s) a string descriptor in C style ASCII.
// Wrapper around libusb_get_string_descriptor(). Uses the first language
// supported by the device. (Source: libusb docs)
func (dh *DeviceHandle) StringDescriptorASCII(
	descIndex uint8,
) (string, error) {
	if dh == nil || dh.libusbDeviceHandle == nil {
		return "", ErrorCode(errorInvalidParam)
	}

	// TODO(mdr): Should the length be a constant? Why did I pick 256 bytes?
	length := 256
	data := make([]byte, length)
	bytesRead, err := C.libusb_get_string_descriptor_ascii(
		dh.libusbDeviceHandle,
		C.uint8_t(descIndex),
		// Unsafe pointer -> https://stackoverflow.com/a/16376039/95592
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.int(length),
	)

	// Check both bytesRead and err
	if err != nil {
		return "", err
	}
	if bytesRead < 0 {
		return "", ErrorCode(bytesRead)
	}
	return string(data[0:bytesRead]), nil
}

// Close implements libusb_close to close the device handle.
func (dh *DeviceHandle) Close() error {
	if dh == nil || dh.libusbDeviceHandle == nil {
		return ErrorCode(errorInvalidParam)
	}
	C.libusb_close(dh.libusbDeviceHandle)
	return nil
}

// Device implements libusb_get_device to get the underlying device for a
// handle.
// TODO(mdr): Determine if I actually need this function.
// func (dh *DeviceHandle) Device() (*Device, error) {
// }

// Configuration implements the libusb_get_configuration function to
// determine the bConfigurationValue of the currently active configuration.
func (dh *DeviceHandle) Configuration() (int, error) {
	if dh == nil || dh.libusbDeviceHandle == nil {
		return 0, ErrorCode(errorInvalidParam)
	}

	// Allocate memory for the configuration value
	var configuration C.int
	err := C.libusb_get_configuration(dh.libusbDeviceHandle, &configuration)
	if err != 0 {
		return 0, ErrorCode(err)
	}
	return int(configuration), nil
}

// SetConfiguration implements libusb_set_configuration to set the active
// configuration for the device.
func (dh *DeviceHandle) SetConfiguration(configuration int) error {
	if dh == nil || dh.libusbDeviceHandle == nil {
		return ErrorCode(errorInvalidParam)
	}
	err := C.libusb_set_configuration(dh.libusbDeviceHandle,
		C.int(configuration))
	if err != 0 {
		return ErrorCode(err)
	}
	return nil
}

// ClaimInterface implements libusb_claim_interface to claim an interface on a
// given device handle. You must claim the interface you wish to use before you
// can perform I/O on any of its endpoints.
func (dh *DeviceHandle) ClaimInterface(interfaceNum int) error {
	if dh == nil || dh.libusbDeviceHandle == nil {
		return ErrorCode(errorInvalidParam)
	}
	err := C.libusb_claim_interface(dh.libusbDeviceHandle, C.int(interfaceNum))
	if err != 0 {
		return ErrorCode(err)
	}
	return nil
}

// ReleaseInterface implements libusb_release_interface to release an interface
// previously claimed with libusb_claim_interface() (i.e., ClaimInterface()).
func (dh *DeviceHandle) ReleaseInterface(interfaceNum int) error {
	if dh == nil || dh.libusbDeviceHandle == nil {
		return ErrorCode(errorInvalidParam)
	}
	err := C.libusb_release_interface(dh.libusbDeviceHandle, C.int(interfaceNum))
	if err != 0 {
		return ErrorCode(err)
	}
	return nil
}

// SetInterfaceAltSetting activates an alternate setting for an interface.
func (dh *DeviceHandle) SetInterfaceAltSetting(
	interfaceNum int,
	alternateSetting int,
) error {
	if dh == nil || dh.libusbDeviceHandle == nil {
		return ErrorCode(errorInvalidParam)
	}
	err := C.libusb_set_interface_alt_setting(
		dh.libusbDeviceHandle,
		C.int(interfaceNum),
		C.int(alternateSetting),
	)
	if err != 0 {
		return ErrorCode(err)
	}
	return nil
}

// FIXME(mdr): libusb_clear_halt takes an endpoint as an unsigned char. Need to
// determine, what I should pass into this function as the endpoint.
// func (dh *DeviceHandle) ClearHalt(endpoint int) error {
// return nil
// }

// ResetDevice implements libusb_reset_device to perform a USB port reset to
// reinitialize a device.
func (dh *DeviceHandle) ResetDevice() error {
	if dh == nil || dh.libusbDeviceHandle == nil {
		return ErrorCode(errorInvalidParam)
	}
	err := C.libusb_reset_device(dh.libusbDeviceHandle)
	if err != 0 {
		return ErrorCode(err)
	}
	return nil
}

// KernelDriverActive implements libusb_kernel_driver_active to determine if a
// kernel driver is active on an interface.
func (dh *DeviceHandle) KernelDriverActive(interfaceNum int) (bool, error) {
	if dh == nil || dh.libusbDeviceHandle == nil {
		return false, ErrorCode(errorInvalidParam)
	}
	ret := C.libusb_kernel_driver_active(
		dh.libusbDeviceHandle, C.int(interfaceNum))
	if ret == 1 {
		return true, nil
	} else if ret != 0 {
		return false, ErrorCode(ret)
	}
	return false, nil
}

// DetachKernelDriver implements libusb_detach_kernel_driver to detach a kernel
// driver from an interface.
func (dh *DeviceHandle) DetachKernelDriver(interfaceNum int) error {
	if dh == nil || dh.libusbDeviceHandle == nil {
		return ErrorCode(errorInvalidParam)
	}
	err := C.libusb_detach_kernel_driver(
		dh.libusbDeviceHandle, C.int(interfaceNum))
	if err != 0 {
		return ErrorCode(err)
	}
	return nil
}

// AttachKernelDriver implements libusb_attach_kernel_driver to re-attach an
// interface's kernel driver, which was previously detached using
// libusb_detach_kernel_driver().
func (dh *DeviceHandle) AttachKernelDriver(interfaceNum int) error {
	if dh == nil || dh.libusbDeviceHandle == nil {
		return ErrorCode(errorInvalidParam)
	}
	err := C.libusb_attach_kernel_driver(
		dh.libusbDeviceHandle, C.int(interfaceNum))
	if err != 0 {
		return ErrorCode(err)
	}
	return nil
}

// SetAutoDetachKernelDriver implements libusb_set_auto_detach_kernel_driver to
// enable/disable libusb's automatic kernel driver detachment.
func (dh *DeviceHandle) SetAutoDetachKernelDriver(enable bool) error {
	if dh == nil || dh.libusbDeviceHandle == nil {
		return ErrorCode(errorInvalidParam)
	}
	cEnable := C.int(0)
	if enable {
		cEnable = C.int(1)
	}
	err := C.libusb_set_auto_detach_kernel_driver(dh.libusbDeviceHandle, cEnable)
	if err != 0 {
		return ErrorCode(err)
	}
	return nil
}
