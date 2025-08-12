// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import (
	"testing"
)

func TestEndpointAddressDirection(t *testing.T) {
	testCases := []struct {
		address  endpointAddress
		expected EndpointDirection
	}{
		{0x81, endpointIn},  // IN endpoint (bit 7 set)
		{0x01, endpointOut}, // OUT endpoint (bit 7 clear)
		{0x82, endpointIn},  // IN endpoint 2
		{0x02, endpointOut}, // OUT endpoint 2
		{0xFF, endpointIn},  // IN endpoint 15
		{0x0F, endpointOut}, // OUT endpoint 15
	}

	for _, tc := range testCases {
		result := tc.address.direction()
		if result != tc.expected {
			t.Errorf("endpointAddress(0x%02x).direction() = %d, want %d",
				tc.address, result, tc.expected)
		}
	}
}

func TestEndpointAddressNumber(t *testing.T) {
	testCases := []struct {
		address  endpointAddress
		expected byte
	}{
		{0x00, 0},  // Endpoint 0
		{0x01, 1},  // Endpoint 1
		{0x0F, 15}, // Endpoint 15 (max)
		{0x81, 1},  // IN endpoint 1
		{0x82, 2},  // IN endpoint 2
		{0x8F, 15}, // IN endpoint 15
		{0xFF, 15}, // All bits set, still endpoint 15
	}

	for _, tc := range testCases {
		result := tc.address.endpointNumber()
		if result != tc.expected {
			t.Errorf("endpointAddress(0x%02x).endpointNumber() = %d, want %d",
				tc.address, result, tc.expected)
		}
	}
}

func TestEndpointAttributesTransferType(t *testing.T) {
	testCases := []struct {
		attributes endpointAttributes
		expected   TransferType
	}{
		{0x00, ControlTransfer},     // Control (0b00)
		{0x01, IsochronousTransfer}, // Isochronous (0b01)
		{0x02, BulkTransfer},        // Bulk (0b10)
		{0x03, InterruptTransfer},   // Interrupt (0b11)
		{0xFC, ControlTransfer},     // Other bits set, still control
		{0xFD, IsochronousTransfer}, // Other bits set, still isochronous
		{0xFE, BulkTransfer},        // Other bits set, still bulk
		{0xFF, InterruptTransfer},   // All bits set, still interrupt
	}

	for _, tc := range testCases {
		result := tc.attributes.transferType()
		if result != tc.expected {
			t.Errorf("endpointAttributes(0x%02x).transferType() = %d, want %d",
				tc.attributes, result, tc.expected)
		}
	}
}

func TestEndpointDescriptorMethods(t *testing.T) {
	// Test EndpointDescriptor methods
	desc := &EndpointDescriptor{
		EndpointAddress: 0x81, // IN endpoint 1
		Attributes:      0x02, // Bulk transfer
	}

	// Test Direction method
	if dir := desc.Direction(); dir != endpointIn {
		t.Errorf("EndpointDescriptor.Direction() = %d, want %d", dir, endpointIn)
	}

	// Test Number method
	if num := desc.Number(); num != 1 {
		t.Errorf("EndpointDescriptor.Number() = %d, want 1", num)
	}

	// Test TransferType method
	if tt := desc.TransferType(); tt != BulkTransfer {
		t.Errorf("EndpointDescriptor.TransferType() = %d, want %d", tt, BulkTransfer)
	}
}

func TestEndpointDescriptorEdgeCases(t *testing.T) {
	// Test with various endpoint configurations
	testCases := []struct {
		name         string
		address      endpointAddress
		attributes   endpointAttributes
		expectedDir  EndpointDirection
		expectedNum  byte
		expectedType TransferType
	}{
		{
			name:         "Control endpoint 0 OUT",
			address:      0x00,
			attributes:   0x00,
			expectedDir:  endpointOut,
			expectedNum:  0,
			expectedType: ControlTransfer,
		},
		{
			name:         "Bulk endpoint 1 IN",
			address:      0x81,
			attributes:   0x02,
			expectedDir:  endpointIn,
			expectedNum:  1,
			expectedType: BulkTransfer,
		},
		{
			name:         "Interrupt endpoint 15 OUT",
			address:      0x0F,
			attributes:   0x03,
			expectedDir:  endpointOut,
			expectedNum:  15,
			expectedType: InterruptTransfer,
		},
		{
			name:         "Isochronous endpoint 8 IN",
			address:      0x88,
			attributes:   0x01,
			expectedDir:  endpointIn,
			expectedNum:  8,
			expectedType: IsochronousTransfer,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			desc := &EndpointDescriptor{
				EndpointAddress: tc.address,
				Attributes:      tc.attributes,
			}

			if dir := desc.Direction(); dir != tc.expectedDir {
				t.Errorf("Direction() = %d, want %d", dir, tc.expectedDir)
			}

			if num := desc.Number(); num != tc.expectedNum {
				t.Errorf("Number() = %d, want %d", num, tc.expectedNum)
			}

			if tt := desc.TransferType(); tt != tc.expectedType {
				t.Errorf("TransferType() = %d, want %d", tt, tc.expectedType)
			}
		})
	}
}
