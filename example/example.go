// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"fmt"
	"log"

	"github.com/gotmc/libusb"
)

func showVersion() {
	version := libusb.GetVersion()
	fmt.Printf(
		"Using libusb version %d.%d.%d (%d)\n",
		version.Major,
		version.Minor,
		version.Micro,
		version.Nano,
	)
}

func main() {
	showVersion()
	ctx, err := libusb.Init()
	if err != nil {
		log.Fatal("Couldn't create USB context. Ending now.")
	}
	defer ctx.Exit()
	devices, _ := ctx.GetDeviceList()
	fmt.Printf("Found %v USB devices.\n", len(devices))
	for _, usbDevice := range devices {
		deviceAddress, _ := usbDevice.GetDeviceAddress()
		deviceSpeed, _ := usbDevice.GetDeviceSpeed()
		busNumber, _ := usbDevice.GetBusNumber()
		usbDeviceDescriptor, _ := usbDevice.GetDeviceDescriptor()
		fmt.Printf("Device address %v is on bus number %v\n=> %v\n",
			deviceAddress,
			busNumber,
			deviceSpeed,
		)
		fmt.Printf("=> Vendor: %v \tProduct: %v\n=> Class: %v\n",
			usbDeviceDescriptor.VendorID,
			usbDeviceDescriptor.ProductID,
			usbDeviceDescriptor.DeviceClass,
		)
		fmt.Printf("=> USB: %v\tMax Packet 0: %v\tSN Index: %v\n",
			usbDeviceDescriptor.USBSpecification,
			usbDeviceDescriptor.MaxPacketSize0,
			usbDeviceDescriptor.SerialNumberIndex,
		)
	}
	showInfo(ctx, "Agilent 33220A", 2391, 1031)
	// showInfo(ctx, "Nike SportWatch", 4524, 21588)
	// showInfo(ctx, "Nike FuelBand", 4524, 25957)

}

func showInfo(ctx *libusb.Context, name string, vendorID, productID uint16) {
	fmt.Printf("Let's open the %s using the Vendor and Product IDs\n", name)
	usbDevice, usbDeviceHandle, err := ctx.OpenDeviceWithVendorProduct(vendorID, productID)
	usbDeviceDescriptor, _ := usbDevice.GetDeviceDescriptor()
	if err != nil {
		fmt.Printf("=> Failed opening the %s: %v\n", name, err)
		return
	}
	defer usbDeviceHandle.Close()
	serialnum, _ := usbDeviceHandle.GetStringDescriptorASCII(
		usbDeviceDescriptor.SerialNumberIndex,
	)
	manufacturer, _ := usbDeviceHandle.GetStringDescriptorASCII(
		usbDeviceDescriptor.ManufacturerIndex)
	product, _ := usbDeviceHandle.GetStringDescriptorASCII(
		usbDeviceDescriptor.ProductIndex)
	fmt.Printf("Found %v %v S/N %s using Vendor ID %v and Product ID %v\n",
		manufacturer,
		product,
		serialnum,
		vendorID,
		productID,
	)
	configDescriptor, err := usbDevice.GetActiveConfigDescriptor()
	if err != nil {
		log.Fatalf("Failed getting the active config: %v", err)
	}
	fmt.Printf("=> Max Power = %d mA\n",
		configDescriptor.MaxPowerMilliAmperes)
	var singularPlural string
	if configDescriptor.NumInterfaces == 1 {
		singularPlural = "interface"
	} else {
		singularPlural = "interfaces"
	}
	fmt.Printf("=> Found %d %s\n",
		configDescriptor.NumInterfaces, singularPlural)
	fmt.Printf("=> The first interface has %d alternate settings.\n",
		configDescriptor.SupportedInterfaces[0].NumAltSettings)
	firstDescriptor := configDescriptor.SupportedInterfaces[0].InterfaceDescriptors[0]
	fmt.Printf("=> The first interface descriptor has a length of %d.\n", firstDescriptor.Length)
	fmt.Printf("=> The first interface descriptor is interface number %d.\n", firstDescriptor.InterfaceNumber)
	fmt.Printf("=> The first interface descriptor has %d endpoint(s).\n", firstDescriptor.NumEndpoints)
	fmt.Printf("   => USB-IF class %d, subclass %d, protocol %d.\n",
		firstDescriptor.InterfaceClass, firstDescriptor.InterfaceSubClass, firstDescriptor.InterfaceProtocol)
}
