package libusb

type bInterfaceClass byte

const reservedField = 0x00

/*
 * usbtmc.c by Agilent/Stefan Kopp sets this to 2048 with the comment:
 * Size of driver internal IO buffer. Must be multiple of 4 and at least as
 * large as wMaxPacketSize (which is usually 512 bytes).
 */
const ioBufferSize = 1024 * 1024 // Set to 1MB
const headerSize = 12

// The base class codes are part of the USB class codes "used to identify a
// device's functionality and to nominally load a device driver based on that
// functionality." [Source: http://www.usb.org/developers/defined_class]
// The USBTMC standard refers to these as the bInterfaceClass. The only base
// class code required for USBTMC is the Application Specific Base Class 0xFE.
const (
	applicationSpecificBaseClass bInterfaceClass = 0xfe
)

type bInterfaceSubClass byte

// The sub class is the second of the three bytes comprising the USB class
// code identifying a device's functionality.
const (
	usbtmcSubClass bInterfaceSubClass = 0x03
)

type bInterfaceProtocol byte

// The interface protocol is the third of three bytes comprising the USB class
// code identifying a device's functionality.
const (
	usbtmcProtocol bInterfaceProtocol = 0x00
	usb488Protocol bInterfaceProtocol = 0x01
)

type msgID uint8

// The following msgID values are found in Table 2 under the MACRO column of
// the USBTMC Specificiation 1.0, April 14, 2003. The end of line comment shows
// the MACRO names as given in the USBTMC specification.
// The trigger msgID comes from Table 1 -- USB488 defined msgID values of the
// USBTMC-USB488 Specification 1.0, April 14, 2003.
const (
	devDepMsgOut            msgID = 1   // DEV_DEP_MSG_OUT
	requestDevDepMsgIn      msgID = 2   // REQUEST_DEV_DEP_MSG_IN
	devDepMsgIn             msgID = 2   // DEV_DEP_MSG_IN
	vendorSpecificOut       msgID = 126 // VENDOR_SPECIFIC_OUT
	requestVendorSpecificIn msgID = 127 // REQUEST_VENDOR_SPECIFIC_IN
	vendorSpecificIn        msgID = 127 // VENDOR_SPECIFIC_IN
	trigger                 msgID = 128 // TRIGGER
)

type bRequest uint8

// The USBMTC bRequest constants come from Table 15 -- USBTMC bRequest values in
// the USBTMC Specificiation 1.0, April 14, 2003.
const (
	initiateAbortBulkOut    bRequest = 1  // INITIATE_ABORT_BULK_OUT
	checkAbortBulkOutStatus bRequest = 2  // CHECK_ABORT_BULK_OUT_STATUS
	initiateAbortBulkIn     bRequest = 3  // INITIATE_ABORT_BULK_IN
	checkAbortBulkInStatus  bRequest = 4  // CHECK_ABORT_BULK_IN_STATUS
	initiateClear           bRequest = 5  // INITIATE_CLEAR
	checkClearStatus        bRequest = 6  // CHECK_CLEAR_STATUS
	getCapabilities         bRequest = 7  // GET_CAPABILITIES
	indicatorPulse          bRequest = 64 // INDICATOR_PULSE
)

// The USB488 bRequest constants come from Table 9 -- USB488 defined bRequest
// values in the USBTMC-USB488 Specification 1.0, April 14, 2003
const (
	readStatusByte bRequest = 128 // READ_STATUS_BYTE
	renControl     bRequest = 160 // REN_CONTROL
	goToLocal      bRequest = 161 // GO_TO_LOCAL
	localLockout   bRequest = 162 // LOCAL_LOCKOUT
)

var requestDescription = map[bRequest]string{
	initiateAbortBulkOut:    "Aborts a Bulk-OUT transfer.",
	checkAbortBulkOutStatus: "Returns the status of the previously sent initiateAbortBulkOut request.",
	initiateAbortBulkIn:     "Aborts a Bulk-IN transfer.",
	checkAbortBulkInStatus:  "Returns the status of the previously sent initiateAbortBulkIn request.",
	initiateClear:           "Clears all previously sent pending and unprocessed Bulk-OUT USBTMC message content and clears all pending Bulk-IN transfers from the USBTMC interface.",
	checkClearStatus:        "Returns the status of the previously sent initiateClear request.",
	getCapabilities:         "Returns attributes and capabilities of the USBTMC interface.",
	indicatorPulse:          "A mechanism to turn on an activity indicator for identification purposes. The device indicates whether or not it supports this request in the getCapabilities response packet.",
	readStatusByte:          "Returns the IEEE 488 Status Byte.",
	renControl:              "Mechanism to enable or disable local controls on a device.",
	goToLocal:               "Mechanism to enable local controls on a device.",
	localLockout:            "Mechanism to disable local controls on a device.",
}

func (req bRequest) String() string {
	return requestDescription[req]
}

type status byte

// The status constant values come from Table 16 -- USBTMC_status values USBTMC
// Specificiation 1.0, April 14, 2003, and from Table 10 -- USB488 defined
// USBTMC_status values in the USBTMC-USB488 Specification 1.0, April 14, 2003
const (
	statusSuccess               status = 0x01 // STATUS_SUCCESS
	statusPending               status = 0x02 // STATUS_PENDING
	statusInterruptInBusy       status = 0x20 // STATUS_INTERRUPT_IN_BUSY
	statusFailed                status = 0x80 // STATUS_FAILED
	statusTransferNotInProgress status = 0x81 // STATUS_TRANSFER_NOT_IN_PROGRESS
	statusSplitNotInProgress    status = 0x82 // STATUS_SPLIT_NOT_IN_PROGRESS
	statusSplitInProgress       status = 0x83 // STATUS_SPLIT_IN_PROGRESS
)
