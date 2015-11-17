// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import (
	"fmt"
	"testing"

	c "github.com/smartystreets/goconvey/convey"
)

const (
	failCheck = `✗` // UTF-8 u2717
	passCheck = `✓` // UTF-8 u2713
)

func TestInitContext(t *testing.T) {
	c.Convey("Given the need to test creating a new libusb context.", t, func() {
		c.Convey("When initializing a new context", func() {
			c.Convey("Then an error should not be received", func() {
				_, err := Init()
				c.So(err, c.ShouldBeNil)
			})
		})
	})
}

func TestExitContext(t *testing.T) {
	c.Convey("Given the need to test exiting a context.", t, func() {
		context, err := Init()
		err = context.Exit()
		c.Convey("When exiting a context", func() {
			c.Convey("Then an error should not be received.", func() {
				c.So(err, c.ShouldBeNil)
			})
		})
		c.Convey("After exiting a context", func() {
			c.Convey("The context field should be nil.", func() {
				c.So(context.libusbContext, c.ShouldBeNil)
			})
		})
	})
}

func TestSetDebugLevel(t *testing.T) {
	c.Convey("Given the need to set the debug/log level", t, func() {
		c.Convey("When the log level is set to error", func() {
			c.Convey("The LogLevel should be set to LogLevelError", func() {
				context, _ := Init()
				context.SetDebug(LogLevelError)
				c.So(context.LogLevel, c.ShouldEqual, LogLevelError)
			})
		})
	})
}

func TestLogLevelStringMethod(t *testing.T) {
	testCases := []struct {
		logLevel logLevel
		expected string
	}{
		{LogLevelNone, "No messages ever printed by the library (default)"},
		{LogLevelError, "Error messages are printed to stderr"},
		{LogLevelWarning, "Warning and error messages are printed to stderr"},
		{LogLevelInfo, "Informational messages are printed to stdout, warning and error messages are printed to stderr"},
		{LogLevelDebug, "Debug and informational messages are printed to stdout, warnings and errors to stderr"},
	}
	c.Convey("Given the need to test the logleve.String() method", t, func() {
		for _, testCase := range testCases {
			conveyance := fmt.Sprintf("When getting logLevel %d's string", testCase.logLevel)
			c.Convey(conveyance, func() {
				computed := testCase.logLevel.String()
				conveyance := fmt.Sprintf("Then the logLevel string should be %s", testCase.expected)
				c.Convey(conveyance, func() {
					c.So(computed, c.ShouldEqual, testCase.expected)
				})
			})
		}
	})
}

func TestGetDeviceList(t *testing.T) {
	c.Convey("Given the desire to get a device list", t, func() {
		c.Convey("When the GetDeviceList() method is called", func() {
			c.Convey("Then an array of pointers to a Device should be returned", func() {
				context, _ := Init()
				defer context.Exit()
				devices, err := context.GetDeviceList()
				c.So(err, c.ShouldBeNil)
				c.So(len(devices), c.ShouldBeGreaterThan, 0)
			})
		})
	})
}
