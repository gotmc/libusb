// Copyright (c) 2015-2020 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// Config models the USB configuration.
type Config struct {
	*ConfigDescriptor
	Device *Device
}

// ConfigDescriptor models the descriptor for the USB configuration
type ConfigDescriptor struct {
	Length               int
	DescriptorType       descriptorType
	TotalLength          uint16
	NumInterfaces        int
	ConfigurationValue   uint8
	ConfigurationIndex   uint8
	Attributes           uint8
	MaxPowerMilliAmperes uint
	SupportedInterfaces
}
