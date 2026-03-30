// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import (
	"testing"
)

func TestBulkTransferNilHandle(t *testing.T) {
	var dh *DeviceHandle
	_, err := dh.BulkTransfer(0, nil, 0, 0)
	if err != ErrorCode(errorInvalidParam) {
		t.Errorf("BulkTransfer: got %v, want errorInvalidParam", err)
	}
}

func TestBulkTransferNilInternalPointer(t *testing.T) {
	dh := &DeviceHandle{}
	_, err := dh.BulkTransfer(0, nil, 0, 0)
	if err != ErrorCode(errorInvalidParam) {
		t.Errorf("BulkTransfer: got %v, want errorInvalidParam", err)
	}
}

func TestBulkTransferOutNilHandle(t *testing.T) {
	var dh *DeviceHandle
	_, err := dh.BulkTransferOut(0, nil, 0)
	if err != ErrorCode(errorInvalidParam) {
		t.Errorf("BulkTransferOut: got %v, want errorInvalidParam", err)
	}
}

func TestBulkTransferInNilHandle(t *testing.T) {
	var dh *DeviceHandle
	_, _, err := dh.BulkTransferIn(0, 64, 0)
	if err != ErrorCode(errorInvalidParam) {
		t.Errorf("BulkTransferIn: got %v, want errorInvalidParam", err)
	}
}

func TestControlTransferNilHandle(t *testing.T) {
	var dh *DeviceHandle
	_, err := dh.ControlTransfer(0, 0, 0, 0, nil, 0, 0)
	if err != ErrorCode(errorInvalidParam) {
		t.Errorf("ControlTransfer: got %v, want errorInvalidParam", err)
	}
}

func TestControlTransferWithTypesNilHandle(t *testing.T) {
	var dh *DeviceHandle
	_, err := dh.ControlTransferWithTypes(
		HostToDevice, Standard, DeviceRecipient,
		0, 0, 0, nil, 0, 0,
	)
	if err != ErrorCode(errorInvalidParam) {
		t.Errorf(
			"ControlTransferWithTypes: got %v, want errorInvalidParam", err,
		)
	}
}

func TestControlOutNilHandle(t *testing.T) {
	var dh *DeviceHandle
	_, err := dh.ControlOut(Standard, DeviceRecipient, 0, 0, 0, nil, 0)
	if err != ErrorCode(errorInvalidParam) {
		t.Errorf("ControlOut: got %v, want errorInvalidParam", err)
	}
}

func TestControlInNilHandle(t *testing.T) {
	var dh *DeviceHandle
	_, err := dh.ControlIn(Standard, DeviceRecipient, 0, 0, 0, nil, 64, 0)
	if err != ErrorCode(errorInvalidParam) {
		t.Errorf("ControlIn: got %v, want errorInvalidParam", err)
	}
}

func TestInterruptTransferNilHandle(t *testing.T) {
	var dh *DeviceHandle
	_, err := dh.InterruptTransfer(0, nil, 0, 0)
	if err != ErrorCode(errorInvalidParam) {
		t.Errorf("InterruptTransfer: got %v, want errorInvalidParam", err)
	}
}

func TestInterruptTransferNilInternalPointer(t *testing.T) {
	dh := &DeviceHandle{}
	_, err := dh.InterruptTransfer(0, nil, 0, 0)
	if err != ErrorCode(errorInvalidParam) {
		t.Errorf("InterruptTransfer: got %v, want errorInvalidParam", err)
	}
}

func TestBitmapRequestType(t *testing.T) {
	// Test building request type bitmap
	testCases := []struct {
		direction TransferDirection
		reqType   RequestType
		recipient RequestRecipient
		expected  byte
	}{
		{DeviceToHost, Standard, DeviceRecipient, 0x80},
		{HostToDevice, Standard, DeviceRecipient, 0x00},
		{DeviceToHost, Class, InterfaceRecipient, 0xA1},
		{HostToDevice, Vendor, EndpointRecipient, 0x42},
	}

	for _, tc := range testCases {
		result := BitmapRequestType(tc.direction, tc.reqType, tc.recipient)
		if result != tc.expected {
			t.Errorf("BitmapRequestType(%v, %v, %v) = 0x%02x, want 0x%02x",
				tc.direction, tc.reqType, tc.recipient, result, tc.expected)
		}
	}
}

func TestTransferDirectionString(t *testing.T) {
	testCases := []struct {
		direction TransferDirection
		expected  string
	}{
		{HostToDevice, "Host-to-device"},
		{DeviceToHost, "Device-to-host"},
	}

	for _, tc := range testCases {
		result := tc.direction.String()
		if result != tc.expected {
			t.Errorf(
				"TransferDirection(%d).String() = %q, want %q",
				tc.direction, result, tc.expected,
			)
		}
	}
}

func TestRequestTypeString(t *testing.T) {
	testCases := []struct {
		reqType  RequestType
		expected string
	}{
		{Standard, "Standard"},
		{Class, "Class"},
		{Vendor, "Vendor"},
		{Reserved, "Reserved"},
	}

	for _, tc := range testCases {
		result := tc.reqType.String()
		if result != tc.expected {
			t.Errorf("RequestType(%d).String() = %q, want %q", tc.reqType, result, tc.expected)
		}
	}
}

func TestRequestRecipientString(t *testing.T) {
	testCases := []struct {
		recipient RequestRecipient
		expected  string
	}{
		{DeviceRecipient, "Device"},
		{InterfaceRecipient, "Interface"},
		{EndpointRecipient, "Endpoint"},
		{OtherRecipient, "Other"},
	}

	for _, tc := range testCases {
		result := tc.recipient.String()
		if result != tc.expected {
			t.Errorf(
				"RequestRecipient(%d).String() = %q, want %q",
				tc.recipient, result, tc.expected,
			)
		}
	}
}
