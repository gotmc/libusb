// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

/*
Package libusb provides a Go bindings for the  libusb C library.
*/
package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"
import "fmt"

// Context represents a libusb session/context.
type context struct {
	context *C.libusb_context
}

// Init intializes a new libusb session/context by creating a new Context and
// returning a pointer to that Context.
func Init() (*context, error) {
	newContext := &context{}
	errnum := C.libusb_init(&newContext.context)
	if errnum != 0 {
		return nil, fmt.Errorf(
			"Failed to initialize new libusb context. Received error %d", errnum)
	}
	return newContext, nil
}

// Exit deinitializes the libusb session/context.
func (ctx *context) Exit() error {
	C.libusb_exit(ctx.context)
	ctx = nil
	return nil
}

// SetDebug sets the log message verbosity.
func (ctx *context) SetDebug(level int) error {
	C.libusb_set_debug(ctx.context, C.int(level))
	return nil
}
