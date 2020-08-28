// Copyright (c) 2015-2020 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import "testing"

func TestVersion(t *testing.T) {
	const major = 1
	const minor = 0
	const minMicro = 17
	releaseCandidate := ""
	version := Version()
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
	if version.Micro < minMicro {
		t.Errorf(
			"Micro version == %d, need at least %d",
			version.Micro, minMicro)
	}
	if version.ReleaseCandidate != releaseCandidate {
		t.Errorf(
			"ReleaseCandidate == %v, want %v",
			version.ReleaseCandidate, releaseCandidate)
	}
}
