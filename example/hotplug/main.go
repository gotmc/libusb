package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gotmc/libusb/v2"
)

var ctx *libusb.Context

func main() {
	var err error
	ctx, err = libusb.NewContext()
	if err != nil {
		log.Fatal("Couldn't create USB context. Ending now.")
	}
	defer ctx.Close()

	// By vID & pID
	const (
		vendorID  = 0x0930
		productID = 0x6545
	)
	log.Println("Connect or disconnect any USB device...")
	ctx.HotplugRegisterCallbackEvent(vendorID, productID, libusb.HotplugArrived|libusb.HotplugLeft, cb)
	time.Sleep(time.Second * 10)
	ctx.HotplugDeregisterCallback(vendorID, productID)

	// All devices
	log.Println("Connect or disconnect any USB device...")
	ctx.HotplugRegisterCallbackEvent(0, 0, libusb.HotplugArrived|libusb.HotplugLeft, cb)
	defer ctx.HotplugDeregisterAllCallbacks()
	time.Sleep(time.Second * 10)
}

func cb(vID, pID uint16, eventType libusb.HotPlugEventType) {
	fmt.Printf("VendorID: %04x, ProductID: %04x, eventType: %d\r\n", vID, pID, eventType)
}
