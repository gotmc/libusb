# libusb

Go bindings for the [libusb C library][libusb-c].

[![GoDoc][godoc image]][godoc link]
[![Build Status][travis image]][travis link]
[![License Badge][license image]][LICENSE.txt]

## Installation

```bash
$ go get github.com/gotmc/libusb
```

### Installing C libusb library

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
$ sudo apt-get install -y libusb-dev libusb-1.0-0 libusb-1.0-0-dev
```

## Documentation

Documentation can be found at either:

- <https://godoc.org/github.com/gotmc/libusb>
- <http://localhost:6060/pkg/github.com/gotmc/libusb/> after running `$
  godoc -http=:6060`

## Contributing

[libusb][] is developed using [Scott Chacon][]'s [GitHub Flow][]. To
contribute, fork [libusb][], create a feature branch, and then
submit a [pull request][].  [GitHub Flow][] is summarized as:

- Anything in the `master` branch is deployable
- To work on something new, create a descriptively named branch off of
  `master` (e.g., `new-oauth2-scopes`)
- Commit to that branch locally and regularly push your work to the same
  named branch on the server
- When you need feedback or help, or you think the branch is ready for
  merging, open a [pull request][].
- After someone else has reviewed and signed off on the feature, you can
  merge it into master.
- Once it is merged and pushed to `master`, you can and *should* deploy
  immediately.

## Testing

Prior to submitting a [pull request][], please run:

```bash
$ gofmt
$ golint
$ go vet
$ go test
```

To update and view the test coverage report:

```bash
$ go test -coverprofile coverage.out
$ go tool cover -html coverage.out
```

## License

[libusb][] is released under the MIT license. Please see the
[LICENSE.txt][] file for more information.

[GitHub Flow]: http://scottchacon.com/2011/08/31/github-flow.html
[godoc image]: https://godoc.org/github.com/gotmc/libusb?status.svg
[godoc link]: https://godoc.org/github.com/gotmc/libusb
[libusb]: https://github.com/gotmc/libusb
[libusb-c]: http://libusb.info
[LICENSE.txt]: https://github.com/gotmc/libusb/blob/master/LICENSE.txt
[license image]: https://img.shields.io/badge/license-MIT-blue.svg
[pull request]: https://help.github.com/articles/using-pull-requests
[Scott Chacon]: http://scottchacon.com/about.html
[travis image]: http://img.shields.io/travis/gotmc/libusb/master.svg
[travis link]: https://travis-ci.org/gotmc/libusb
