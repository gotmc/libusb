// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import (
	"fmt"
	"testing"
)

func TestBcdType(t *testing.T) {
	testCases := []struct {
		bcdValue   uint16
		bcdString  string
		bcdDecimal float64
	}{
		{0x0300, "0x0300 (3.00)", 3.0},
		{0x0200, "0x0200 (2.00)", 2.0},
		{0x0110, "0x0110 (1.10)", 1.1},
		{0x0100, "0x0100 (1.00)", 1.0},
	}
	t.Log("Given the need to test the bcd type")
	{
		for _, testCase := range testCases {
			b := bcd(testCase.bcdValue)
			t.Logf("\tWhen getting the string for bcd %#04x", testCase.bcdValue)
			computedString := fmt.Sprintf("%s", b)
			if computedString != testCase.bcdString {
				t.Errorf("\t%v Should have computed %s but got %s", failCheck, testCase.bcdString, computedString)
			} else {
				t.Logf("\t%v Should compute %s", passCheck, computedString)
			}
			// t.Logf("\tWhen getting the decimal value for bcd %#x", testCase.bcdValue)
			// computedDecimal := testCase.bcdValue.AsDecimal()
			// if computedDecimal != testCase.bcdDecimal {
			// t.Errorf("\t%v Should return decimal %2.2f but got %2.2f",
			// failCheck,
			// testCase.bcdDecimal,
			// computedDecimal,
			// )
			// } else {
			// t.Logf("\t%v Should return decimal %2.2f", passCheck, computedDecimal)
			// }
		}
	}
}
