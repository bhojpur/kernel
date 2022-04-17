# Bhojpur Kernel - Operating System Library

The `Bhojpur Kernel` is an operating system as a library (i.e., Unikernel framework) used by the
[Bhojpur.NET Platform](https://github.com/bhojpur/platform) ecosystem for delivery of fast and
secure applications or services. It has the ability to convert normal `Go` program into an ELF
unikernel, which could be run on a bare metal.

## Unikernel Framework

You need `QEMU` to run it on the macOS locally.

```bash
$ brew install x86_64-elf-binutils x86_64-elf-gcc x86_64-elf-gdb
$ brew install qemu
```

You can start using following command

```bash
$ builder qemu
```

Firstly, get the `kernel` command line tool.

```
$ go install github.com/bhojpur/kernel
```

Run the following command to build your `Unikernel` application.

```
$ kernel build -o kernel.elf
```

The following command can package your custom application into a `Unikernel` ISO file.

```bash
$ kernel pack -o bhojpur-kernel.iso -k kernel.elf
```

Then, you can use [https://github.com/ventoy/Ventoy](https://github.com/ventoy/Ventoy) to
run the ISO file on a bare metal compute server.

### Source Code Build

The `Bhojpur Kernel` is compiled into binary image by the `builder` tool using the configuration settings
defined in the `builderfile.go` file.

```bash
$ go mod tidy
$ go get
$ builder iso
```
