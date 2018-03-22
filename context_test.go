// Copyright (c) 2015-2017 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import "testing"

func TestNewContext(t *testing.T) {
	if _, err := NewContext(); err != nil {
		t.Errorf(
			"Error initializing new libusb context:\n\tgot %v want %v",
			err,
			nil,
		)
	}
}

func TestCloseContext(t *testing.T) {
	context, _ := NewContext()
	if err := context.Close(); err != nil {
		t.Errorf(
			"Error exiting context:\n\tgot %v want %v",
			err,
			nil,
		)
	}
	if context.libusbContext != nil {
		t.Errorf(
			"Context field still exists after exiting:\n\tgot %v want %v",
			context.libusbContext,
			nil,
		)
	}
}

func TestSetDebugLevel(t *testing.T) {
	testCases := []struct {
		lev LogLevel
	}{
		{LogLevelNone},
		{LogLevelError},
		{LogLevelWarning},
		{LogLevelInfo},
		{LogLevelDebug},
	}
	for _, tc := range testCases {
		context, _ := NewContext()
		context.SetDebug(tc.lev)
		if got := context.LogLevel; got != tc.lev {
			t.Errorf("got %v; want %v", got, tc.lev)
		}
	}
}

func TestLogLevelStringMethod(t *testing.T) {
	testCases := []struct {
		logLevel LogLevel
		want     string
	}{
		{LogLevelNone, "No messages ever printed by the library (default)"},
		{LogLevelError, "Error messages are printed to stderr"},
		{LogLevelWarning, "Warning and error messages are printed to stderr"},
		{LogLevelInfo, "Informational messages are printed to stdout, warning and error messages are printed to stderr"},
		{LogLevelDebug, "Debug and informational messages are printed to stdout, warnings and errors to stderr"},
	}
	for _, tc := range testCases {
		if got := tc.logLevel.String(); got != tc.want {
			t.Errorf("got %s; want %s", got, tc.want)
		}
	}
}

func TestGetDeviceList(t *testing.T) {
	context, _ := NewContext()
	defer context.Close()
	devices, err := context.GetDeviceList()
	if err != nil {
		t.Errorf(
			"Error on GetDeviceList:\n\tgot %v; want %v",
			err,
			nil,
		)
	}
	if got := len(devices); got < 1 {
		t.Errorf(
			"got %v device; want at least one",
			got,
		)
	}
}
