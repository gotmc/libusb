// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import (
	"testing"
)

func TestDeviceHandleNilChecks(t *testing.T) {
	// Test with completely nil handle
	var dh *DeviceHandle

	// Test StringDescriptor with nil handle
	_, err := dh.StringDescriptor(0, 0)
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("StringDescriptor should return errorInvalidParam for nil handle, got %v", err)
	}

	// Test StringDescriptorASCII with nil handle
	_, err = dh.StringDescriptorASCII(0)
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("StringDescriptorASCII should return errorInvalidParam for nil handle, got %v", err)
	}

	// Test Close with nil handle
	err = dh.Close()
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("Close should return errorInvalidParam for nil handle, got %v", err)
	}

	// Test Configuration with nil handle
	_, err = dh.Configuration()
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("Configuration should return errorInvalidParam for nil handle, got %v", err)
	}

	// Test SetConfiguration with nil handle
	err = dh.SetConfiguration(1)
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("SetConfiguration should return errorInvalidParam for nil handle, got %v", err)
	}

	// Test ClaimInterface with nil handle
	err = dh.ClaimInterface(0)
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("ClaimInterface should return errorInvalidParam for nil handle, got %v", err)
	}

	// Test ReleaseInterface with nil handle
	err = dh.ReleaseInterface(0)
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("ReleaseInterface should return errorInvalidParam for nil handle, got %v", err)
	}

	// Test SetInterfaceAltSetting with nil handle
	err = dh.SetInterfaceAltSetting(0, 0)
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("SetInterfaceAltSetting should return errorInvalidParam for nil handle, got %v", err)
	}

	// Test ResetDevice with nil handle
	err = dh.ResetDevice()
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("ResetDevice should return errorInvalidParam for nil handle, got %v", err)
	}

	// Test KernelDriverActive with nil handle
	_, err = dh.KernelDriverActive(0)
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("KernelDriverActive should return errorInvalidParam for nil handle, got %v", err)
	}

	// Test DetachKernelDriver with nil handle
	err = dh.DetachKernelDriver(0)
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("DetachKernelDriver should return errorInvalidParam for nil handle, got %v", err)
	}

	// Test AttachKernelDriver with nil handle
	err = dh.AttachKernelDriver(0)
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("AttachKernelDriver should return errorInvalidParam for nil handle, got %v", err)
	}

	// Test SetAutoDetachKernelDriver with nil handle
	err = dh.SetAutoDetachKernelDriver(true)
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("SetAutoDetachKernelDriver should return errorInvalidParam for nil handle, got %v", err)
	}
}

func TestDeviceHandleWithNilInternalPointer(t *testing.T) {
	// Test with handle struct but nil internal pointer
	dh := &DeviceHandle{libusbDeviceHandle: nil}

	// Test StringDescriptor
	_, err := dh.StringDescriptor(0, 0)
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("StringDescriptor should return errorInvalidParam for nil internal pointer, got %v", err)
	}

	// Test StringDescriptorASCII
	_, err = dh.StringDescriptorASCII(0)
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("StringDescriptorASCII should return errorInvalidParam for nil internal pointer, got %v", err)
	}

	// Test Close
	err = dh.Close()
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("Close should return errorInvalidParam for nil internal pointer, got %v", err)
	}

	// Test Configuration
	_, err = dh.Configuration()
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("Configuration should return errorInvalidParam for nil internal pointer, got %v", err)
	}
}

func TestDeviceHandleCloseIdempotent(t *testing.T) {
	// Create a handle with non-nil pointer (will be invalid but that's ok for this test)
	dh := &DeviceHandle{}

	// Close should set the pointer to nil
	err := dh.Close()
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("Close should return errorInvalidParam for nil internal pointer, got %v", err)
	}

	// Verify the handle is now nil
	if dh.libusbDeviceHandle != nil {
		t.Error("Close should set internal pointer to nil")
	}

	// Second close should still be safe
	err = dh.Close()
	if err == nil || err != ErrorCode(errorInvalidParam) {
		t.Errorf("Second Close should return errorInvalidParam, got %v", err)
	}
}
