// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"
import "fmt"

// ErrorCode is the tyep for the libusb_error C enum.
type ErrorCode int

// Error implements the Go error interface for ErrorCode.
func (err ErrorCode) Error() string {
	return fmt.Sprintf("%v: %v",
		ErrorName(err),
		StrError(err),
	)
}

// ErrorName implements the libusb_error_name function.
func ErrorName(err ErrorCode) string {
	return C.GoString(C.libusb_error_name(C.int(err)))
}

// StrError implements the libusb_strerror function.
func StrError(err ErrorCode) string {
	return C.GoString(C.libusb_strerror(int32(err)))
}

// SetLocale sets the locale for libusb errors.
func SetLocale(locale string) ErrorCode {
	return ErrorCode(C.libusb_setlocale(C.CString(locale)))
}

const (
	success           ErrorCode = C.LIBUSB_SUCCESS
	errorIo           ErrorCode = C.LIBUSB_ERROR_IO
	errorInvalidParam ErrorCode = C.LIBUSB_ERROR_INVALID_PARAM
	errorAccess       ErrorCode = C.LIBUSB_ERROR_ACCESS
	errorNoDevice     ErrorCode = C.LIBUSB_ERROR_NO_DEVICE
	errorNotFound     ErrorCode = C.LIBUSB_ERROR_NOT_FOUND
	errorBusy         ErrorCode = C.LIBUSB_ERROR_BUSY
	errorTimeout      ErrorCode = C.LIBUSB_ERROR_TIMEOUT
	errorOverflow     ErrorCode = C.LIBUSB_ERROR_OVERFLOW
	errorPipe         ErrorCode = C.LIBUSB_ERROR_PIPE
	errorInterrupted  ErrorCode = C.LIBUSB_ERROR_INTERRUPTED
	errorNoMem        ErrorCode = C.LIBUSB_ERROR_NO_MEM
	errorNotSupported ErrorCode = C.LIBUSB_ERROR_NOT_SUPPORTED
	errorOther        ErrorCode = C.LIBUSB_ERROR_OTHER
)
