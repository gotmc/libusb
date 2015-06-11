// Package libusb provides a Go bindings for the  libusb C library.
package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
import "C"
import "strings"

type Version struct {
	Major            uint16
	Minor            uint16
	Micro            uint16
	Nano             uint16
	ReleaseCandidate string
	Describe         string
}

func GetVersion() Version {
	var cVersion C.struct_libusb_version
	cVersion = *C.libusb_get_version()
	version := Version{
		Major:            uint16(cVersion.major),
		Minor:            uint16(cVersion.minor),
		Micro:            uint16(cVersion.micro),
		Nano:             uint16(cVersion.nano),
		ReleaseCandidate: string(*cVersion.rc),
		Describe:         string(*cVersion.describe),
	}
	version.ReleaseCandidate = strings.TrimRight(version.ReleaseCandidate, "\x00")
	return version
}
