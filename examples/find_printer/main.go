// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"fmt"
	"log"

	"github.com/gotmc/libusb/v2"
)

func main() {
	// Create a new libusb context
	ctx, err := libusb.NewContext()
	if err != nil {
		log.Fatalf("Error creating context: %s", err)
	}
	defer ctx.Close()

	// Get list of devices
	devices, err := ctx.DeviceList()
	if err != nil {
		log.Fatalf("Error getting device list: %s", err)
	}
	// Clean up device references when done
	defer func() {
		for _, device := range devices {
			device.Close()
		}
	}()

	fmt.Printf("Found %d USB devices\n", len(devices))

	// Loop through all devices
	for i, device := range devices {
		desc, err := device.DeviceDescriptor()
		if err != nil {
			log.Printf("Error getting device descriptor for device %d: %s", i, err)
			continue
		}

		// Print the device class information
		fmt.Printf("\nDevice %d: VID=0x%04x, PID=0x%04x\n", i, desc.VendorID, desc.ProductID)
		fmt.Printf("  Device Class: %s (0x%02x)\n", desc.DeviceClass, byte(desc.DeviceClass))

		// If this is a per-interface class device (class 0), inspect the interfaces
		if desc.DeviceClass == 0 {
			fmt.Println("  This device uses per-interface class codes. Checking interfaces...")

			// Find printer interfaces (class 7)
			printerIfaces, err := device.FindInterfacesByClass(libusb.InterfaceClassPrinter)
			if err != nil {
				log.Printf("  Error getting interface info: %s", err)
				continue
			}

			if len(printerIfaces) > 0 {
				fmt.Printf("  Found %d printer interface(s)!\n", len(printerIfaces))

				// Print details about each printer interface
				for j, iface := range printerIfaces {
					fmt.Printf("  Printer Interface %d:\n", j)
					fmt.Printf("    Interface Number: %d\n", iface.InterfaceNumber)
					fmt.Printf("    Interface Class: 0x%02x (Printer)\n", iface.InterfaceClass)
					fmt.Printf("    Interface SubClass: 0x%02x\n", iface.InterfaceSubClass)
					fmt.Printf("    Interface Protocol: 0x%02x\n", iface.InterfaceProtocol)
					fmt.Printf("    Endpoints: %d\n", iface.NumEndpoints)
				}
			}
		}
	}
}

