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

func (err libusbError) Error() string {
	return fmt.Sprintf("libusb error code %d: %s", err, errorName(err))
}

func errorName(err libusbError) string {
	if errorString, ok := libusbErrorMap[err]; ok {
		return errorString
	}
	return "UNKNOWN"
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

var libusbErrorMap = map[libusbError]string{
	C.LIBUSB_SUCCESS:             "Success (no error)",
	C.LIBUSB_ERROR_IO:            "Input/output error.",
	C.LIBUSB_ERROR_INVALID_PARAM: "Invalid parameter.",
	C.LIBUSB_ERROR_ACCESS:        "Access denied (insufficient permissions)",
	C.LIBUSB_ERROR_NO_DEVICE:     "No such device (it may have been disconnected)",
	C.LIBUSB_ERROR_NOT_FOUND:     "Entity not found.",
	C.LIBUSB_ERROR_BUSY:          "Resource busy.",
	C.LIBUSB_ERROR_TIMEOUT:       "Operation timed out.",
	C.LIBUSB_ERROR_OVERFLOW:      "Overflow.",
	C.LIBUSB_ERROR_PIPE:          "Pipe error.",
	C.LIBUSB_ERROR_INTERRUPTED:   "System call interrupted (perhaps due to signal)",
	C.LIBUSB_ERROR_NO_MEM:        "Insufficient memory.",
	C.LIBUSB_ERROR_NOT_SUPPORTED: "Operation not supported or unimplemented on this platform.",
	C.LIBUSB_ERROR_OTHER:         "Other error.",
}
