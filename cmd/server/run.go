package cmd

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
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/bhojpur/kernel/cmd/server/assets"
	"github.com/bhojpur/kernel/cmd/server/build"
	"github.com/google/shlex"
	"github.com/spf13/cobra"
)

const (
	qemu64 = "qemu-system-x86_64"
)

var (
	ports []string
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run <kernel>",
	Short: "run running a Bhojpur Kernel in qemu",
	Run: func(cmd *cobra.Command, args []string) {
		err := runKernel(args)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func runKernel(args []string) error {
	base, err := ioutil.TempDir("", "bhojpur-run")
	if err != nil {
		return err
	}
	defer os.RemoveAll(base)

	var kernelFile string

	if len(args) == 0 || args[0] == "" {
		kernelFile = filepath.Join(base, "kernel.elf")

		b := build.NewBuilder(build.Config{
			GoRoot:        goroot,
			Basedir:       base,
			BuildTest:     false,
			KernelVersion: kernelVersion,
			GoArgs: []string{
				"-o", kernelFile,
			},
		})
		if err := b.Build(); err != nil {
			return fmt.Errorf("error building Bhojpur Kernel: %s", err)
		}
	} else {
		kernelFile = args[0]
	}

	var runArgs []string

	ext := filepath.Ext(kernelFile)
	switch ext {
	case ".elf", "":
		loaderFile := filepath.Join(base, "loader.elf")
		mustLoaderFile(loaderFile)
		runArgs = append(runArgs, "-kernel", loaderFile)
		runArgs = append(runArgs, "-initrd", kernelFile)
	case ".iso":
		runArgs = append(runArgs, "-cdrom", kernelFile)
	}

	var qemuArgs []string
	if qemuArgs, err = shlex.Split(os.Getenv("QEMU_OPTS")); err != nil {
		return fmt.Errorf("error parsing QEMU_OPTS: %s", err)
	}

	runArgs = append(runArgs, "-m", "256M", "-no-reboot", "-serial", "mon:stdio")
	runArgs = append(runArgs, "-netdev", "user,id=eth0"+portMapingArgs())
	runArgs = append(runArgs, "-device", "e1000,netdev=eth0")
	runArgs = append(runArgs, "-device", "isa-debug-exit")
	runArgs = append(runArgs, qemuArgs...)

	cmd := exec.Command(qemu64, runArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err == nil {
		return nil
	}
	switch e := err.(type) {
	case *exec.ExitError:
		code := e.ExitCode()
		if code == 0 || code == 1 {
			return nil
		}
		return err
	default:
		return err
	}
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func mustLoaderFile(fname string) {
	content, err := assets.Boot.ReadFile("boot/multiboot.elf")
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile(fname, content, 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func portMapingArgs() string {
	var ret []string
	for _, mapping := range ports {
		fs := strings.Split(mapping, ":")
		if len(fs) < 2 {
			continue
		}
		arg := fmt.Sprintf(",hostfwd=tcp::%s-:%s", fs[0], fs[1])
		ret = append(ret, arg)
	}
	return strings.Join(ret, "")
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringSliceVarP(&ports, "port", "p", nil, "port mapping from host to Bhojpur Kernel, format $host_port:$kernel_port")
}
