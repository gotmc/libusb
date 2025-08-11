// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import (
	"testing"
)

// Note: Most syncio functions call libusb C functions directly and will segfault
// with nil or invalid pointers. Testing error conditions would require valid
// libusb context and device setup, which is beyond unit testing scope.

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
			t.Errorf("TransferDirection(%d).String() = %q, want %q", tc.direction, result, tc.expected)
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
			t.Errorf("RequestRecipient(%d).String() = %q, want %q", tc.recipient, result, tc.expected)
		}
	}
}
