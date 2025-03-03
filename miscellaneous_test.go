// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import (
	"fmt"
	"testing"
)

func ExampleErrorCode_Error() {
	SetLocale("en")
	fmt.Println(success.Error())
	fmt.Println(errorIo.Error())
	fmt.Println(errorInvalidParam.Error())
	fmt.Println(errorAccess.Error())
	fmt.Println(errorNoDevice.Error())
	fmt.Println(errorNotFound.Error())
	fmt.Println(errorBusy.Error())
	fmt.Println(errorTimeout.Error())
	fmt.Println(errorOverflow.Error())
	fmt.Println(errorPipe.Error())
	fmt.Println(errorInterrupted.Error())
	fmt.Println(errorNoMem.Error())
	fmt.Println(errorNotSupported.Error())
	fmt.Println(errorOther.Error())

	// Output:
	// LIBUSB_SUCCESS / LIBUSB_TRANSFER_COMPLETED: Success
	// LIBUSB_ERROR_IO: Input/Output Error
	// LIBUSB_ERROR_INVALID_PARAM: Invalid parameter
	// LIBUSB_ERROR_ACCESS: Access denied (insufficient permissions)
	// LIBUSB_ERROR_NO_DEVICE: No such device (it may have been disconnected)
	// LIBUSB_ERROR_NOT_FOUND: Entity not found
	// LIBUSB_ERROR_BUSY: Resource busy
	// LIBUSB_ERROR_TIMEOUT: Operation timed out
	// LIBUSB_ERROR_OVERFLOW: Overflow
	// LIBUSB_ERROR_PIPE: Pipe error
	// LIBUSB_ERROR_INTERRUPTED: System call interrupted (perhaps due to signal)
	// LIBUSB_ERROR_NO_MEM: Insufficient memory
	// LIBUSB_ERROR_NOT_SUPPORTED: Operation not supported or unimplemented on this platform
	// LIBUSB_ERROR_OTHER: Other error
}

func TestBcdToDecimal(t *testing.T) {
	testCases := []struct {
		bcdValue uint16
		want     float64
	}{
		{0x0110, 1.1},
		{0x0200, 2.0},
		{0x0210, 2.1},
		{0x0300, 3.0},
		{0x0310, 3.1},
	}
	for _, tc := range testCases {
		if got := bcdToDecimal(tc.bcdValue); got != tc.want {
			t.Errorf(
				"Error converting BCD to decimal\n\tgot %v; want %v",
				got,
				tc.want,
			)
		}
	}
}
