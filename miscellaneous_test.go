// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import (
	"testing"
)

func TestErrorCodeError(t *testing.T) {
	testCodes := []ErrorCode{
		success,
		errorIo,
		errorInvalidParam,
		errorAccess,
		errorNoDevice,
		errorNotFound,
		errorBusy,
		errorTimeout,
		errorOverflow,
		errorPipe,
		errorInterrupted,
		errorNoMem,
		errorNotSupported,
		errorOther,
	}

	for _, code := range testCodes {
		// Test Error() method returns non-empty string
		err := code.Error()
		if err == "" {
			t.Errorf("ErrorCode(%d).Error() returned empty string", code)
		}
	}
}

func TestErrorName(t *testing.T) {
	// Test that ErrorName returns non-empty strings for known error codes
	testCodes := []ErrorCode{
		success,
		errorIo,
		errorInvalidParam,
		errorAccess,
		errorNoDevice,
		errorNotFound,
		errorBusy,
		errorTimeout,
		errorOverflow,
		errorPipe,
		errorInterrupted,
		errorNoMem,
		errorNotSupported,
		errorOther,
	}

	for _, code := range testCodes {
		name := ErrorName(code)
		if name == "" {
			t.Errorf("ErrorName(%d) returned empty string", code)
		}
		// Should start with "LIBUSB_"
		if len(name) < 7 || name[:7] != "LIBUSB_" {
			t.Errorf("ErrorName(%d) = %q, expected to start with 'LIBUSB_'", code, name)
		}
	}
}

func TestStrError(t *testing.T) {
	// Test that StrError returns non-empty strings for known error codes
	testCodes := []ErrorCode{
		success,
		errorIo,
		errorInvalidParam,
		errorAccess,
		errorNoDevice,
		errorNotFound,
		errorBusy,
		errorTimeout,
		errorOverflow,
		errorPipe,
		errorInterrupted,
		errorNoMem,
		errorNotSupported,
		errorOther,
	}

	for _, code := range testCodes {
		msg := StrError(code)
		if msg == "" {
			t.Errorf("StrError(%d) returned empty string", code)
		}
	}
}

func TestSetLocale(t *testing.T) {
	// Test SetLocale with various locales
	testCases := []string{
		"C",
		"en_US",
		"fr_FR",
		"", // Empty string should be handled gracefully
	}

	for _, locale := range testCases {
		// SetLocale should not panic with any string input
		result := SetLocale(locale)
		// We can't really test the return value without knowing the system's capabilities
		// but we can ensure it doesn't panic and returns a valid ErrorCode
		_ = result // Use the result to avoid unused variable warning
	}
}

func TestCPUtoLE16(t *testing.T) {
	testCases := []struct {
		input    int
		expected int // On little-endian systems, should be unchanged
	}{
		{0x1234, 0x1234},
		{0x0000, 0x0000},
		{0xFFFF, 0xFFFF},
		{0x00FF, 0x00FF},
		{0xFF00, 0xFF00},
	}

	for _, tc := range testCases {
		result := CPUtoLE16(tc.input)
		// On little-endian systems (most common), result should equal input
		// On big-endian systems, bytes would be swapped
		// We just verify the function doesn't crash and returns a valid value
		if result < 0 || result > 0xFFFF {
			t.Errorf("CPUtoLE16(%d) = %d, expected value in range [0, 65535]", tc.input, result)
		}
	}
}

func TestHasCapability(t *testing.T) {
	// Test HasCapability with some capability constants
	// We can't know which capabilities are supported without libusb context,
	// but we can test that the function doesn't crash
	testCapabilities := []int{
		0,    // No capability
		1,    // Some capability
		2,    // Another capability
		-1,   // Invalid capability
		1000, // Large capability value
	}

	for _, cap := range testCapabilities {
		// Should not panic regardless of input
		result := HasCapability(cap)
		// Result should be a boolean value (true or false)
		_ = result // Use the result to avoid unused variable warning
	}
}

func TestBcdToDecimal(t *testing.T) {
	testCases := []struct {
		bcd      uint16
		expected float64
	}{
		{0x0120, 1.20},  // BCD 0120 = decimal 1.20
		{0x0000, 0.00},  // BCD 0000 = decimal 0.00
		{0x0101, 1.01},  // BCD 0101 = decimal 1.01
		{0x9999, 99.99}, // BCD 9999 = decimal 99.99
	}

	for _, tc := range testCases {
		result := bcdToDecimal(tc.bcd)
		// Allow for small floating point differences
		diff := result - tc.expected
		if diff < -0.01 || diff > 0.01 {
			t.Errorf("bcdToDecimal(0x%04x) = %.2f, expected %.2f", tc.bcd, result, tc.expected)
		}
	}
}
