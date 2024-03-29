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
	"log"

	"github.com/bhojpur/kernel/cmd/server/build"
	"github.com/spf13/cobra"
)

var (
	output string
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:                "build",
	Short:              "build like Go build command",
	DisableFlagParsing: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := runBuild(args)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func runBuild(args []string) error {
	b := build.NewBuilder(build.Config{
		GoRoot:        goroot,
		KernelVersion: kernelVersion,
		GoArgs:        args,
	})
	err := b.Build()
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
