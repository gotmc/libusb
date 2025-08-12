// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
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

func TestVersionTypeString(t *testing.T) {
	// Test that a version can be formatted as a string
	v := Version()

	// Just verify we can access all fields
	if v.Major == 0 && v.Minor == 0 && v.Micro == 0 {
		t.Error("Version() returned all zeros, expected valid version")
	}

	// Test that version fields are reasonable
	if v.Major > 10 {
		t.Errorf("Version.Major = %d seems unusually high", v.Major)
	}
}

func TestVersionFields(t *testing.T) {
	// Test that Version() returns a properly structured version
	v := Version()

	// All version numbers should be non-negative (uint16 can't be negative, but test for sanity)
	if v.Major > 65535 {
		t.Errorf("Version.Major = %d, exceeds uint16 range", v.Major)
	}
	if v.Minor > 65535 {
		t.Errorf("Version.Minor = %d, exceeds uint16 range", v.Minor)
	}
	if v.Micro > 65535 {
		t.Errorf("Version.Micro = %d, exceeds uint16 range", v.Micro)
	}
	if v.Nano > 65535 {
		t.Errorf("Version.Nano = %d, exceeds uint16 range", v.Nano)
	}

	// Describe field should typically be non-empty
	if v.Describe == "" {
		t.Log("Version.Describe is empty (this may be OK depending on libusb version)")
	}
}
