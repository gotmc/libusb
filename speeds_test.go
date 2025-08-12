// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import (
	"testing"
)

func TestSpeedTypeString(t *testing.T) {
	testCases := []struct {
		speed    SpeedType
		expected string
	}{
		{speedUnknown, "The OS doesn't report or know the device speed."},
		{speedLow, "The device is operating at low speed (1.5MBit/s)"},
		{speedFull, "The device is operating at full speed (12MBit/s)"},
		{speedHigh, "The device is operating at high speed (480MBit/s)"},
		{speedSuper, "The device is operating at super speed (5000MBit/s)"},
	}

	for _, tc := range testCases {
		result := tc.speed.String()
		if result != tc.expected {
			t.Errorf("SpeedType(%d).String() = %q, want %q", tc.speed, result, tc.expected)
		}
	}
}

func TestUnknownSpeedType(t *testing.T) {
	// Test an unknown speed type
	unknown := SpeedType(255)
	result := unknown.String()
	// Should return empty string for unknown speed types
	if result != "" {
		t.Errorf("Unknown speed type should return empty string, got %q", result)
	}
}

func TestSpeedTypeConstants(t *testing.T) {
	// Verify that speed constants have expected values
	// These are based on libusb constants
	testCases := []struct {
		speed SpeedType
		name  string
	}{
		{speedUnknown, "speedUnknown"},
		{speedLow, "speedLow"},
		{speedFull, "speedFull"},
		{speedHigh, "speedHigh"},
		{speedSuper, "speedSuper"},
	}

	for _, tc := range testCases {
		// Just verify the constants are defined and have non-negative values
		if tc.speed < 0 {
			t.Errorf("%s has invalid value: %d", tc.name, tc.speed)
		}
	}
}

func TestSupportedSpeedString(t *testing.T) {
	testCases := []struct {
		speed    supportedSpeed
		expected string
	}{
		{lowSpeedOperation, "Low speed operation supported (1.5MBit/s)."},
		{fullSpeedOperation, "Full speed operation supported (12MBit/s)."},
		{highSpeedOperation, "High speed operation supported (480MBit/s)."},
		{superSpeedOperation, "Superspeed operation supported (5000MBit/s)."},
	}

	for _, tc := range testCases {
		result := tc.speed.String()
		if result != tc.expected {
			t.Errorf("supportedSpeed(%d).String() = %q, want %q", tc.speed, result, tc.expected)
		}
	}
}

func TestUnknownSupportedSpeed(t *testing.T) {
	// Test an unknown supported speed
	unknown := supportedSpeed(255)
	result := unknown.String()
	// Should return empty string for unknown speed types
	if result != "" {
		t.Errorf("Unknown supported speed should return empty string, got %q", result)
	}
}
