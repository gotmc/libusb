// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"

// The Version struct represents the libusb version.
type version struct {
	Major            uint16
	Minor            uint16
	Micro            uint16
	Nano             uint16
	ReleaseCandidate string
	Describe         string
}

// GetVersion gets the libusb version and returns a Version struct.
func GetVersion() version {
	var cVersion C.struct_libusb_version
	cVersion = *C.libusb_get_version()
	version := version{
		Major:            uint16(cVersion.major),
		Minor:            uint16(cVersion.minor),
		Micro:            uint16(cVersion.micro),
		Nano:             uint16(cVersion.nano),
		ReleaseCandidate: C.GoString(cVersion.rc),
		Describe:         C.GoString(cVersion.describe),
	}
	return version
}
