// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"fmt"

	"github.com/gotmc/libusb"
)

func main() {
	version := libusb.GetVersion()
	fmt.Printf(
		"Using libusb version %d.%d.%d (%d)\n",
		version.Major,
		version.Minor,
		version.Micro,
		version.Nano,
	)
	ctx, _ := libusb.Init()
	fmt.Println("Made it past libusb.Init()")
	devices, _ := ctx.GetDeviceList()
	fmt.Printf("Found %v USB devices.\n", len(devices))
	for _, usbDevice := range devices {
		deviceAddress, _ := usbDevice.GetDeviceAddress()
		deviceSpeed, _ := usbDevice.GetDeviceSpeed()
		busNumber, _ := usbDevice.GetBusNumber()
		deviceDescriptor, _ := usbDevice.GetDeviceDescriptor()
		fmt.Printf("Device address %v is on bus number %v. %v\n",
			deviceAddress,
			busNumber,
			deviceSpeed,
		)
		fmt.Printf("\tVendor: %v \tProduct: %v \tClass: %v\n",
			deviceDescriptor.VendorID,
			deviceDescriptor.ProductID,
			deviceDescriptor.DeviceClass,
		)
	}
	fmt.Println("Let's open the Agilent 33220A")
	agilent, _ := ctx.OpenDeviceWithVendorProduct(2391, 1031)
	agilent.Close()
	ctx.Exit()

}
