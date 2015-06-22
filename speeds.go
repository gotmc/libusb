// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"

type speed int

const (
	speedUnknown speed = C.LIBUSB_SPEED_UNKNOWN
	speedLow     speed = C.LIBUSB_SPEED_LOW
	speedFull    speed = C.LIBUSB_SPEED_FULL
	speedHigh    speed = C.LIBUSB_SPEED_HIGH
	speedSuper   speed = C.LIBUSB_SPEED_SUPER
)

var speedCodes = map[speed]string{
	speedUnknown: "The OS doesn't report or know the device speed.",
	speedLow:     "The device is operating at low speed (1.5MBit/s)",
	speedFull:    "The device is operating at full speed (12MBit/s)",
	speedHigh:    "The device is operating at high speed (480MBit/s)",
	speedSuper:   "The device is operating at super speed (5000MBit/s)",
}

func (speed speed) String() string {
	return speedCodes[speed]
}

type supportedSpeed int

const (
	lowSpeedOperation   supportedSpeed = C.LIBUSB_LOW_SPEED_OPERATION
	fullSpeedOperation  supportedSpeed = C.LIBUSB_FULL_SPEED_OPERATION
	highSpeedOperation  supportedSpeed = C.LIBUSB_HIGH_SPEED_OPERATION
	superSpeedOperation supportedSpeed = C.LIBUSB_SUPER_SPEED_OPERATION
)

var supportedSpeeds = map[supportedSpeed]string{
	lowSpeedOperation:   "Low speed operation supported (1.5MBit/s).",
	fullSpeedOperation:  "Full speed operation supported (12MBit/s).",
	highSpeedOperation:  "High speed operation supported (480MBit/s).",
	superSpeedOperation: "Superspeed operation supported (5000MBit/s).",
}

func (speed supportedSpeed) String() string {
	return supportedSpeeds[speed]
}
