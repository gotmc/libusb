// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"
import (
	"fmt"
	"reflect"
	"unsafe"
)

type logLevel int

// Log level enumeration
const (
	LogLevelNone    logLevel = C.LIBUSB_LOG_LEVEL_NONE
	LogLevelError   logLevel = C.LIBUSB_LOG_LEVEL_ERROR
	LogLevelWarning logLevel = C.LIBUSB_LOG_LEVEL_WARNING
	LogLevelInfo    logLevel = C.LIBUSB_LOG_LEVEL_INFO
	LogLevelDebug   logLevel = C.LIBUSB_LOG_LEVEL_DEBUG
)

var logLevels = map[logLevel]string{
	LogLevelNone:    "No messages ever printed by the library (default)",
	LogLevelError:   "Error messages are printed to stderr",
	LogLevelWarning: "Warning and error messages are printed to stderr",
	LogLevelInfo:    "Informational messages are printed to stdout, warning and error messages are printed to stderr",
	LogLevelDebug:   "Debug and informational messages are printed to stdout, warnings and errors to stderr",
}

func (level logLevel) String() string {
	return logLevels[level]
}

// Context represents a libusb session/context.
type Context struct {
	libusbContext *C.libusb_context
}

// Init intializes a new libusb session/context by creating a new Context and
// returning a pointer to that Context.
func Init() (*Context, error) {
	newContext := &Context{}
	errnum := C.libusb_init(&newContext.libusbContext)
	if errnum != 0 {
		return nil, fmt.Errorf(
			"Failed to initialize new libusb context. Received error %d", errnum)
	}
	return newContext, nil
}

// Exit deinitializes the libusb session/context.
func (ctx *Context) Exit() error {
	C.libusb_exit(ctx.libusbContext)
	ctx.libusbContext = nil
	return nil
}

// SetDebug sets the log message verbosity.
func (ctx *Context) SetDebug(level logLevel) {
	C.libusb_set_debug(ctx.libusbContext, C.int(level))
	return
}

// GetDeviceList returns an array of devices for the context.
func (ctx *Context) GetDeviceList() ([]*Device, error) {
	var devices []*Device
	var list **C.libusb_device
	const unrefDevices = 1
	numDevicesFound := int(C.libusb_get_device_list(ctx.libusbContext, &list))
	if numDevicesFound < 0 {
		return nil, ErrorCode(numDevicesFound)
	}
	defer C.libusb_free_device_list(list, unrefDevices)
	var libusbDevices []*C.libusb_device
	*(*reflect.SliceHeader)(unsafe.Pointer(&libusbDevices)) = reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(list)),
		Len:  numDevicesFound,
		Cap:  numDevicesFound,
	}
	for _, thisLibusbDevice := range libusbDevices {
		thisDevice := Device{
			libusbDevice: thisLibusbDevice,
		}
		devices = append(devices, &thisDevice)
	}
	return devices, nil
}

// OpenDeviceWithVendorProduct opens a USB device using the VendorID and
// productID and then returns a device handle.
func (ctx *Context) OpenDeviceWithVendorProduct(
	vendorID,
	productID uint16,
) (*Device, error) {
	var deviceHandle DeviceHandle
	deviceHandle.libusbDeviceHandle = C.libusb_open_device_with_vid_pid(
		ctx.libusbContext, C.uint16_t(vendorID), C.uint16_t(productID))
	if deviceHandle.libusbDeviceHandle == nil {
		return nil, fmt.Errorf("Could not open USB device %v:%v",
			vendorID,
			productID,
		)
	}
	device := Device{
		libusbDevice:     C.libusb_get_device(deviceHandle.libusbDeviceHandle),
		DeviceDescriptor: nil,
		DeviceHandle:     &deviceHandle,
	}
	err := device.GetDeviceDescriptor()
	if err != nil {
		return nil, err
	}
	return &device, nil
}
