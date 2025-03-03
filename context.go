// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
// /* Define LIBUSB_OPTION_LOG_LEVEL if not available in the libusb version */
// #if !defined(LIBUSB_OPTION_LOG_LEVEL) && defined(LIBUSB_API_VERSION) && (LIBUSB_API_VERSION >= 0x01000107)
// #define LIBUSB_OPTION_LOG_LEVEL 1
// #endif
//
// int set_debug(libusb_context * ctx, int level) {
// #if defined(HAVE_LIBUSB_SET_OPTION) && defined(LIBUSB_OPTION_LOG_LEVEL)
//    return libusb_set_option(ctx, LIBUSB_OPTION_LOG_LEVEL, level);
// #else
//    libusb_set_debug(ctx, level);
//    return 0;
// #endif
// }
import "C"

import (
	"fmt"
	"unsafe"
)

// LogLevel is an enum for the C libusb log message levels.
type LogLevel int

// Log message levels
//
// http://bit.ly/enum_libusb_log_level
const (
	LogLevelNone    LogLevel = C.LIBUSB_LOG_LEVEL_NONE
	LogLevelError   LogLevel = C.LIBUSB_LOG_LEVEL_ERROR
	LogLevelWarning LogLevel = C.LIBUSB_LOG_LEVEL_WARNING
	LogLevelInfo    LogLevel = C.LIBUSB_LOG_LEVEL_INFO
	LogLevelDebug   LogLevel = C.LIBUSB_LOG_LEVEL_DEBUG
)

var logLevels = map[LogLevel]string{
	LogLevelNone:    "No messages ever printed by the library (default)",
	LogLevelError:   "Error messages are printed to stderr",
	LogLevelWarning: "Warning and error messages are printed to stderr",
	LogLevelInfo:    "Informational messages are printed to stdout, warning and error messages are printed to stderr",
	LogLevelDebug:   "Debug and informational messages are printed to stdout, warnings and errors to stderr",
}

func (level LogLevel) String() string {
	return logLevels[level]
}

// Context represents a libusb session/context.
type Context struct {
	libusbContext *C.libusb_context
	LogLevel      LogLevel
}

// NewContext intializes a new libusb session/context by creating a new
// Context and returning a pointer to that Context.
func NewContext() (*Context, error) {
	newContext := &Context{
		LogLevel: LogLevelNone,
	}
	errnum := C.libusb_init(&newContext.libusbContext)
	if errnum != 0 {
		return nil, fmt.Errorf(
			"failed to initialize new libusb context; received error %d", errnum)
	}
	return newContext, nil
}

// Close deinitializes the libusb session/context.
func (ctx *Context) Close() error {
	C.libusb_exit(ctx.libusbContext)
	ctx.libusbContext = nil
	return nil
}

// SetDebug sets the log message verbosity.
func (ctx *Context) SetDebug(level LogLevel) {
	C.set_debug(ctx.libusbContext, C.int(level))
	ctx.LogLevel = level
}

// DeviceList returns an array of devices for the context.
func (ctx *Context) DeviceList() ([]*Device, error) {
	var devices []*Device
	var list **C.libusb_device
	const unrefDevices = 1
	numDevicesFound := int(C.libusb_get_device_list(ctx.libusbContext, &list))
	if numDevicesFound < 0 {
		return nil, ErrorCode(numDevicesFound)
	}
	defer C.libusb_free_device_list(list, unrefDevices)
	libusbDevices := unsafe.Slice(list, numDevicesFound)
	// var libusbDevices []*C.libusb_device
	// *(*reflect.SliceHeader)(unsafe.Pointer(&libusbDevices)) = reflect.SliceHeader{
	// 	Data: uintptr(unsafe.Pointer(list)),
	// 	Len:  numDevicesFound,
	// 	Cap:  numDevicesFound,
	// }

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
	vendorID uint16,
	productID uint16,
) (*Device, *DeviceHandle, error) {
	var deviceHandle DeviceHandle
	deviceHandle.libusbDeviceHandle = C.libusb_open_device_with_vid_pid(
		ctx.libusbContext, C.uint16_t(vendorID), C.uint16_t(productID))
	if deviceHandle.libusbDeviceHandle == nil {
		return nil, nil, fmt.Errorf("could not open USB device %v:%v",
			vendorID,
			productID,
		)
	}
	device := Device{
		libusbDevice: C.libusb_get_device(deviceHandle.libusbDeviceHandle),
	}
	return &device, &deviceHandle, nil
}
