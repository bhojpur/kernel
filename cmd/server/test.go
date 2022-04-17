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
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/bhojpur/kernel/cmd/server/build"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "tests likes Go test, but running in QEMU",
	Run: func(cmd *cobra.Command, args []string) {
		err := runTest()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func runTest() error {
	base, err := ioutil.TempDir("", "bhojpur-test")
	if err != nil {
		return err
	}
	defer os.RemoveAll(base)

	outfile := filepath.Join(base, "kernel.test.elf")

	b := build.NewBuilder(build.Config{
		GoRoot:        goroot,
		Basedir:       base,
		BuildTest:     true,
		KernelVersion: kernelVersion,
		GoArgs: []string{
			"-o", outfile,
			"-vet=off",
		},
	})
	err = b.Build()
	if err != nil {
		return err
	}

	return runKernel([]string{outfile})
}

func init() {
	rootCmd.AddCommand(testCmd)
}
