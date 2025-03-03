// Copyright (c) 2015–2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"log"

	libusb "github.com/gotmc/libusb/v2"
)

func main() {
	ctx, err := libusb.NewContext()
	if err != nil {
		log.Fatal("Couldn't create USB context. Ending now.")
	}
	defer ctx.Close()
	devices, err := ctx.DeviceList()
	if err != nil {
		log.Fatalf("Couldn't get devices")
	}
	for _, device := range devices {
		usbDeviceDescriptor, _ := device.DeviceDescriptor()
		handle, err := device.Open()
		if err != nil {
			log.Fatalf("Error opening device: %s", err)
		}
		defer handle.Close()
		serialNumber, err := handle.StringDescriptorASCII(usbDeviceDescriptor.SerialNumberIndex)
		if err != nil {
			serialNumber = "N/A"
		}
		log.Printf("Found S/N: %s", serialNumber)
	}

}
