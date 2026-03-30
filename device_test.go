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

func TestDeviceNilChecks(t *testing.T) {
	var dev *Device

	if _, err := dev.BusNumber(); err != ErrorCode(errorInvalidParam) {
		t.Errorf("BusNumber: got %v, want errorInvalidParam", err)
	}
	if _, err := dev.PortNumber(); err != ErrorCode(errorInvalidParam) {
		t.Errorf("PortNumber: got %v, want errorInvalidParam", err)
	}
	if _, err := dev.MaxPacketSize(0); err != ErrorCode(errorInvalidParam) {
		t.Errorf("MaxPacketSize: got %v, want errorInvalidParam", err)
	}
	if _, err := dev.DeviceAddress(); err != ErrorCode(errorInvalidParam) {
		t.Errorf("DeviceAddress: got %v, want errorInvalidParam", err)
	}
	if _, err := dev.Speed(); err != ErrorCode(errorInvalidParam) {
		t.Errorf("Speed: got %v, want errorInvalidParam", err)
	}
	if _, err := dev.Open(); err != ErrorCode(errorInvalidParam) {
		t.Errorf("Open: got %v, want errorInvalidParam", err)
	}
	if _, err := dev.DeviceDescriptor(); err != ErrorCode(errorInvalidParam) {
		t.Errorf("DeviceDescriptor: got %v, want errorInvalidParam", err)
	}
	if _, err := dev.ActiveConfigDescriptor(); err != ErrorCode(errorInvalidParam) {
		t.Errorf(
			"ActiveConfigDescriptor: got %v, want errorInvalidParam", err,
		)
	}
	if _, err := dev.ConfigDescriptor(0); err != ErrorCode(errorInvalidParam) {
		t.Errorf("ConfigDescriptor: got %v, want errorInvalidParam", err)
	}
	if _, err := dev.ConfigDescriptorByValue(1); err != ErrorCode(errorInvalidParam) {
		t.Errorf(
			"ConfigDescriptorByValue: got %v, want errorInvalidParam", err,
		)
	}
	if _, err := dev.FindInterfacesByClass(0x07); err != ErrorCode(errorInvalidParam) {
		t.Errorf(
			"FindInterfacesByClass: got %v, want errorInvalidParam", err,
		)
	}
}

func TestDeviceNilInternalPointer(t *testing.T) {
	dev := &Device{}

	if _, err := dev.BusNumber(); err != ErrorCode(errorInvalidParam) {
		t.Errorf("BusNumber: got %v, want errorInvalidParam", err)
	}
	if _, err := dev.Open(); err != ErrorCode(errorInvalidParam) {
		t.Errorf("Open: got %v, want errorInvalidParam", err)
	}
	if _, err := dev.DeviceDescriptor(); err != ErrorCode(errorInvalidParam) {
		t.Errorf("DeviceDescriptor: got %v, want errorInvalidParam", err)
	}
	if _, err := dev.ActiveConfigDescriptor(); err != ErrorCode(errorInvalidParam) {
		t.Errorf(
			"ActiveConfigDescriptor: got %v, want errorInvalidParam", err,
		)
	}
}
