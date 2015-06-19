// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"
import "fmt"

type libusbError int

// Error implements the Go error interface for libusbError.
func (err libusbError) Error() string {
	return fmt.Sprintf("%v: %v",
		errorName(err),
		strError(err),
	)
}

// errorName implements the libusb_error_name function.
func errorName(err libusbError) string {
	return C.GoString(C.libusb_error_name(C.int(err)))
}

// strerror implements the libusb_strerror function.
func strError(err libusbError) string {
	return C.GoString(C.libusb_strerror(int32(err)))
}

func SetLocale(locale string) libusbError {
	return libusbError(C.libusb_setlocale(C.CString(locale)))
}

const (
	success           libusbError = C.LIBUSB_SUCCESS
	errorIo           libusbError = C.LIBUSB_ERROR_IO
	errorInvalidParam libusbError = C.LIBUSB_ERROR_INVALID_PARAM
	errorAccess       libusbError = C.LIBUSB_ERROR_ACCESS
	errorNoDevice     libusbError = C.LIBUSB_ERROR_NO_DEVICE
	errorNotFound     libusbError = C.LIBUSB_ERROR_NOT_FOUND
	errorBusy         libusbError = C.LIBUSB_ERROR_BUSY
	errorTimeout      libusbError = C.LIBUSB_ERROR_TIMEOUT
	errorOverflow     libusbError = C.LIBUSB_ERROR_OVERFLOW
	errorPipe         libusbError = C.LIBUSB_ERROR_PIPE
	errorInterrupted  libusbError = C.LIBUSB_ERROR_INTERRUPTED
	errorNoMem        libusbError = C.LIBUSB_ERROR_NO_MEM
	errorNotSupported libusbError = C.LIBUSB_ERROR_NOT_SUPPORTED
	errorOther        libusbError = C.LIBUSB_ERROR_OTHER
)
