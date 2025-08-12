// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"bytes"
	"encoding/binary"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	libusb "github.com/gotmc/libusb/v2"
)

const reservedField = 0x00

const (
	devDepMsgOut msgID = 1 // DEV_DEP_MSG_OUT
)

type msgID uint8

func showVersion() {
	version := libusb.Version()
	fmt.Printf(
		"Using C libusb version %d.%d.%d (%d)\n",
		version.Major,
		version.Minor,
		version.Micro,
		version.Nano,
	)
}

func main() {
	// Define command line flags
	var vendorID uint
	var productID uint
	var listOnly bool
	
	flag.UintVar(&vendorID, "vid", 2391, "USB Vendor ID (decimal or hex with 0x prefix)")
	flag.UintVar(&productID, "pid", 1031, "USB Product ID (decimal or hex with 0x prefix)")
	flag.BoolVar(&listOnly, "list", false, "List all USB devices and exit")
	
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -vid 2391 -pid 1031\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -vid 0x0957 -pid 0x0407\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -list\n", os.Args[0])
	}
	
	flag.Parse()
	
	// Validate vendor and product IDs
	if vendorID > 0xFFFF {
		log.Fatalf("Invalid vendor ID: %d (must be 0-65535)", vendorID)
	}
	if productID > 0xFFFF {
		log.Fatalf("Invalid product ID: %d (must be 0-65535)", productID)
	}
	
	showVersion()
	ctx, err := libusb.NewContext()
	if err != nil {
		log.Fatal("Couldn't create USB context. Ending now.")
	}
	defer ctx.Close()
	
	if listOnly {
		start := time.Now()
		devices, _ := ctx.DeviceList()
		fmt.Printf("Found %v USB devices (%.4fs elapsed).\n",
			len(devices),
			time.Since(start).Seconds(),
		)
		for _, usbDevice := range devices {
			deviceAddress, _ := usbDevice.DeviceAddress()
			deviceSpeed, _ := usbDevice.Speed()
			busNumber, _ := usbDevice.BusNumber()
			usbDeviceDescriptor, _ := usbDevice.DeviceDescriptor()
			fmt.Printf("Device address %v is on bus number %v\n=> %v\n",
				deviceAddress,
				busNumber,
				deviceSpeed,
			)
			fmt.Printf("=> Vendor: 0x%04X (%d) \tProduct: 0x%04X (%d)\n=> Class: %v\n",
				usbDeviceDescriptor.VendorID,
				usbDeviceDescriptor.VendorID,
				usbDeviceDescriptor.ProductID,
				usbDeviceDescriptor.ProductID,
				usbDeviceDescriptor.DeviceClass,
			)
			fmt.Printf("=> USB: %v\tMax Packet 0: %v\tSN Index: %v\n",
				usbDeviceDescriptor.USBSpecification,
				usbDeviceDescriptor.MaxPacketSize0,
				usbDeviceDescriptor.SerialNumberIndex,
			)
		}
		return
	}
	
	// Connect to the specified device
	showInfo(ctx, uint16(vendorID), uint16(productID))
}

func showInfo(ctx *libusb.Context, vendorID, productID uint16) {
	fmt.Printf("Opening device with VID:0x%04X PID:0x%04X\n", vendorID, productID)
	usbDevice, usbDeviceHandle, err := ctx.OpenDeviceWithVendorProduct(vendorID, productID)
	if err != nil {
		fmt.Printf(
			"=> Failed opening device (VID:0x%04x PID:0x%04x): %v\n",
			vendorID,
			productID,
			err,
		)
		return
	}
	usbDeviceDescriptor, err := usbDevice.DeviceDescriptor()
	if err != nil {
		fmt.Printf("=> Failed getting device descriptor: %v\n", err)
		usbDeviceHandle.Close()
		return
	}
	defer usbDeviceHandle.Close()

	// Get string descriptors with proper error handling
	serialnum, err := usbDeviceHandle.StringDescriptorASCII(usbDeviceDescriptor.SerialNumberIndex)
	if err != nil || serialnum == "" {
		serialnum = fmt.Sprintf("Device_%04X_%04X", vendorID, productID)
	}

	manufacturer, err := usbDeviceHandle.StringDescriptorASCII(
		usbDeviceDescriptor.ManufacturerIndex,
	)
	if err != nil {
		manufacturer = "<unavailable>"
	}

	product, err := usbDeviceHandle.StringDescriptorASCII(usbDeviceDescriptor.ProductIndex)
	if err != nil {
		product = "<unavailable>"
	}
	
	// Use serial number as the device identifier from now on
	fmt.Printf("=> Connected to device S/N: %s\n", serialnum)
	fmt.Printf("=> Manufacturer: %s, Product: %s\n", manufacturer, product)
	fmt.Printf("=> VID: 0x%04X, PID: 0x%04X\n", vendorID, productID)
	configDescriptor, err := usbDevice.ActiveConfigDescriptor()
	if err != nil {
		fmt.Printf("=> Failed getting the active config for S/N %s: %v\n", serialnum, err)
		return
	}
	fmt.Printf("=> Max Power = %d mA\n", configDescriptor.MaxPowerMilliAmperes)

	var singularPlural string
	if configDescriptor.NumInterfaces == 1 {
		singularPlural = "interface"
	} else {
		singularPlural = "interfaces"
	}
	fmt.Printf("=> Found %d %s\n", configDescriptor.NumInterfaces, singularPlural)

	// Check if we have interfaces before accessing them
	if len(configDescriptor.SupportedInterfaces) == 0 {
		fmt.Printf("=> No supported interfaces found for S/N %s\n", serialnum)
		return
	}

	firstInterface := configDescriptor.SupportedInterfaces[0]
	fmt.Printf("=> The first interface has %d alternate settings.\n", firstInterface.NumAltSettings)

	// Check if we have interface descriptors before accessing them
	if len(firstInterface.InterfaceDescriptors) == 0 {
		fmt.Printf("=> No interface descriptors found for S/N %s\n", serialnum)
		return
	}

	firstDescriptor := firstInterface.InterfaceDescriptors[0]
	fmt.Printf("=> The first interface descriptor has a length of %d.\n", firstDescriptor.Length)
	fmt.Printf(
		"=> The first interface descriptor is interface number %d.\n",
		firstDescriptor.InterfaceNumber,
	)
	fmt.Printf(
		"=> The first interface descriptor has %d endpoint(s).\n",
		firstDescriptor.NumEndpoints,
	)
	fmt.Printf(
		"   => USB-IF class %d, subclass %d, protocol %d.\n",
		firstDescriptor.InterfaceClass,
		firstDescriptor.InterfaceSubClass,
		firstDescriptor.InterfaceProtocol,
	)

	// Check if we have endpoint descriptors before accessing them
	if len(firstDescriptor.EndpointDescriptors) == 0 {
		fmt.Printf("=> No endpoint descriptors found for S/N %s\n", serialnum)
		return
	}

	for i, endpoint := range firstDescriptor.EndpointDescriptors {
		fmt.Printf(
			"   => Endpoint index %d on Interface %d has the following properties:\n",
			i, firstDescriptor.InterfaceNumber)
		fmt.Printf(
			"     => Address: %d (b%08b)\n",
			endpoint.EndpointAddress,
			endpoint.EndpointAddress,
		)
		fmt.Printf("       => Endpoint #: %d\n", endpoint.Number())
		fmt.Printf("       => Direction: %s (%d)\n", endpoint.Direction(), endpoint.Direction())
		fmt.Printf("     => Attributes: %d (b%08b) \n", endpoint.Attributes, endpoint.Attributes)
		fmt.Printf(
			"       => Transfer Type: %s (%d) \n",
			endpoint.TransferType(),
			endpoint.TransferType(),
		)
		fmt.Printf("     => Max packet size: %d\n", endpoint.MaxPacketSize)
	}

	err = usbDeviceHandle.ClaimInterface(0)
	if err != nil {
		log.Printf("Error claiming interface %s", err)
	}

	// Get Capabilities
	p := make([]byte, 64)
	idx := uint16(0x0000)
	n, err := usbDeviceHandle.ControlTransfer(0xA1, 7, 0x0000, idx, p, 0x18, 2000)
	if err != nil {
		log.Printf("Error sending control transfer: %s", err)
	}
	log.Printf("Sent %d bytes on control transfer", n)
	log.Printf("capabilities = %q", p)
	log.Printf("capabilities = %v", p)
	log.Printf("cap[14] := %b (%d)", p[14], p[14])
	log.Printf("cap[15] := %b (%d)", p[15], p[15])

	// Send USBTMC message to device
	if len(firstDescriptor.EndpointDescriptors) == 0 {
		fmt.Printf("=> No endpoints available for bulk transfer on S/N %s\n", serialnum)
		return
	}

	bulkOutput := firstDescriptor.EndpointDescriptors[0]
	address := bulkOutput.EndpointAddress
	fmt.Printf("Set frequency/amplitude on endpoint address %d\n", address)
	data := createGotmcMessage("apply:sinusoid 2340, 0.1, 0.0")
	transferred, err := usbDeviceHandle.BulkTransfer(address, data, len(data), 5000)
	if err != nil {
		fmt.Printf("=> Error on bulk transfer to S/N %s: %s\n", serialnum, err)
	} else {
		fmt.Printf("Sent %d bytes to S/N %s\n", transferred, serialnum)
	}
	err = usbDeviceHandle.ReleaseInterface(0)
	if err != nil {
		log.Printf("Error releasing interface %s", err)
	}
}

func createDevDepMsgOutBulkOutHeader(
	transferSize uint32, eom bool, bTag byte,
) [12]byte {
	// Offset 0-3: See Table 1.
	prefix := encodeBulkHeaderPrefix(devDepMsgOut, bTag)
	// Offset 4-7: TransferSize
	// Per USBTMC Table 3, the TransferSize is the "total number of USBTMC
	// message data bytes to be sent in this USB transfer. This does not include
	// the number of bytes in this Bulk-OUT Header or alignment bytes. Sent least
	// significant byte first, most significant byte last. TransferSize must be >
	// 0x00000000."
	packedTransferSize := make([]byte, 4)
	binary.LittleEndian.PutUint32(packedTransferSize, transferSize)
	// Offset 8: bmTransferAttributes
	// Per USBTMC Table 3, D0 of bmTransferAttributes:
	//   1 - The last USBTMC message data byte in the transfer is the last byte
	//       of the USBTMC message.
	//   0 - The last USBTMC message data byte in the transfer is not the last
	//       byte of the USBTMC message.
	// All other bits of bmTransferAttributes must be 0.
	bmTransferAttributes := byte(0x00)
	if eom {
		bmTransferAttributes = byte(0x01)
	}
	// Offset 9-11: reservedField. Must be 0x000000.
	return [12]byte{
		prefix[0],
		prefix[1],
		prefix[2],
		prefix[3],
		packedTransferSize[0],
		packedTransferSize[1],
		packedTransferSize[2],
		packedTransferSize[3],
		bmTransferAttributes,
		reservedField,
		reservedField,
		reservedField,
	}
}

// Create the first four bytes of the USBTMC meassage Bulk-OUT Header as shown
// in USBTMC Table 1. The msgID value must match USBTMC Table 2.
func encodeBulkHeaderPrefix(msgID msgID, bTag byte) [4]byte {
	return [4]byte{
		byte(msgID),
		bTag,
		invertbTag(bTag),
		reservedField,
	}
}

func invertbTag(bTag byte) byte {
	return bTag ^ 0xff
}

func createGotmcMessage(input string) []byte {
	message := []byte(input + "\n")
	header := createDevDepMsgOutBulkOutHeader(uint32(len(message)), true, 1)
	data := append(header[:], message...)
	if moduloFour := len(data) % 4; moduloFour > 0 {
		numAlignment := 4 - moduloFour
		alignment := bytes.Repeat([]byte{0x00}, numAlignment)
		data = append(data, alignment...)
	}
	return data
}
