# libusb

Go bindings for the [libusb C library][libusb-c].

[![GoDoc][godoc badge]][godoc link]
[![Go Report Card][report badge]][report card]
[![Build Status][travis image]][travis link]
[![License Badge][license image]][LICENSE.txt]

# Installation

```bash
$ go get github.com/gotmc/libusb/v2
```

## Installing C libusb library

To use [libusb][] package, you'll need to install the [libusb C
library][libusb-c] first.

### OS X

```bash
$ brew install libusb
```

### Windows

Download and install the latest Windows libusb binaries from
[libusb.info][libusb-c].

### Linux

```bash
$ sudo apt-get install -y libusb-dev libusb-1.0-0-dev
```

# Documentation

Documentation can be found at either:

- <https://godoc.org/github.com/gotmc/libusb>
- <http://localhost:6060/pkg/github.com/gotmc/libusb/> after running `$
godoc -http=:6060`

## Contributing

Contributions are welcome! To contribute please:

1. Fork the repository
2. Create a feature branch
3. Code
4. Submit a [pull request][]

### Testing

Prior to submitting a [pull request][], please run:

```bash
$ make check
$ make lint
```

To update and view the test coverage report:

```bash
$ make cover
```

## Alternatives

There are other USB Go libraries besides [libusb][]. Below are a few
alternatives:

- [google/gousb][] — Wraps the [libusb C library][libusb-c] to provde
  Go-bindings. This library supersedes [kylelemons/gousb][], which is not
  archived. Apachage-2.0 license.
- [karalabe/usb][] — Does not require the [libusb C library][libusb-c] to be
  installed. Written in C to be a cross platform, fully self-contained library
  for accessing and communicating with USB devices either via HID or low level
  interrupts. LGPL-3.0 license.
- [deadsy/libusb][] — Wraps the [libusb C library][libusb-c]. MIT license. As of
  05-Apr-24, this package hasn't been updated in six years.

## License

[libusb][] is released under the MIT license. Please see the
[LICENSE.txt][] file for more information.

[deadsy/libusb]: https://github.com/deadsy/libusb
[godoc badge]: https://godoc.org/github.com/gotmc/libusb?status.svg
[godoc link]: https://godoc.org/github.com/gotmc/libusb
[google/gousb]: https://github.com/google/gousb
[karalabe/usb]: https://github.com/karalabe/usb
[kylelemons/gousb]: https://github.com/kylelemons/gousb
[libusb]: https://github.com/gotmc/libusb
[libusb-c]: http://libusb.info
[LICENSE.txt]: https://github.com/gotmc/libusb/blob/master/LICENSE.txt
[license image]: https://img.shields.io/badge/license-MIT-blue.svg
[pull request]: https://help.github.com/articles/using-pull-requests
[report badge]: https://goreportcard.com/badge/github.com/gotmc/libusb
[report card]: https://goreportcard.com/report/github.com/gotmc/libusb
[travis image]: http://img.shields.io/travis/gotmc/libusb/master.svg
[travis link]: https://travis-ci.org/gotmc/libusb
