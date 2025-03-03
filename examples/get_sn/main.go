// Copyright (c) 2015–2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"fmt"
	"log"

	libusb "github.com/gotmc/libusb/v2"
)

func showVersion() {
	version := libusb.Version()
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
	ctx, err := libusb.NewContext()
	if err != nil {
		log.Fatal("Couldn't create USB context. Ending now.")
	}
	defer ctx.Close()
	devices, err := ctx.DeviceList()
	if err != nil {
		log.Fatalf("Couldn't get devices")
	}
	log.Printf("Found %v USB devices.\n", len(devices))
	for _, device := range devices {
		usbDeviceDescriptor, err := device.DeviceDescriptor()
		if err != nil {
			log.Printf("Error getting device descriptor: %s", err)
			continue
		}
		addr, err := device.DeviceAddress()
		if err != nil {
			log.Printf("Error getting device address: %s", err)
			continue
		}
		log.Printf("Device address: %d", addr)
		handle, err := device.Open()
		if err != nil {
			log.Printf("Error opening device: %s", err)
			continue
		}
		defer handle.Close()
		serialNumber, err := handle.StringDescriptorASCII(usbDeviceDescriptor.SerialNumberIndex)
		if err != nil {
			serialNumber = "N/A"
		}
		manufacturer, err := handle.StringDescriptorASCII(usbDeviceDescriptor.ManufacturerIndex)
		if err != nil {
			manufacturer = "N/A"
		}
		product, err := handle.StringDescriptorASCII(usbDeviceDescriptor.ProductIndex)
		if err != nil {
			product = "N/A"
		}
		log.Printf("Found %s (S/N: %s) manufactured by %s",
			product,
			serialNumber,
			manufacturer,
		)
	}

}
