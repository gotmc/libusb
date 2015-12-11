// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

type Config struct {
	*ConfigDescriptor
	Device *Device
}

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
