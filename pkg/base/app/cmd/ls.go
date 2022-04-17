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
	"io"
	"os"
	"text/tabwriter"

	"github.com/bhojpur/kernel/pkg/base/app"
)

func printfiles(w io.Writer, files ...os.FileInfo) {
	tw := tabwriter.NewWriter(w, 0, 4, 1, ' ', 0)
	for _, file := range files {
		fmt.Fprintf(tw, "%s\t%d\t%s\n", file.Mode(), file.Size(), file.Name())
	}
	tw.Flush()
}

func lsmain(ctx *app.Context) error {
	err := ctx.ParseFlags()
	if err != nil {
		return err
	}

	var name string
	if ctx.Flag().NArg() > 0 {
		name = ctx.Flag().Arg(0)
	}
	f, err := ctx.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	stat, err := f.Stat()
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		printfiles(ctx.Stdout, stat)
		return nil
	}
	stats, err := f.Readdir(-1)
	if err != nil {
		return err
	}
	printfiles(ctx.Stdout, stats...)
	return nil
}

func init() {
	app.Register("ls", lsmain)
}
