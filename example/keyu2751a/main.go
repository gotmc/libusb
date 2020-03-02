// Copyright (c) 2015-2020 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gotmc/libusb"
)

const reservedField = 0x00

const (
	devDepMsgOut    msgID = 1 // DEV_DEP_MSG_OUT
	reqDevDepMsgOut msgID = 2 // REQUEST_DEV_DEP_MSG_IN
)

type msgID uint8

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
	showInfo(ctx, "Agilent U2751A", 2391, 15896)

}

func showInfo(ctx *libusb.Context, name string, vendorID, productID uint16) {
	fmt.Printf("Let's open the %s using vendor ID %d and product ID %d\n", name, vendorID, productID)
	usbDevice, usbDeviceHandle, err := ctx.OpenDeviceWithVendorProduct(vendorID, productID)
	if err != nil {
		fmt.Printf("=> Failed opening the %s: %v\n", name, err)
		return
	}
	usbDeviceDescriptor, err := usbDevice.DeviceDescriptor()
	if err != nil {
		fmt.Printf("=> Failed getting the device descriptor for %s: %v\n", name, err)
		return
	}
	defer usbDeviceHandle.Close()
	serialnum, _ := usbDeviceHandle.StringDescriptorASCII(
		usbDeviceDescriptor.SerialNumberIndex,
	)
	manufacturer, _ := usbDeviceHandle.StringDescriptorASCII(
		usbDeviceDescriptor.ManufacturerIndex)
	product, _ := usbDeviceHandle.StringDescriptorASCII(
		usbDeviceDescriptor.ProductIndex)
	fmt.Printf("Manufacturer = %s\n", strings.TrimSpace(manufacturer))
	fmt.Printf("Product = %s\n", strings.TrimSpace(product))
	fmt.Printf("S/N = %s\n", strings.TrimSpace(serialnum))
	configDescriptor, err := usbDevice.ActiveConfigDescriptor()
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
	fmt.Printf(
		"   => USB-IF class %d, subclass %d, protocol %d.\n",
		firstDescriptor.InterfaceClass, firstDescriptor.InterfaceSubClass, firstDescriptor.InterfaceProtocol,
	)
	for i, endpoint := range firstDescriptor.EndpointDescriptors {
		fmt.Printf(
			"   => Endpoint index %d on Interface %d has the following properties:\n",
			i, firstDescriptor.InterfaceNumber)
		fmt.Printf("     => Address: %d (b%08b)\n", endpoint.EndpointAddress, endpoint.EndpointAddress)
		fmt.Printf("       => Endpoint #: %d\n", endpoint.Number())
		fmt.Printf("       => Direction: %s (%d)\n", endpoint.Direction(), endpoint.Direction())
		fmt.Printf("     => Attributes: %d (b%08b) \n", endpoint.Attributes, endpoint.Attributes)
		fmt.Printf("       => Transfer Type: %s (%d) \n", endpoint.TransferType(), endpoint.TransferType())
		fmt.Printf("     => Max packet size: %d\n", endpoint.MaxPacketSize)
	}

	// Initiate clear
	p1 := make([]byte, 1)
	_, err = usbDeviceHandle.ControlTransfer(0xC0, 0x0C, 0x0000, 0x047E, p1, 0x0001, 2000)
	if err != nil {
		log.Printf("Error sending first control transfer: %s", err)
	}
	p2 := make([]byte, 6)
	_, err = usbDeviceHandle.ControlTransfer(0xC0, 0x0C, 0x0000, 0x047D, p2, 0x0006, 2000)
	if err != nil {
		log.Printf("Error sending second control transfer: %s", err)
	}
	p3 := make([]byte, 5)
	_, err = usbDeviceHandle.ControlTransfer(0xC0, 0x0C, 0x0000, 0x0484, p3, 0x0005, 2000)
	if err != nil {
		log.Printf("Error sending third control transfer: %s", err)
	}
	p4 := make([]byte, 0x000C)
	_, err = usbDeviceHandle.ControlTransfer(0xC0, 0x0C, 0x0000, 0x0472, p4, 0x000C, 2000)
	if err != nil {
		log.Printf("Error sending fourth control transfer: %s", err)
	}
	p5 := make([]byte, 0x0001)
	_, err = usbDeviceHandle.ControlTransfer(0xC0, 0x0C, 0x0000, 0x047A, p5, 0x0001, 2000)
	if err != nil {
		log.Printf("Error sending fifth control transfer: %s", err)
	}
	p6 := []byte{0x00, 0x00, 0x01, 0x01, 0x00, 0x00, 0x08, 0x01}
	_, err = usbDeviceHandle.ControlTransfer(0x40, 0x0C, 0x0000, 0x0475, p6, len(p6), 2000)
	if err != nil {
		log.Printf("Error sending sixth control transfer: %s", err)
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
