// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import "fmt"

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
