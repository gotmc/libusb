// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import (
	"testing"
)

const (
	failCheck = `✗` // UTF-8 u2717
	passCheck = `✓` // UTF-8 u2713
)

func TestBcdType(t *testing.T) {
	testCases := []struct {
		bcdValue   uint16
		bcdString  string
		bcdDecimal float64
	}{
		{0x0300, "0x0300 (3.00)", 3.0},
		{0x0200, "0x0200 (2.00)", 2.0},
		{0x0110, "0x0110 (1.10)", 1.1},
		{0x0100, "0x0100 (1.00)", 1.0},
	}
	t.Log("Given the need to test the bcd type")
	{
		for _, testCase := range testCases {
			b := bcd(testCase.bcdValue)
			t.Logf("\tWhen getting the string for bcd %#04x", testCase.bcdValue)
			computedString := b.String()
			if computedString != testCase.bcdString {
				t.Errorf(
					"\t%v Should have computed %s but got %s",
					failCheck,
					testCase.bcdString,
					computedString,
				)
			} else {
				t.Logf("\t%v Should compute %s", passCheck, computedString)
			}
		}
	}
}

func TestClassCodeStringMethod(t *testing.T) {
	testCases := []struct {
		class    classCode
		expected string
	}{
		{
			perInterface,
			"Each interface specifies its own class information and all interfaces operate independently.",
		},
		{audio, "Audio class."},
		{comm, "Communications class."},
		{hid, "Human Interface Device class."},
		{physical, "Physical."},
		{printer, "Printer class."},
		{image, "Image class."},
		{massStorage, "Mass storage class."},
		{hub, "Hub class."},
		{data, "Data class."},
		{smartCard, "Smart Card."},
		{contentSecurity, "Content Security."},
		{video, "Video."},
		{personalHealthcare, "Personal Healthcare."},
		{diagnosticDevice, "Diagnostic Device."},
		{wireless, "Wireless class."},
		{application, "Application class."},
		{vendorSpec, "Class is vendor-specific."},
	}
	t.Log("Given the need to test the classCode.String() method")
	{
		for _, testCase := range testCases {
			t.Logf("\tWhen getting classCode %d's string", testCase.class)
			computed := testCase.class.String()
			if computed != testCase.expected {
				t.Errorf(
					"\t%v Should have yielded: %s, but got %s",
					failCheck,
					testCase.expected,
					computed,
				)
			} else {
				t.Logf("\t%v Should yield: %s", passCheck, computed)
			}
		}
	}
}

func TestDescriptortypeStringMethod(t *testing.T) {
	testCases := []struct {
		desc     descriptorType
		expected string
	}{
		{descDevice, "Device descriptor."},
		{descConfig, "Configuration descriptor."},
		{descString, "String descriptor."},
		{descInterface, "Interface descriptor."},
		{descEndpoint, "Endpoint descriptor."},
		{descBos, "BOS descriptor."},
		{descDeviceCapability, "Device Capability descriptor."},
		{descHid, "HID descriptor."},
		{descReport, "HID report descriptor."},
		{descPhysical, "Physical descriptor."},
		{descHub, "Hub descriptor."},
		{descSuperspeedHub, "SuperSpeed Hub descriptor."},
		{descEndpointCompanion, "SuperSpeed Endpoint Companion descriptor."},
	}
	t.Log("Given the need to test the descriptorType.String() method")
	{
		for _, testCase := range testCases {
			t.Logf("\tWhen getting descriptorType %d's string", testCase.desc)
			computed := testCase.desc.String()
			if computed != testCase.expected {
				t.Errorf(
					"\t%v Should have yielded: %s, but got %s",
					failCheck,
					testCase.expected,
					computed,
				)
			} else {
				t.Logf("\t%v Should yield: %s", passCheck, computed)
			}
		}
	}
}

func TestEndpointDirectionStringMethod(t *testing.T) {
	testCases := []struct {
		end      EndpointDirection
		expected string
	}{
		{endpointOut, "Out: host-to-device."},
		{endpointIn, "In: device-to-host."},
	}
	t.Log("Given the need to test the endpointDirection.String() method")
	{
		for _, testCase := range testCases {
			t.Logf("\tWhen getting endpointDirection %d's string", testCase.end)
			computed := testCase.end.String()
			if computed != testCase.expected {
				t.Errorf(
					"\t%v Should have yielded: %s, but got %s",
					failCheck,
					testCase.expected,
					computed,
				)
			} else {
				t.Logf("\t%v Should yield: %s", passCheck, computed)
			}
		}
	}
}

func TestTransferTypeStringMethod(t *testing.T) {
	testCases := []struct {
		transfer TransferType
		expected string
	}{
		{ControlTransfer, "Control endpoint."},
		{IsochronousTransfer, "Isochronous endpoint."},
		{BulkTransfer, "Bulk endpoint."},
		{InterruptTransfer, "Interrupt endpoint."},
	}
	t.Log("Given the need to test the endpointDirection.String() method")
	{
		for _, testCase := range testCases {
			t.Logf("\tWhen getting endpointDirection %d's string", testCase.transfer)
			computed := testCase.transfer.String()
			if computed != testCase.expected {
				t.Errorf(
					"\t%v Should have yielded: %s, but got %s",
					failCheck,
					testCase.expected,
					computed,
				)
			} else {
				t.Logf("\t%v Should yield: %s", passCheck, computed)
			}
		}
	}
}

func TestInterfaceDescriptorsGetInterfacesByClass(t *testing.T) {
	// Create test interface descriptors
	printerIface := &InterfaceDescriptor{
		InterfaceClass:  0x07, // Printer class
		InterfaceNumber: 0,
	}
	audioIface := &InterfaceDescriptor{
		InterfaceClass:  0x01, // Audio class
		InterfaceNumber: 1,
	}
	anotherPrinterIface := &InterfaceDescriptor{
		InterfaceClass:  0x07, // Printer class
		InterfaceNumber: 2,
	}

	ifaces := InterfaceDescriptors{printerIface, audioIface, anotherPrinterIface}

	// Test finding printer interfaces
	printerIfaces := ifaces.GetInterfacesByClass(0x07)
	if len(printerIfaces) != 2 {
		t.Errorf("Expected 2 printer interfaces, got %d", len(printerIfaces))
	}

	// Test finding audio interfaces
	audioIfaces := ifaces.GetInterfacesByClass(0x01)
	if len(audioIfaces) != 1 {
		t.Errorf("Expected 1 audio interface, got %d", len(audioIfaces))
	}

	// Test finding non-existent class
	massStorageIfaces := ifaces.GetInterfacesByClass(0x08)
	if len(massStorageIfaces) != 0 {
		t.Errorf("Expected 0 mass storage interfaces, got %d", len(massStorageIfaces))
	}
}

func TestSupportedInterfacesGetAllInterfacesByClass(t *testing.T) {
	// Create test supported interfaces
	printerIface := &InterfaceDescriptor{
		InterfaceClass:  0x07, // Printer class
		InterfaceNumber: 0,
	}
	audioIface := &InterfaceDescriptor{
		InterfaceClass:  0x01, // Audio class
		InterfaceNumber: 1,
	}

	supportedIface1 := &SupportedInterface{
		InterfaceDescriptors: InterfaceDescriptors{printerIface},
	}
	supportedIface2 := &SupportedInterface{
		InterfaceDescriptors: InterfaceDescriptors{audioIface},
	}

	supportedInterfaces := SupportedInterfaces{supportedIface1, supportedIface2}

	// Test finding printer interfaces
	printerIfaces := supportedInterfaces.GetAllInterfacesByClass(0x07)
	if len(printerIfaces) != 1 {
		t.Errorf("Expected 1 printer interface, got %d", len(printerIfaces))
	}

	// Test finding non-existent class
	massStorageIfaces := supportedInterfaces.GetAllInterfacesByClass(0x08)
	if len(massStorageIfaces) != 0 {
		t.Errorf("Expected 0 mass storage interfaces, got %d", len(massStorageIfaces))
	}
}
