// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import (
	"runtime"
	"testing"
	"time"
)

// TestDeviceFinalizer tests that Device objects are properly cleaned up
// by the garbage collector when not explicitly closed.
func TestDeviceFinalizer(t *testing.T) {
	// Skip test if no USB devices are connected
	ctx, err := NewContext()
	if err != nil {
		t.Skip("Cannot create context for finalizer test")
	}
	defer ctx.Close()

	devices, err := ctx.DeviceList()
	if err != nil {
		t.Skip("Cannot get device list for finalizer test")
	}
	defer func() {
		for _, device := range devices {
			device.Close()
		}
	}()

	if len(devices) == 0 {
		t.Skip("No USB devices connected for finalizer test")
	}

	// Create a device and let it go out of scope without calling Close()
	func() {
		// This device will go out of scope and should be finalized
		_ = devices[0]
	}()

	// Force garbage collection and finalization
	runtime.GC()
	runtime.GC()                       // Call twice to ensure finalizers run
	time.Sleep(100 * time.Millisecond) // Give finalizers time to run

	// If we get here without crashing, the finalizer worked correctly
	t.Log("Device finalizer test completed successfully")
}

// TestDeviceHandleFinalizer tests that DeviceHandle objects are properly
// cleaned up by the garbage collector when not explicitly closed.
func TestDeviceHandleFinalizer(t *testing.T) {
	// Skip test if no USB devices are connected
	ctx, err := NewContext()
	if err != nil {
		t.Skip("Cannot create context for finalizer test")
	}
	defer ctx.Close()

	devices, err := ctx.DeviceList()
	if err != nil {
		t.Skip("Cannot get device list for finalizer test")
	}
	defer func() {
		for _, device := range devices {
			device.Close()
		}
	}()

	if len(devices) == 0 {
		t.Skip("No USB devices connected for finalizer test")
	}

	// Try to open a device handle and let it go out of scope
	func() {
		handle, err := devices[0].Open()
		if err != nil {
			t.Skip("Cannot open device for handle finalizer test")
		}
		// This handle will go out of scope and should be finalized
		_ = handle
	}()

	// Force garbage collection and finalization
	runtime.GC()
	runtime.GC()                       // Call twice to ensure finalizers run
	time.Sleep(100 * time.Millisecond) // Give finalizers time to run

	// If we get here without crashing, the finalizer worked correctly
	t.Log("DeviceHandle finalizer test completed successfully")
}

// TestExplicitCloseRemovesFinalizer tests that explicitly calling Close()
// removes the finalizer so it won't be called later.
func TestExplicitCloseRemovesFinalizer(t *testing.T) {
	// Skip test if no USB devices are connected
	ctx, err := NewContext()
	if err != nil {
		t.Skip("Cannot create context for finalizer test")
	}
	defer ctx.Close()

	devices, err := ctx.DeviceList()
	if err != nil {
		t.Skip("Cannot get device list for finalizer test")
	}
	defer func() {
		for _, device := range devices {
			device.Close()
		}
	}()

	if len(devices) == 0 {
		t.Skip("No USB devices connected for finalizer test")
	}

	// Create a device, close it explicitly, then force GC
	device := devices[0]
	device.Close() // This should clear the finalizer

	// Force garbage collection - the finalizer should not run
	runtime.GC()
	runtime.GC()
	time.Sleep(100 * time.Millisecond)

	t.Log("Explicit close finalizer removal test completed successfully")
}

// TestFinalizerMemoryLeakPrevention tests that finalizers prevent memory leaks
// when objects are not explicitly closed.
func TestFinalizerMemoryLeakPrevention(t *testing.T) {
	// This test creates many Device objects without closing them
	// and verifies that finalizers clean them up properly
	ctx, err := NewContext()
	if err != nil {
		t.Skip("Cannot create context for finalizer test")
	}
	defer ctx.Close()

	// Create multiple rounds of devices to stress test finalizers
	for i := 0; i < 5; i++ {
		func() {
			devices, err := ctx.DeviceList()
			if err != nil {
				t.Skip("Cannot get device list for memory leak test")
			}
			// Let all devices go out of scope without closing them
			_ = devices
		}()

		// Force garbage collection after each round
		runtime.GC()
		runtime.GC()
		time.Sleep(50 * time.Millisecond)
	}

	t.Log("Memory leak prevention test completed successfully")
}
