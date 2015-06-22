// Copyright (c) 2015 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

import "testing"

const (
	failCheck = `✗` // UTF-8 u2717
	passCheck = `✓` // UTF-8 u2713
)

func TestInitContext(t *testing.T) {
	t.Log("Given the need to test creating a new context.")
	{
		t.Log("\tWhen initializing a new context")
		_, err := Init()
		if err != nil {
			t.Error("\t", failCheck, "Should not have received error:", err)
		}
		t.Log("\t", passCheck, "Should not receive an error.")
	}
}

func TestExitContext(t *testing.T) {
	t.Log("Given the need to test exiting a context.")
	{
		t.Log("\tWhen exiting a context")
		context, err := Init()
		err = context.Exit()
		if err != nil {
			t.Error("\t", failCheck, "Should not have received an error:", err)
		}
		t.Log("\t", passCheck, "Should not receive an error.")
		if context.context != nil {
			t.Fatal("\t", failCheck, "The context field in the context struct should be nil.")
		}
		t.Log("\t", passCheck, "The context field in the context struct should be nil.")
	}
}
