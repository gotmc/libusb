// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
// int libusbHotplugCallback (libusb_context *ctx, libusb_device *device, libusb_hotplug_event event, void *user_data);
// typedef struct libusb_device_descriptor libusb_device_descriptor_struct;
// static int libusb_hotplug_register_callback_wrapper (
//	libusb_context *ctx,
//	int events, int flags,
//	int vendor_id, int product_id, int dev_class,
//	libusb_hotplug_callback_fn cb_fn, void *user_data,
//	libusb_hotplug_callback_handle *callback_handle)
//	{
// 		return libusb_hotplug_register_callback(ctx, events, flags, vendor_id, product_id, dev_class, cb_fn, user_data, callback_handle);
// }
import "C"
import (
	"fmt"
	"log"
	"sync"
	"unsafe"
)

// HotPlugEventType ...
type HotPlugEventType uint8

// HotPlugCbFunc callback
type HotPlugCbFunc func(vID, pID uint16, eventType HotPlugEventType)

// HotPlug Event Types
const (
	HotplugUndefined HotPlugEventType = iota
	HotplugArrived
	HotplugLeft
)

// HotPlugEvent callback message
type HotPlugEvent struct {
	VendorID  uint16
	ProductID uint16
	Event     HotPlugEventType
}

type hotplugCallback struct {
	handler *C.libusb_hotplug_callback_handle
	fn      HotPlugCbFunc
}

// HotplugCallbackStorage ...
type HotplugCallbackStorage struct {
	callbackMap map[uint32]hotplugCallback
	done        chan struct{}
	mu          sync.RWMutex // Protects the callbackMap from concurrent access
}

var hotplugCallbackStorage HotplugCallbackStorage

func (ctx *Context) newHotPlugHandler() {
	hotplugCallbackStorage.callbackMap = make(map[uint32]hotplugCallback)
	hotplugCallbackStorage.done = make(chan struct{})

	go hotplugCallbackStorage.handleEvents(ctx.libusbContext)
}

func (s *HotplugCallbackStorage) isEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.callbackMap == nil
}

// HotplugRegisterCallbackEvent ...
func (ctx *Context) HotplugRegisterCallbackEvent(vendorID, productID uint16, eventType HotPlugEventType, cb HotPlugCbFunc) error {
	if hotplugCallbackStorage.isEmpty() {
		ctx.newHotPlugHandler()
	}

	var event C.int
	switch eventType {
	case HotplugArrived:
		event = C.LIBUSB_HOTPLUG_EVENT_DEVICE_ARRIVED
	case HotplugLeft:
		event = C.LIBUSB_HOTPLUG_EVENT_DEVICE_LEFT
	default:
		event = C.LIBUSB_HOTPLUG_EVENT_DEVICE_ARRIVED | C.LIBUSB_HOTPLUG_EVENT_DEVICE_LEFT
	}

	var vID C.int = C.LIBUSB_HOTPLUG_MATCH_ANY
	var pID C.int = C.LIBUSB_HOTPLUG_MATCH_ANY

	if vendorID != 0 {
		vID = C.int(vendorID)
	}
	if productID != 0 {
		pID = C.int(productID)
	}

	var cbHandle C.libusb_hotplug_callback_handle

	rc := C.libusb_hotplug_register_callback_wrapper(
		ctx.libusbContext,
		event,
		C.LIBUSB_HOTPLUG_NO_FLAGS,
		vID,
		pID,
		C.LIBUSB_HOTPLUG_MATCH_ANY,
		C.libusb_hotplug_callback_fn(unsafe.Pointer(C.libusbHotplugCallback)),
		nil,
		&cbHandle,
	)
	if rc != C.LIBUSB_SUCCESS {
		return fmt.Errorf("libusb_hotplug_register_callback error: %s", ErrorCode(rc))
	}

	hotplugCallbackStorage.mu.Lock()
	hotplugCallbackStorage.callbackMap[vidPidToUint32(vendorID, productID)] = hotplugCallback{
		handler: &cbHandle,
		fn:      cb,
	}
	hotplugCallbackStorage.mu.Unlock()

	return nil
}

// HotplugDeregisterCallback ...
func (ctx *Context) HotplugDeregisterCallback(vendorID, productID uint16) error {
	if hotplugCallbackStorage.isEmpty() {
		return nil
	}

	key := vidPidToUint32(vendorID, productID)

	hotplugCallbackStorage.mu.RLock()
	cb, ok := hotplugCallbackStorage.callbackMap[key]
	hotplugCallbackStorage.mu.RUnlock()

	if !ok {
		return nil
	}

	C.libusb_hotplug_deregister_callback(ctx.libusbContext, *cb.handler)

	hotplugCallbackStorage.mu.Lock()
	delete(hotplugCallbackStorage.callbackMap, key)
	mapEmpty := len(hotplugCallbackStorage.callbackMap) == 0
	hotplugCallbackStorage.mu.Unlock()

	if mapEmpty {
		ctx.hotplugHandleEventsCompleteAll()
	}
	return nil
}

// HotplugDeregisterAllCallbacks ...
func (ctx *Context) HotplugDeregisterAllCallbacks() error {
	hotplugCallbackStorage.mu.RLock()
	mapExists := hotplugCallbackStorage.callbackMap != nil

	if mapExists {
		// Make a copy of the handlers to avoid holding the lock during C function calls
		handlers := make([]*C.libusb_hotplug_callback_handle, 0, len(hotplugCallbackStorage.callbackMap))
		for _, cb := range hotplugCallbackStorage.callbackMap {
			handlers = append(handlers, cb.handler)
		}
		hotplugCallbackStorage.mu.RUnlock()

		// Deregister callbacks without holding the lock
		for _, handler := range handlers {
			C.libusb_hotplug_deregister_callback(ctx.libusbContext, *handler)
		}
	} else {
		hotplugCallbackStorage.mu.RUnlock()
	}

	ctx.hotplugHandleEventsCompleteAll()

	return nil
}

func (ctx *Context) hotplugHandleEventsCompleteAll() {
	if hotplugCallbackStorage.isEmpty() {
		return
	}

	// Signal the event handler to stop
	hotplugCallbackStorage.done <- struct{}{}

	// Clear the callbackMap and close the channel
	hotplugCallbackStorage.mu.Lock()
	hotplugCallbackStorage.callbackMap = nil
	hotplugCallbackStorage.mu.Unlock()

	close(hotplugCallbackStorage.done)
}

func (storage *HotplugCallbackStorage) handleEvents(libCtx *C.libusb_context) {
	for {
		select {
		case <-storage.done:
			return
		default:
		}
		if errno := C.libusb_handle_events_completed(libCtx, nil); errno < 0 {
			log.Printf("handle_events error: %s", ErrorCode(errno))
		}
	}
}

//export libusbHotplugCallback
func libusbHotplugCallback(ctx *C.libusb_context, dev *C.libusb_device, event C.libusb_hotplug_event, p unsafe.Pointer) C.int {
	var desc C.libusb_device_descriptor_struct
	rc := C.libusb_get_device_descriptor(dev, &desc)
	if rc != C.LIBUSB_SUCCESS {
		return rc
	}

	var vendorID = uint16(desc.idVendor)
	var productID = uint16(desc.idProduct)

	var e HotPlugEventType
	switch event {
	case C.LIBUSB_HOTPLUG_EVENT_DEVICE_ARRIVED:
		e = HotplugArrived
	case C.LIBUSB_HOTPLUG_EVENT_DEVICE_LEFT:
		e = HotplugLeft
	default:
		e = HotplugUndefined
	}

	// Read callback map with a read lock
	hotplugCallbackStorage.mu.RLock()
	// Get device-specific callback
	cb, ok := hotplugCallbackStorage.callbackMap[vidPidToUint32(vendorID, productID)]
	var deviceCallback HotPlugCbFunc
	if ok {
		deviceCallback = cb.fn
	}

	// Get the callback for all devices
	cb, ok = hotplugCallbackStorage.callbackMap[0]
	var allCallback HotPlugCbFunc
	if ok {
		allCallback = cb.fn
	}
	hotplugCallbackStorage.mu.RUnlock()

	// Call callbacks outside the lock
	if deviceCallback != nil {
		deviceCallback(vendorID, productID, e)
	}

	if allCallback != nil {
		allCallback(vendorID, productID, e)
	}

	return C.LIBUSB_SUCCESS
}

func vidPidToUint32(vID, pID uint16) uint32 {
	return (uint32(vID) << 16) | (uint32(pID))
}
