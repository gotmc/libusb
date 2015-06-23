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

func main() {
	version := libusb.GetVersion()
	fmt.Printf(
		"Using libusb version %d.%d.%d (%d)\n",
		version.Major,
		version.Minor,
		version.Micro,
		version.Nano,
	)
	ctx, err := libusb.Init()
	if err != nil {
		log.Fatal("Couldn't create USB context. Ending now.")
	}
	defer ctx.Exit()
	fmt.Println("Made it past libusb.Init()")
	devices, _ := ctx.GetDeviceList()
	fmt.Printf("Found %v USB devices.\n", len(devices))
	for _, usbDevice := range devices {
		deviceAddress, _ := usbDevice.GetDeviceAddress()
		deviceSpeed, _ := usbDevice.GetDeviceSpeed()
		busNumber, _ := usbDevice.GetBusNumber()
		_ = usbDevice.GetDeviceDescriptor()
		fmt.Printf("Device address %v is on bus number %v. %v\n",
			deviceAddress,
			busNumber,
			deviceSpeed,
		)
		fmt.Printf("\tVendor: %v \tProduct: %v \tClass: %v\n",
			usbDevice.VendorID,
			usbDevice.ProductID,
			usbDevice.DeviceClass,
		)
		fmt.Printf("\tUSB: %v\tRelease Num: %v\tSN Index: %v\n",
			usbDevice.USBSpecification,
			usbDevice.DeviceReleaseNumber,
			usbDevice.SerialNumberIndex,
		)
	}
	fmt.Println("Let's open the Agilent 33220A")
	agilent, err := ctx.OpenDeviceWithVendorProduct(2391, 1031)
	if err != nil {
		fmt.Println("Couldn't find the Agilent 33220A")
	} else {
		defer agilent.Close()
		serialnum, _ := agilent.GetStringDescriptorASCII(
			agilent.SerialNumberIndex,
		)
		manufacturer, _ := agilent.GetStringDescriptorASCII(
			agilent.ManufacturerIndex)
		product, _ := agilent.GetStringDescriptorASCII(
			agilent.ProductIndex)
		fmt.Printf("Found %v %v S/N %s\n",
			manufacturer,
			product,
			serialnum,
		)
	}

}
