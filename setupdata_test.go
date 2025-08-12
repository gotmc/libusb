// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import (
	"testing"
)

func TestTransferDirectionStringer(t *testing.T) {
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

func TestUnknownTransferDirection(t *testing.T) {
	// Test an unknown transfer direction
	unknown := TransferDirection(0x55)
	result := unknown.String()
	// Should return empty string for unknown directions
	if result != "" {
		t.Errorf("Unknown transfer direction should return empty string, got %q", result)
	}
}

func TestRequestTypeStringer(t *testing.T) {
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

func TestUnknownRequestType(t *testing.T) {
	// Test an unknown request type
	unknown := RequestType(0x55)
	result := unknown.String()
	// Should return empty string for unknown types
	if result != "" {
		t.Errorf("Unknown request type should return empty string, got %q", result)
	}
}

func TestRequestRecipientStringer(t *testing.T) {
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

func TestUnknownRequestRecipient(t *testing.T) {
	// Test an unknown request recipient
	unknown := RequestRecipient(0x55)
	result := unknown.String()
	// Should return empty string for unknown recipients
	if result != "" {
		t.Errorf("Unknown request recipient should return empty string, got %q", result)
	}
}

func TestBitmapRequestTypeComprehensive(t *testing.T) {
	// Test all combinations of direction, type, and recipient
	testCases := []struct {
		name      string
		direction TransferDirection
		reqType   RequestType
		recipient RequestRecipient
		expected  byte
	}{
		// Host-to-device combinations
		{"H2D Standard Device", HostToDevice, Standard, DeviceRecipient, 0x00},
		{"H2D Standard Interface", HostToDevice, Standard, InterfaceRecipient, 0x01},
		{"H2D Standard Endpoint", HostToDevice, Standard, EndpointRecipient, 0x02},
		{"H2D Standard Other", HostToDevice, Standard, OtherRecipient, 0x03},
		{"H2D Class Device", HostToDevice, Class, DeviceRecipient, 0x20},
		{"H2D Class Interface", HostToDevice, Class, InterfaceRecipient, 0x21},
		{"H2D Class Endpoint", HostToDevice, Class, EndpointRecipient, 0x22},
		{"H2D Class Other", HostToDevice, Class, OtherRecipient, 0x23},
		{"H2D Vendor Device", HostToDevice, Vendor, DeviceRecipient, 0x40},
		{"H2D Vendor Interface", HostToDevice, Vendor, InterfaceRecipient, 0x41},
		{"H2D Vendor Endpoint", HostToDevice, Vendor, EndpointRecipient, 0x42},
		{"H2D Vendor Other", HostToDevice, Vendor, OtherRecipient, 0x43},

		// Device-to-host combinations
		{"D2H Standard Device", DeviceToHost, Standard, DeviceRecipient, 0x80},
		{"D2H Standard Interface", DeviceToHost, Standard, InterfaceRecipient, 0x81},
		{"D2H Standard Endpoint", DeviceToHost, Standard, EndpointRecipient, 0x82},
		{"D2H Standard Other", DeviceToHost, Standard, OtherRecipient, 0x83},
		{"D2H Class Device", DeviceToHost, Class, DeviceRecipient, 0xA0},
		{"D2H Class Interface", DeviceToHost, Class, InterfaceRecipient, 0xA1},
		{"D2H Class Endpoint", DeviceToHost, Class, EndpointRecipient, 0xA2},
		{"D2H Class Other", DeviceToHost, Class, OtherRecipient, 0xA3},
		{"D2H Vendor Device", DeviceToHost, Vendor, DeviceRecipient, 0xC0},
		{"D2H Vendor Interface", DeviceToHost, Vendor, InterfaceRecipient, 0xC1},
		{"D2H Vendor Endpoint", DeviceToHost, Vendor, EndpointRecipient, 0xC2},
		{"D2H Vendor Other", DeviceToHost, Vendor, OtherRecipient, 0xC3},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := BitmapRequestType(tc.direction, tc.reqType, tc.recipient)
			if result != tc.expected {
				t.Errorf("BitmapRequestType(%v, %v, %v) = 0x%02x, want 0x%02x",
					tc.direction, tc.reqType, tc.recipient, result, tc.expected)
			}
		})
	}
}

func TestTransferDirectionConstants(t *testing.T) {
	// Verify the constant values are correct
	if HostToDevice != 0x00 {
		t.Errorf("HostToDevice = 0x%02x, want 0x00", HostToDevice)
	}
	if DeviceToHost != 0x80 {
		t.Errorf("DeviceToHost = 0x%02x, want 0x80", DeviceToHost)
	}
}
