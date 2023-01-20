// Copyright (c) 2015-2022 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"

// The VersionType struct represents the libusb version.
type VersionType struct {
	Major            uint16
	Minor            uint16
	Micro            uint16
	Nano             uint16
	ReleaseCandidate string
	Describe         string
}

// Version gets the libusb version and returns a Version struct.
func Version() VersionType {
	cVersion := *C.libusb_get_version()
	version := VersionType{
		Major:            uint16(cVersion.major),
		Minor:            uint16(cVersion.minor),
		Micro:            uint16(cVersion.micro),
		Nano:             uint16(cVersion.nano),
		ReleaseCandidate: C.GoString(cVersion.rc),
		Describe:         C.GoString(cVersion.describe),
	}
	return version
}
