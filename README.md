# Bhojpur Kernel - Library Operating System

The `Bhojpur Kernel` is an `operating system` as a library (i.e., Unikernel framework) used by
the [Bhojpur.NET Platform](https://github.com/bhojpur/platform) ecosystem for delivery of fast
and secure `applications` or `services`. It has the ability to convert normal `Go` program into
an ELF unikernel, which could be run on a bare metal. It is a tool for compiling application
sources into `Unikernel` (i.e., lightweight bootable disk images) and `MicroVM` rather than
binaries.

## Unikernel Framework

The `Bhojpur Kernel` runs and manages instances of compiled images across a variety of cloud
providers as well as locally. It utilizes a simple `Docker`-like command line interface (i.e.,
`kernctl`), making building `Unikernel` and `MicroVM` as easy as building containers.

### Supported Providers

* [Firecracker](https://firecracker-microvm.github.io/)
* Virtualbox
* AWS
* Google Cloud
* vSphere
* QEMU
* UKVM
* Xen
* OpenStack
* Photon Controller

The `Instance Listener` is a special component of th eBhojpur Kernel that bootstraps unikernel
instances running on certain providers (currently `vSphere` and `VirtualBox`).

### Build Unikernel Image

To start with, you need `QEMU` to run it locally on a macOS.

```bash
$ brew install x86_64-elf-binutils x86_64-elf-gcc x86_64-elf-gdb
$ brew install qemu
```

You can start using the following command

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

### Unikernel Types Supported

* **Firecracker**: the `kernctl` supports compiling Go source code into [Firecracker](https://firecracker-microvm.github.io/) MicroVMs
* **rump**: the `kernctl` supports compiling Python, Node.js, and Go source code into [rumprun](https://github.com/rumpkernel) unikernels
* **OSv**: the `kernctl` supports compiling Java, Node.js, C and C++ source code into [OSv](http://osv.io/) unikernels
* **IncludeOS**: the `kernctl` supports compiling C++ source code into [IncludeOS](https://github.com/hioa-cs/IncludeOS) unikernels
* **MirageOS**: the `kernctl` supports compiling OCaml source code into [MirageOS](https://mirage.io) unikernels

### Source Code Build

The `Bhojpur Kernel` is compiled into binary image by the `builder` tool using the configuration settings
defined in the `builderfile.go` file.

```bash
$ go mod tidy
$ go get
$ builder iso
```
