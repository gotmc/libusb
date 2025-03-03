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
	ctx, err := libusb.NewContext()
	if err != nil {
		log.Fatalf("Error creating USB context: %v", err)
	}
	defer ctx.Close()

	// For demo purposes only - replace with your actual values
	var vendorID uint16 = 0x1234
	var productID uint16 = 0x5678
	
	// Open device
	device, handle, err := ctx.OpenDeviceWithVendorProduct(vendorID, productID)
	if err != nil {
		log.Printf("Could not open device %04x:%04x, using example code only", vendorID, productID)
		// Continue with example code anyway
	} else {
		defer handle.Close()
		log.Printf("Opened device %04x:%04x", vendorID, productID)
	}

	// Example 1: Using the original approach with raw byte
	// This requires bit manipulation and can be error-prone
	fmt.Println("\nExample 1: Original approach with raw bytes")
	fmt.Println("----------------------------------------")
	
	// Construct requestType using bit manipulation: Class request (0x01<<5) to interface (0x01)
	rawRequestType := byte((0x01 << 5) + 0x01) 
	fmt.Printf("Raw requestType: 0x%02x\n", rawRequestType)
	
	// Example code - would actually be executed if handle != nil
	fmt.Printf("handle.ControlTransfer(0x%02x, 0x09, 0x0300, 0, data, len, 5000)\n", rawRequestType)

	// Example 2: Using the new type-safe approach with constants
	fmt.Println("\nExample 2: Type-safe approach with exported types")
	fmt.Println("----------------------------------------")
	
	// Use library constants for better readability and type safety
	direction := libusb.HostToDevice
	reqType := libusb.Class
	recipient := libusb.InterfaceRecipient
	
	fmt.Printf("Direction: %v\n", direction)
	fmt.Printf("RequestType: %v\n", reqType)
	fmt.Printf("Recipient: %v\n", recipient)
	
	// Create combined requestType using BitmapRequestType
	bmRequestType := libusb.BitmapRequestType(direction, reqType, recipient)
	fmt.Printf("Combined bmRequestType: 0x%02x\n", bmRequestType)
	
	// Example code
	fmt.Printf("handle.ControlTransferWithTypes(%v, %v, %v, 0x09, 0x0300, 0, data, len, 5000)\n",
		direction, reqType, recipient)
	
	// Example 3: Using the new helper methods
	fmt.Println("\nExample 3: Using convenience helper methods")
	fmt.Println("----------------------------------------")
	
	fmt.Printf("handle.ControlOut(%v, %v, 0x09, 0x0300, 0, data, 5000)\n",
		reqType, recipient)
}