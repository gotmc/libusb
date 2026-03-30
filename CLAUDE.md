# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Go (cgo) bindings for the [libusb C library](http://libusb.info) (libusb-1.0). Module path: `github.com/gotmc/libusb/v2`. Requires the libusb C library installed on the system (`brew install libusb` on macOS, `apt-get install libusb-dev libusb-1.0-0-dev` on Linux).

## Build & Test Commands

This project uses [Just](https://github.com/casey/just) as a command runner.

- `just check` — format (`go fmt`) and vet (`go vet`)
- `just lint` — lint with golangci-lint (config in `.golangci.yaml`)
- `just unit` — run unit tests (runs `check` first, uses `-short -race -cover`)
- `just unit -v` — verbose unit tests
- `just int` — run integration tests (tests matching `Integration`)
- `just e2e` — run end-to-end tests (tests matching `E2E`)
- `just cover` — generate and open HTML coverage report

Run a single test directly: `go test -run TestName -v`

## Architecture

This is a single Go package (`package libusb`) that wraps the libusb-1.0 C library via cgo. Every `.go` file uses `#cgo pkg-config: libusb-1.0` and `#include <libusb.h>`.

Key types mirror the libusb C API:

- **Context** (`context.go`) — wraps `libusb_context`. Entry point: `NewContext()` returns a session. Provides `DeviceList()` and `OpenDeviceWithVendorProduct()`.
- **Device** (`device.go`) — wraps `libusb_device`. Obtained from Context. Provides descriptors (`DeviceDescriptor()`, `ActiveConfigDescriptor()`, `ConfigDescriptor()`). Uses `runtime.SetFinalizer` for safety with `libusb_ref_device`/`libusb_unref_device`.
- **DeviceHandle** (`handle.go`) — wraps `libusb_device_handle`. Obtained from `Device.Open()`. Provides interface claiming, kernel driver management, and string descriptors. Also uses finalizers.
- **Synchronous I/O** (`syncio.go`) — `BulkTransfer`, `ControlTransfer`, `InterruptTransfer` and helper variants on DeviceHandle.
- **Async I/O** (`asyncio.go`) — asynchronous transfer support.
- **Hotplug** (`hotplug.go`) — hotplug callback registration via cgo export.
- **Setup Data** (`setupdata.go`) — USB request type/direction/recipient constants and `BitmapRequestType` helper.
- **Descriptors** (`descriptors.go`) — internal types (`bcd`, `descriptorType`, `classCode`, `endpointAddress`, etc.).
- **Interfaces** (`interfaces.go`) — `SupportedInterface`, `InterfaceDescriptor`, interface class constants, class-based filtering.
- **Configuration** (`configuration.go`) — `ConfigDescriptor` and `EndpointDescriptor` types.
- **Error handling** (`miscellaneous.go`) — `ErrorCode` type wrapping libusb error codes; `ErrorName()`, `StrError()`.

## Conventions

- C memory management follows libusb conventions: `libusb_ref_device`/`libusb_unref_device` paired with Go finalizers as a safety net.
- Types that wrap C pointers (`Device`, `DeviceHandle`) have explicit `Close()` methods and `runtime.SetFinalizer` for GC cleanup.
- Tests are in `*_test.go` files alongside source. Test names use `Test` prefix for unit, `Integration` for integration, `E2E` for end-to-end.
