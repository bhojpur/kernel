# Bhojpur Kernel - Builder Tool

The `builder` is a make-like build tool using Go. You write plain-old Go functions, and
`builder` automatically uses them as Makefile-like runnable targets.

## Installation

The `builder` has no dependencies outside the Go standard library, and builds with Go 1.7
and above.

**Using GOPATH**

```
go get -u -d github.com/bhojpur/kernel/pkg/builder
cd $GOPATH/src/github.com/bhojpur/kernel/pkg/builder
go run bootstrap.go
```

**Using Go Modules**

```
git clone https://github.com/bhojpur/kernel
cd pkg/builder
go run bootstrap.go
```

It will download the source code and then run the `bootstrap` script to build `builder`
with version infomation embedded in it. A normal `go get` (without -d) or `go install`
will build the binary correctly, but no version info will be embedded. If you've done
this, no worries, just go to `$GOPATH/src/github.com/bhojpur/kernel/pkg/builder` and
run `builder install` or `go run bootstrap.go` and a new binary will be created with
the correct version information.

The `builder` binary will be created in your $GOPATH/bin directory.
