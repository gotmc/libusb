// Copyright (c) 2016 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"log"

	"github.com/gotmc/libusb"
)

func main() {
	ctx, err := libusb.NewContext()
	if err != nil {
		log.Fatal("Couldn't create USB context. Ending now.")
	}
	defer ctx.Close()
	devices, err := ctx.GetDeviceList()
	if err != nil {
		log.Fatalf("Couldn't get devices")
	}
	log.Printf("Found %v USB devices.\n", len(devices))
	for _, device := range devices {
		usbDeviceDescriptor, err := device.GetDeviceDescriptor()
		if err != nil {
			log.Printf("Error getting device descriptor: %s", err)
			continue
		}
		handle, err := device.Open()
		if err != nil {
			log.Printf("Error opening device: %s", err)
			continue
		}
		defer handle.Close()
		serialNumber, err := handle.GetStringDescriptorASCII(usbDeviceDescriptor.SerialNumberIndex)
		if err != nil {
			serialNumber = "N/A"
		}
		manufacturer, err := handle.GetStringDescriptorASCII(usbDeviceDescriptor.ManufacturerIndex)
		if err != nil {
			manufacturer = "N/A"
		}
		product, err := handle.GetStringDescriptorASCII(usbDeviceDescriptor.ProductIndex)
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
