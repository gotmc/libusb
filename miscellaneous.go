// Copyright (c) 2015-2020 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"
import (
	"fmt"
	"math"
)

// ErrorCode is the type for the libusb_error C enum.
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

// CPUtoLE16 converts "a 16-bit value from host-endian to little-endian format.
// On little endian systems, this function does nothing. On big endian systems,
// the bytes are swapped.
func CPUtoLE16(value int) int {
	return int(C.libusb_cpu_to_le16(C.uint16_t(value)))
}

// HasCapability checks "at runtime if the loaded library has a given
// capability. This call should be performed after libusb_init(), to ensure
// the backend has updated its capability set." (Source: libusb docs)
func HasCapability(capability int) bool {
	if C.libusb_has_capability(C.uint32_t(capability)) != 0 {
		return true
	}
	return false
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

	errorTransferError    ErrorCode = C.LIBUSB_TRANSFER_ERROR
	errorTransferTimedOut ErrorCode = C.LIBUSB_TRANSFER_TIMED_OUT
	errorTransferCanceled ErrorCode = C.LIBUSB_TRANSFER_CANCELLED
	errorTransferStall    ErrorCode = C.LIBUSB_TRANSFER_STALL
	errorTransferNoDevice ErrorCode = C.LIBUSB_TRANSFER_NO_DEVICE
	errorTransferOverflow ErrorCode = C.LIBUSB_TRANSFER_OVERFLOW
)

func bcdToDecimal(bcdValue uint16) float64 {
	bcdPowersByPosition := []string{"hundreths", "tenths", "ones", "tens"}

	var bcdMap map[string]uint16
	bcdMap = make(map[string]uint16)
	for i, power := range bcdPowersByPosition {
		bcdMap[power] = bcdValue & (0xf << uint(4*i)) / uint16(math.Pow(16, float64(i)))
	}
	return 10*float64(bcdMap["tens"]) + float64(bcdMap["ones"]) +
		0.1*float64(bcdMap["tenths"]) + 0.01*float64(bcdMap["hundreths"])
}
