// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package main

import (
	"fmt"

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
	ctx, _ := libusb.Init()
	ctx.Exit()

}
