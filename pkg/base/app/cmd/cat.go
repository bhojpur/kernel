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
	"errors"
	"fmt"
	"io"

	"github.com/bhojpur/kernel/pkg/base/app"
	"github.com/spf13/afero"
)

func catmain(ctx *app.Context) error {
	err := ctx.ParseFlags()
	if err != nil {
		return err
	}

	if ctx.Flag().NArg() == 0 {
		return errors.New("usage: cat $filename")
	}
	name := ctx.Flag().Arg(0)
	if isdir, _ := afero.IsDir(ctx, name); isdir {
		return errors.New("can't cat dir")
	}

	f, err := ctx.Open(name)
	if err != nil {
		return err
	}
	defer f.Close()

	io.Copy(ctx.Stdout, f)
	fmt.Fprintln(ctx.Stdout)
	return nil
}

func init() {
	app.Register("cat", catmain)
}
