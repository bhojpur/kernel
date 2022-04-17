package app

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
	"flag"
	"fmt"
	"io"

	"github.com/bhojpur/kernel/pkg/base/fs"
	"github.com/bhojpur/kernel/pkg/base/fs/chdir"
	"github.com/peterh/liner"
)

type Context struct {
	Args   []string
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer

	*chdir.Chdirfs

	flag  *flag.FlagSet
	liner *liner.State
}

func (c *Context) Init() {
	c.Chdirfs = chdir.New(fs.Root)
}

func (c *Context) Printf(fmtstr string, args ...interface{}) {
	fmt.Fprintf(c.Stdout, fmtstr, args...)
}

func (c *Context) Flag() *flag.FlagSet {
	if c.flag != nil {
		return c.flag
	}
	c.flag = flag.NewFlagSet(c.Args[0], flag.ContinueOnError)
	return c.flag
}

func (c *Context) ParseFlags() error {
	return c.Flag().Parse(c.Args[1:])
}

func (c *Context) LineReader() LineReader {
	_, ok := c.Stdin.(fs.Ioctler)
	if !ok {
		return newSimpleLineReader(c.Stdin, c.Stdout)
	}
	return newLineEditor()
}
