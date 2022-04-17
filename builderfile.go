//go:build builder
// +build builder

package main

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/bhojpur/kernel/pkg/builder/sh"
	"github.com/bhojpur/kernel/pkg/builder/utils"
)

var (
	TOOLPREFIX = detectToolPrefix()
	CC         = TOOLPREFIX + "gcc"
	LD         = TOOLPREFIX + "ld"

	CFLAGS  = initCflags()
	LDFLAGS = initLdflags()
)

var (
	GOTAGS    = "nes phy prometheus"
	GOGCFLAGS = ""
)

var (
	QEMU64 = "qemu-system-x86_64"
	QEMU32 = "qemu-system-i386"

	QEMU_OPT       = initQemuOpt()
	QEMU_DEBUG_OPT = initQemuDebugOpt()
)

var (
	kernelBin string
)

const (
	goMajorVersionSupported    = 1
	maxGoMinorVersionSupported = 17
)

// Kernel target build the ELF kernel file for Bhojpur Kernel, generate kernel.elf
func Kernel() error {
	utils.Deps(BhojpurKernel)

	detectGoVersion()
	return rundir("app", nil, kernelBin, "build", "-o", "../kernel.elf",
		"-gcflags", GOGCFLAGS,
		"-tags", GOTAGS,
		"./kmain")
}

func Boot64() error {
	compileCfile("boot/boot64.S", "-m64")
	compileCfile("boot/boot64main.c", "-m64")
	ldflags := "-Ttext 0x3200000 -m elf_x86_64 -o boot64.elf boot64.o boot64main.o"
	ldArgs := append([]string{}, LDFLAGS...)
	ldArgs = append(ldArgs, strings.Fields(ldflags)...)
	return sh.RunV(LD, ldArgs...)
}

// Multiboot target build Multiboot specification compatible elf format, generate multiboot.elf
func Multiboot() error {
	utils.Deps(Boot64)
	compileCfile("boot/multiboot.c", "-m32")
	compileCfile("boot/multiboot_header.S", "-m32")
	ldflags := "-Ttext 0x3300000 -m elf_i386 -o multiboot.elf multiboot.o multiboot_header.o -b binary boot64.elf"
	ldArgs := append([]string{}, LDFLAGS...)
	ldArgs = append(ldArgs, strings.Fields(ldflags)...)
	err := sh.RunV(LD, ldArgs...)
	if err != nil {
		return err
	}
	return sh.Copy(
		filepath.Join("cmd", "server", "assets", "boot", "multiboot.elf"),
		"multiboot.elf",
	)
}

func Test() error {
	utils.Deps(BhojpurKernel)

	envs := map[string]string{
		"QEMU_OPTS": quoteArgs(QEMU_OPT),
	}
	return rundir("tests", envs, kernelBin, "test")
}

func TestDebug() error {
	utils.Deps(BhojpurKernel)

	envs := map[string]string{
		"QEMU_OPTS": quoteArgs(QEMU_DEBUG_OPT),
	}
	return rundir("tests", envs, kernelBin, "test")
}

// Qemu run multiboot.elf on qemu.
// If env QEMU_ACCEL is set，QEMU acceleration will be enabled.
// If env QEMU_GRAPHIC is set QEMU will run in graphic mode.
// Use Crtl+a c to switch console, and type `quit`
func Qemu() error {
	utils.Deps(Kernel)

	detectQemu()
	return kernelrun(QEMU_OPT, "kernel.elf")
}

// QemuDebug run multiboot.elf in debug mode.
// Monitor GDB connection on port 1234
func QemuDebug() error {
	GOGCFLAGS += " -N -l"
	utils.Deps(Kernel)

	detectQemu()
	return kernelrun(QEMU_DEBUG_OPT, "kernel.elf")
}

// Iso generate bhojpur-kernel.iso, which can be used with qemu -cdrom option.
func Iso() error {
	utils.Deps(Kernel)
	return sh.RunV(kernelBin, "pack", "-o", "bhojpur-kernel.iso", "-k", "kernel.elf")
}

// Graphic run bhojpur-kernel.iso on qemu, which vbe is enabled.
func Graphic() error {
	detectQemu()

	utils.Deps(Iso)
	return kernelrun(QEMU_OPT, "bhojpur-kernel.iso")
}

// GraphicDebug run bhojpur-kernel.iso on qemu in debug mode.
func GraphicDebug() error {
	detectQemu()

	GOGCFLAGS += " -N -l"
	utils.Deps(Iso)
	return kernelrun(QEMU_DEBUG_OPT, "bhojpur-kernel.iso")
}

func BhojpurKernel() error {
	err := rundir("cmd", nil, "go", "build", "-o", "../kernel", "./kernel")
	if err != nil {
		return err
	}
	current, _ := os.Getwd()
	kernelBin = filepath.Join(current, "kernel")
	return nil
}

func Clean() {
	rmGlob("*.o")
	rmGlob("kernel.elf")
	rmGlob("multiboot.elf")
	rmGlob("qemu.log")
	rmGlob("qemu.pcap")
	rmGlob("bhojpur-kernel.iso")
	rmGlob("kernel")
	rmGlob("boot64.elf")
	rmGlob("bochs.log")
}

func detectToolPrefix() string {
	prefix := os.Getenv("TOOLPREFIX")
	if prefix != "" {
		return prefix
	}

	if hasOutput("elf32-i386", "x86_64-elf-objdump", "-i") {
		return "x86_64-elf-"
	}

	if hasOutput("elf32-i386", "i386-elf-objdump", "-i") {
		return "i386-elf-"
	}

	if hasOutput("elf32-i386", "objdump", "-i") {
		return ""
	}
	panic(`
	*** Error: Couldn't find an i386-*-elf or x86_64-*-elf version of GCC/binutils
	*** Is the directory with i386-elf-gcc or x86_64-elf-gcc in your PATH?
	*** If your i386/x86_64-*-elf toolchain is installed with a command
	*** prefix other than 'i386/x86_64-elf-', set your TOOLPREFIX
	*** environment variable to that prefix and run 'make' again.
	`)
}

var goVersionRegexp = regexp.MustCompile(`go(\d+)\.(\d+)\.?(\d?)`)

func goVersion() (string, int, int, error) {
	versionBytes, err := cmdOutput(gobin(), "version")
	if err != nil {
		return "", 0, 0, err
	}
	version := strings.TrimSpace(string(versionBytes))
	result := goVersionRegexp.FindStringSubmatch(version)
	if len(result) < 3 {
		return "", 0, 0, fmt.Errorf("use of unreleased Go version `%s`, may not work", version)
	}
	major, _ := strconv.Atoi(result[1])
	minor, _ := strconv.Atoi(result[2])
	return version, major, minor, nil
}

func detectGoVersion() {
	version, major, minor, err := goVersion()
	if err != nil {
		fmt.Printf("warning: %s\n", err)
		return
	}
	if !(major == goMajorVersionSupported && minor <= maxGoMinorVersionSupported) {
		fmt.Printf("warning: max supported Go version go%d.%d.x, found Go version `%s`, may not work\n",
			goMajorVersionSupported, maxGoMinorVersionSupported, version,
		)
		return
	}
}

func gobin() string {
	goroot := os.Getenv("BHOJPUR_KERNEL_GOROOT")
	if goroot != "" {
		return filepath.Join(goroot, "bin", "go")
	}
	return "go"
}

func detectQemu() {
	if !hasCommand(QEMU64) {
		panic(QEMU64 + ` command not found`)
	}
}

func accelArg() []string {
	switch runtime.GOOS {
	case "darwin":
		return []string{"-M", "accel=hvf"}
	default:
		// fmt.Printf("accel method not found")
		return nil
	}
}

func initCflags() []string {
	cflags := strings.Fields("-fno-pic -static -fno-builtin -fno-strict-aliasing -O2 -Wall -Werror -fno-omit-frame-pointer -I. -nostdinc")
	if hasOutput("-fno-stack-protector", CC, "--help") {
		cflags = append(cflags, "-fno-stack-protector")
	}
	if hasOutput("[^f]no-pie", CC, "-dumpspecs") {
		cflags = append(cflags, "-fno-pie", "-no-pie")
	}
	if hasOutput("[^f]nopie", CC, "-dumpspecs") {
		cflags = append(cflags, "-fno-pie", "-nopie")
	}
	return cflags
}

func initLdflags() []string {
	ldflags := strings.Fields("-N -e _start")
	return ldflags
}

func initQemuOpt() []string {
	var opts []string
	if os.Getenv("QEMU_ACCEL") != "" {
		opts = append(opts, accelArg()...)
	}
	if os.Getenv("QEMU_GRAPHIC") == "" {
		opts = append(opts, "-nographic")
	}
	return opts
}

func initQemuDebugOpt() []string {
	opts := `
	-d int -D qemu.log
	-object filter-dump,id=f1,netdev=eth0,file=qemu.pcap
	-s -S
	`
	ret := append([]string{}, initQemuOpt()...)
	ret = append(ret, strings.Fields(opts)...)
	return ret
}

func compileCfile(file string, extFlags ...string) {
	args := append([]string{}, CFLAGS...)
	args = append(args, extFlags...)
	args = append(args, "-c", file)
	err := sh.RunV(CC, args...)
	if err != nil {
		panic(err)
	}
}

func rundir(dir string, envs map[string]string, cmd string, args ...string) error {
	current, _ := os.Getwd()
	os.Chdir(dir)
	defer os.Chdir(current)
	return sh.RunWithV(envs, cmd, args...)
}

func kernelrun(qemuArgs []string, flags ...string) error {
	qemuOpts := quoteArgs(qemuArgs)
	var args []string
	args = append(args, "run")
	args = append(args, "-p", "8080:80")
	args = append(args, flags...)
	envs := map[string]string{
		"QEMU_OPTS": qemuOpts,
	}
	return sh.RunWithV(envs, kernelBin, args...)
}

func cmdOutput(cmd string, args ...string) ([]byte, error) {
	return exec.Command(cmd, args...).CombinedOutput()
}

// quote string which has spaces with ""
func quoteArgs(args []string) string {
	var ret []string
	for _, s := range args {
		if strings.Index(s, " ") != -1 {
			ret = append(ret, strconv.Quote(s))
		} else {
			ret = append(ret, s)
		}
	}
	return strings.Join(ret, " ")
}

func hasCommand(cmd string) bool {
	_, err := exec.LookPath(cmd)
	if err != nil {
		return false
	}
	return true
}

func hasOutput(regstr, cmd string, args ...string) bool {
	out, err := cmdOutput(cmd, args...)
	if err != nil {
		return false
	}
	match, err := regexp.Match(regstr, []byte(out))
	if err != nil {
		return false
	}
	return match
}

func rmGlob(patten string) error {
	match, err := filepath.Glob(patten)
	if err != nil {
		return err
	}
	for _, file := range match {
		err = os.Remove(file)
		if err != nil {
			return err
		}
	}
	return nil
}
