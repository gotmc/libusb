// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import (
	"testing"
)

func TestConfigDescriptorStructure(t *testing.T) {
	// Test that ConfigDescriptor struct can be created and accessed
	config := &ConfigDescriptor{
		Length:               9,
		DescriptorType:       descConfig,
		TotalLength:          100,
		NumInterfaces:        2,
		ConfigurationValue:   1,
		ConfigurationIndex:   0,
		Attributes:           0x80,
		MaxPowerMilliAmperes: 500,
		SupportedInterfaces:  nil,
	}

	// Verify basic fields
	if config.Length != 9 {
		t.Errorf("ConfigDescriptor.Length = %d, want 9", config.Length)
	}
	if config.NumInterfaces != 2 {
		t.Errorf("ConfigDescriptor.NumInterfaces = %d, want 2", config.NumInterfaces)
	}
	if config.MaxPowerMilliAmperes != 500 {
		t.Errorf("ConfigDescriptor.MaxPowerMilliAmperes = %d, want 500", config.MaxPowerMilliAmperes)
	}
}

func TestConfigStructure(t *testing.T) {
	// Test that Config struct properly embeds ConfigDescriptor
	desc := &ConfigDescriptor{
		Length:               9,
		ConfigurationValue:   1,
		MaxPowerMilliAmperes: 250,
	}

	config := &Config{
		ConfigDescriptor: desc,
		Device:           nil,
	}

	// Verify embedded fields are accessible
	if config.Length != 9 {
		t.Errorf("Config.Length = %d, want 9", config.Length)
	}
	if config.ConfigurationValue != 1 {
		t.Errorf("Config.ConfigurationValue = %d, want 1", config.ConfigurationValue)
	}
	if config.MaxPowerMilliAmperes != 250 {
		t.Errorf("Config.MaxPowerMilliAmperes = %d, want 250", config.MaxPowerMilliAmperes)
	}
}

func TestConfigDescriptorAttributes(t *testing.T) {
	testCases := []struct {
		name       string
		attributes uint8
		selfPower  bool
		remoteWake bool
	}{
		{"Bus powered only", 0x80, false, false},
		{"Self powered", 0xC0, true, false},
		{"Remote wakeup", 0xA0, false, true},
		{"Self powered + Remote wakeup", 0xE0, true, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			config := &ConfigDescriptor{
				Attributes: tc.attributes,
			}

			// Test self-powered bit (bit 6)
			selfPowered := (config.Attributes & 0x40) != 0
			if selfPowered != tc.selfPower {
				t.Errorf("Self-powered bit = %v, want %v", selfPowered, tc.selfPower)
			}

			// Test remote wakeup bit (bit 5)
			remoteWakeup := (config.Attributes & 0x20) != 0
			if remoteWakeup != tc.remoteWake {
				t.Errorf("Remote wakeup bit = %v, want %v", remoteWakeup, tc.remoteWake)
			}

			// Bit 7 should always be set for USB 2.0
			if (config.Attributes & 0x80) == 0 {
				t.Error("Bit 7 should be set for USB 2.0 compliance")
			}
		})
	}
}

func TestMaxPowerConversion(t *testing.T) {
	testCases := []struct {
		maxPowerRaw uint8 // Raw value from descriptor (in 2mA units)
		expectedMA  uint  // Expected value in milliamps
	}{
		{0, 0},     // 0 mA
		{1, 2},     // 2 mA
		{50, 100},  // 100 mA
		{125, 250}, // 250 mA
		{250, 500}, // 500 mA (max for USB 2.0)
	}

	for _, tc := range testCases {
		// The ConfigDescriptor should store milliamps directly
		// The conversion from 2mA units happens during descriptor parsing
		config := &ConfigDescriptor{
			MaxPowerMilliAmperes: uint(tc.maxPowerRaw) * 2,
		}

		if config.MaxPowerMilliAmperes != tc.expectedMA {
			t.Errorf("MaxPowerMilliAmperes = %d, want %d (from raw value %d)",
				config.MaxPowerMilliAmperes, tc.expectedMA, tc.maxPowerRaw)
		}
	}
}
