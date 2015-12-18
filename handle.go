// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"
import "unsafe"

// DeviceHandle represents the libusb device handle.
type DeviceHandle struct {
	libusbDeviceHandle *C.libusb_device_handle
}

func (dh *DeviceHandle) GetStringDescriptor(
	descIndex uint8,
	langID uint16,
) (string, error) {
	var cData *C.uchar
	length := 512
	usberr := C.libusb_get_string_descriptor(
		dh.libusbDeviceHandle,
		C.uint8_t(descIndex),
		C.uint16_t(langID),
		cData,
		C.int(length),
	)
	if usberr < 0 {
		return "", ErrorCode(usberr)
	}
	data := (*C.char)(unsafe.Pointer(cData))
	return C.GoString(data), nil
}

func (dh *DeviceHandle) GetStringDescriptorASCII(
	descIndex uint8,
) (string, error) {
	length := 256
	data := make([]byte, length)
	usberr := C.libusb_get_string_descriptor_ascii(
		dh.libusbDeviceHandle,
		C.uint8_t(descIndex),
		// Unsafe pointer -> http://stackoverflow.com/a/16376039/95592
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.int(length),
	)
	if usberr < 0 {
		return "", ErrorCode(usberr)
	}
	return string(data), nil
}

// Close implements libusb_close to close the device handle.
func (dh *DeviceHandle) Close() error {
	C.libusb_close(dh.libusbDeviceHandle)
	return nil
}

// GetDevice implements libusb_get_device to get the underlying device for a
// handle.
// TODO(mdr): Determine if I actually need this function.
// func (dh *DeviceHandle) GetDevice() (*Device, error) {
// }

// GetConfiguration implements the libusb_get_configuration function to
// determine the bConfigurationValue of the currently active configuration.
func (dh *DeviceHandle) GetConfiguration() (int, error) {
	var configuration *C.int
	err := C.libusb_get_configuration(dh.libusbDeviceHandle, configuration)
	if err != 0 {
		return 0, ErrorCode(err)
	}
	return int(*configuration), nil
}

// SetConfiguration implements libusb_set_configuration to set the active
// configuration for the device.
func (dh *DeviceHandle) SetConfiguration(configuration int) error {
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
	err := C.libusb_claim_interface(dh.libusbDeviceHandle, C.int(interfaceNum))
	if err != 0 {
		return ErrorCode(err)
	}
	return nil
}

// ReleaseInterface implements libusb_release_interface to release an interface
// previously claimed with libusb_claim_interface() (i.e., ClaimInterface()).
func (dh *DeviceHandle) ReleaseInterface(interfaceNum int) error {
	err := C.libusb_release_interface(dh.libusbDeviceHandle, C.int(interfaceNum))
	if err != 0 {
		return ErrorCode(err)
	}
	return nil
}

func (dh *DeviceHandle) SetInterfaceAltSetting(
	interfaceNum int,
	alternateSetting int,
) error {
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
	err := C.libusb_reset_device(dh.libusbDeviceHandle)
	if err != 0 {
		return ErrorCode(err)
	}
	return nil
}

// KernelDriverActive implements libusb_kernel_driver_active to determine if a
// kernel driver is active on an interface.
func (dh *DeviceHandle) KernelDriverActive(interfaceNum int) error {
	err := C.libusb_kernel_driver_active(
		dh.libusbDeviceHandle, C.int(interfaceNum))
	if err != 0 {
		return ErrorCode(err)
	}
	return nil
}

// DetachKernelDriver implements libusb_detach_kernel_driver to detach a kernel
// driver from an interface.
func (dh *DeviceHandle) DetachKernelDriver(interfaceNum int) error {
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
	cEnable := C.int(0)
	if enable == true {
		cEnable = C.int(1)
	}
	err := C.libusb_set_auto_detach_kernel_driver(dh.libusbDeviceHandle, cEnable)
	if err != 0 {
		return ErrorCode(err)
	}
	return nil
}

func (dh *DeviceHandle) BulkTransferOut(
	endpoint endpointAddress,
	data []byte,
	timeout int,
) (int, error) {
	var transferred C.int
	err := C.libusb_bulk_transfer(
		dh.libusbDeviceHandle,
		C.uchar(endpoint),
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.int(len(data)),
		&transferred,
		C.uint(timeout),
	)
	if err != 0 {
		return 0, ErrorCode(err)
	}
	return int(transferred), nil
}

func (dh *DeviceHandle) BulkTransferIn(
	endpoint endpointAddress,
	maxReceiveBytes int,
	timeout int,
) ([]byte, int, error) {
	data := make([]byte, maxReceiveBytes)
	var transferred C.int
	err := C.libusb_bulk_transfer(
		dh.libusbDeviceHandle,
		C.uchar(endpoint),
		(*C.uchar)(unsafe.Pointer(&data[0])),
		C.int(len(data)),
		&transferred,
		C.uint(timeout),
	)
	if err != 0 {
		return nil, 0, ErrorCode(err)
	}
	return data, int(transferred), nil
}
