// Copyright (c) 2015-2025 The libusb developers. All rights reserved.
// Project site: https://github.com/gotmc/libusb
// Use of this source code is governed by a MIT-style license that
// can be found in the LICENSE.txt file for the project.

package libusb

// #cgo pkg-config: libusb-1.0
// #include <libusb.h>
// #include <sys/time.h>
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
// static int libusb_handle_events_timeout_ms(libusb_context *ctx, int timeout_ms) {
//	struct timeval tv;
//	tv.tv_sec = timeout_ms / 1000;
//	tv.tv_usec = (timeout_ms % 1000) * 1000;
//	return libusb_handle_events_timeout_completed(ctx, &tv, NULL);
// }
import "C"
import (
	"fmt"
	"log"
	"sync"
	"unsafe"
)

// HotPlugEventType represents the type of hotplug event.
type HotPlugEventType uint8

// HotPlugCbFunc is the callback function signature for hotplug events.
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

// hotplugStorage holds the callback map and done channel for a single context.
type hotplugStorage struct {
	callbackMap map[uint32]hotplugCallback
	done        chan struct{}
	mu          sync.RWMutex
}

// hotplugEventTimeoutMs is the timeout in milliseconds for
// libusb_handle_events_timeout_completed in the event loop. This prevents
// busy-spinning while still being responsive to the done signal.
const hotplugEventTimeoutMs = 200

// hotplugRegistry maps context pointers to their hotplug storage, allowing
// multiple contexts to register hotplug callbacks independently.
var (
	hotplugRegistry   = make(map[*C.libusb_context]*hotplugStorage)
	hotplugRegistryMu sync.RWMutex
)

func getHotplugStorage(
	libCtx *C.libusb_context,
) *hotplugStorage {
	hotplugRegistryMu.RLock()
	defer hotplugRegistryMu.RUnlock()
	return hotplugRegistry[libCtx]
}

func (ctx *Context) newHotPlugHandler() *hotplugStorage {
	storage := &hotplugStorage{
		callbackMap: make(map[uint32]hotplugCallback),
		done:        make(chan struct{}),
	}

	hotplugRegistryMu.Lock()
	hotplugRegistry[ctx.libusbContext] = storage
	hotplugRegistryMu.Unlock()

	go storage.handleEvents(ctx.libusbContext)
	return storage
}

func (ctx *Context) getOrCreateHotplugStorage() *hotplugStorage {
	storage := getHotplugStorage(ctx.libusbContext)
	if storage == nil {
		storage = ctx.newHotPlugHandler()
	}
	return storage
}

func removeHotplugStorage(libCtx *C.libusb_context) {
	hotplugRegistryMu.Lock()
	delete(hotplugRegistry, libCtx)
	hotplugRegistryMu.Unlock()
}

// HotplugRegisterCallbackEvent registers a hotplug callback for the given
// vendor/product ID pair and event type.
func (ctx *Context) HotplugRegisterCallbackEvent(
	vendorID, productID uint16,
	eventType HotPlugEventType, cb HotPlugCbFunc,
) error {
	storage := ctx.getOrCreateHotplugStorage()

	var event C.int
	switch eventType {
	case HotplugArrived:
		event = C.LIBUSB_HOTPLUG_EVENT_DEVICE_ARRIVED
	case HotplugLeft:
		event = C.LIBUSB_HOTPLUG_EVENT_DEVICE_LEFT
	default:
		event = C.LIBUSB_HOTPLUG_EVENT_DEVICE_ARRIVED |
			C.LIBUSB_HOTPLUG_EVENT_DEVICE_LEFT
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
		C.libusb_hotplug_callback_fn(
			unsafe.Pointer(C.libusbHotplugCallback),
		),
		nil,
		&cbHandle,
	)
	if rc != C.LIBUSB_SUCCESS {
		return fmt.Errorf(
			"libusb_hotplug_register_callback error: %s", ErrorCode(rc),
		)
	}

	key := vidPidToUint32(vendorID, productID)
	storage.mu.Lock()
	storage.callbackMap[key] = hotplugCallback{
		handler: &cbHandle,
		fn:      cb,
	}
	storage.mu.Unlock()

	return nil
}

// HotplugDeregisterCallback deregisters a hotplug callback for the given
// vendor/product ID pair.
func (ctx *Context) HotplugDeregisterCallback(
	vendorID, productID uint16,
) error {
	storage := getHotplugStorage(ctx.libusbContext)
	if storage == nil {
		return nil
	}

	key := vidPidToUint32(vendorID, productID)

	storage.mu.RLock()
	cb, ok := storage.callbackMap[key]
	storage.mu.RUnlock()

	if !ok {
		return nil
	}

	C.libusb_hotplug_deregister_callback(ctx.libusbContext, *cb.handler)

	storage.mu.Lock()
	delete(storage.callbackMap, key)
	mapEmpty := len(storage.callbackMap) == 0
	storage.mu.Unlock()

	if mapEmpty {
		ctx.hotplugHandleEventsCompleteAll()
	}
	return nil
}

// HotplugDeregisterAllCallbacks deregisters all hotplug callbacks for this
// context and stops the event handler goroutine.
func (ctx *Context) HotplugDeregisterAllCallbacks() error {
	storage := getHotplugStorage(ctx.libusbContext)
	if storage == nil {
		return nil
	}

	storage.mu.RLock()
	handlers := make(
		[]*C.libusb_hotplug_callback_handle,
		0,
		len(storage.callbackMap),
	)
	for _, cb := range storage.callbackMap {
		handlers = append(handlers, cb.handler)
	}
	storage.mu.RUnlock()

	for _, handler := range handlers {
		C.libusb_hotplug_deregister_callback(ctx.libusbContext, *handler)
	}

	ctx.hotplugHandleEventsCompleteAll()

	return nil
}

func (ctx *Context) hotplugHandleEventsCompleteAll() {
	storage := getHotplugStorage(ctx.libusbContext)
	if storage == nil {
		return
	}

	// Signal the event handler goroutine to stop. Closing the channel
	// unblocks all receivers immediately without needing a separate send.
	close(storage.done)

	// Clean up the storage for this context.
	storage.mu.Lock()
	storage.callbackMap = nil
	storage.mu.Unlock()

	removeHotplugStorage(ctx.libusbContext)
}

func (storage *hotplugStorage) handleEvents(
	libCtx *C.libusb_context,
) {
	for {
		select {
		case <-storage.done:
			return
		default:
		}
		errno := C.libusb_handle_events_timeout_ms(
			libCtx, C.int(hotplugEventTimeoutMs),
		)
		if errno < 0 {
			if ErrorCode(errno) == errorInterrupted {
				continue
			}
			log.Printf("handle_events error: %s", ErrorCode(errno))
		}
	}
}

//export libusbHotplugCallback
func libusbHotplugCallback(
	ctx *C.libusb_context, dev *C.libusb_device,
	event C.libusb_hotplug_event, p unsafe.Pointer,
) C.int {
	var desc C.libusb_device_descriptor_struct
	rc := C.libusb_get_device_descriptor(dev, &desc)
	if rc != C.LIBUSB_SUCCESS {
		return rc
	}

	vendorID := uint16(desc.idVendor)
	productID := uint16(desc.idProduct)

	var e HotPlugEventType
	switch event {
	case C.LIBUSB_HOTPLUG_EVENT_DEVICE_ARRIVED:
		e = HotplugArrived
	case C.LIBUSB_HOTPLUG_EVENT_DEVICE_LEFT:
		e = HotplugLeft
	default:
		e = HotplugUndefined
	}

	storage := getHotplugStorage(ctx)
	if storage == nil {
		return C.LIBUSB_SUCCESS
	}

	storage.mu.RLock()
	cb, ok := storage.callbackMap[vidPidToUint32(vendorID, productID)]
	var deviceCallback HotPlugCbFunc
	if ok {
		deviceCallback = cb.fn
	}
	cb, ok = storage.callbackMap[0]
	var allCallback HotPlugCbFunc
	if ok {
		allCallback = cb.fn
	}
	storage.mu.RUnlock()

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
