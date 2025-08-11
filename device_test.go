// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import (
	"testing"
)

func TestDeviceClose(t *testing.T) {
	// Test closing a device multiple times
	dev := &Device{}

	// First close should work
	dev.Close()

	// Second close should be safe (no panic)
	dev.Close()

	// Verify device is nil
	if dev.libusbDevice != nil {
		t.Error("Device pointer should be nil after Close()")
	}
}

// Note: Testing Device methods with nil pointers would require adding nil checks
// to all Device methods, which would be an API change. The current implementation
// passes nil pointers directly to libusb C functions, which can cause segfaults.
// This is consistent with the C library behavior but not ideal for Go.
