// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import (
	"testing"
)

func TestVidPidToUint32(t *testing.T) {
	testCases := []struct {
		vid      uint16
		pid      uint16
		expected uint32
	}{
		{0x0000, 0x0000, 0x00000000},
		{0x1234, 0x5678, 0x12345678},
		{0xFFFF, 0xFFFF, 0xFFFFFFFF},
		{0x04B8, 0x0202, 0x04B80202},
		{0x0001, 0x0000, 0x00010000},
		{0x0000, 0x0001, 0x00000001},
	}
	for _, tc := range testCases {
		result := vidPidToUint32(tc.vid, tc.pid)
		if result != tc.expected {
			t.Errorf(
				"vidPidToUint32(0x%04X, 0x%04X) = 0x%08X, want 0x%08X",
				tc.vid, tc.pid, result, tc.expected,
			)
		}
	}
}

func TestGetHotplugStorageNilContext(t *testing.T) {
	storage := getHotplugStorage(nil)
	if storage != nil {
		t.Error("getHotplugStorage(nil) should return nil")
	}
}

func TestHotplugRegistryOperations(t *testing.T) {
	// Verify empty registry returns nil
	storage := getHotplugStorage(nil)
	if storage != nil {
		t.Error("expected nil for unregistered context")
	}

	// removeHotplugStorage on non-existent key should not panic
	removeHotplugStorage(nil)
}

func TestHotPlugEventTypeConstants(t *testing.T) {
	if HotplugUndefined != 0 {
		t.Errorf("HotplugUndefined = %d, want 0", HotplugUndefined)
	}
	if HotplugArrived != 1 {
		t.Errorf("HotplugArrived = %d, want 1", HotplugArrived)
	}
	if HotplugLeft != 2 {
		t.Errorf("HotplugLeft = %d, want 2", HotplugLeft)
	}
}

func TestHotPlugEventStruct(t *testing.T) {
	event := HotPlugEvent{
		VendorID:  0x1234,
		ProductID: 0x5678,
		Event:     HotplugArrived,
	}
	if event.VendorID != 0x1234 {
		t.Errorf("VendorID = 0x%04X, want 0x1234", event.VendorID)
	}
	if event.ProductID != 0x5678 {
		t.Errorf("ProductID = 0x%04X, want 0x5678", event.ProductID)
	}
	if event.Event != HotplugArrived {
		t.Errorf("Event = %d, want HotplugArrived", event.Event)
	}
}

func TestHotplugCallbackStorageZeroValue(t *testing.T) {
	storage := &HotplugCallbackStorage{}
	if storage.callbackMap != nil {
		t.Error("zero-value HotplugCallbackStorage should have nil callbackMap")
	}
	if storage.done != nil {
		t.Error("zero-value HotplugCallbackStorage should have nil done channel")
	}
}
