// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import "testing"

func TestGetVersion(t *testing.T) {
	const major = 1
	const minor = 0
	const micro = 19
	const nano = 10903
	releaseCandidate := ""
	describe := "http://libusb.info"
	version := GetVersion()
	if version.Major != major {
		t.Errorf(
			"Major version == %d, want %d",
			version.Major, major)
	}
	if version.Minor != minor {
		t.Errorf(
			"Minor version == %d, want %d",
			version.Minor, minor)
	}
	if version.Micro != micro {
		t.Errorf(
			"Micro version == %d, want %d",
			version.Micro, micro)
	}
	if version.Nano != nano {
		t.Errorf(
			"Nano version == %d, want %d",
			version.Nano, nano)
	}
	if version.ReleaseCandidate != releaseCandidate {
		t.Errorf(
			"ReleaseCandidate == %v, want %v",
			version.ReleaseCandidate, releaseCandidate)
	}
	if version.Describe != describe {
		t.Errorf(
			"Describe == %v, want %v",
			version.Describe, describe)
	}
}
