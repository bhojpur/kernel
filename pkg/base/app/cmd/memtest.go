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
	"math/rand"
	"time"

	"github.com/bhojpur/kernel/pkg/base/app"
)

func memtestmain(ctx *app.Context) error {
	if len(ctx.Args) < 2 {
		fmt.Fprintln(ctx.Stderr, "usage: memtest $duration")
		return nil
	}

	dura, err := time.ParseDuration(ctx.Args[1])
	if err != nil {
		return err
	}
	deadline := time.Now().Add(dura)
	for {
		if time.Now().After(deadline) {
			return nil
		}

		buf := make([]byte, 1024)
		rand.Read(buf)
		fmt.Fprintf(ioutil.Discard, "%v", buf)
	}
}

func init() {
	app.Register("memtest", memtestmain)
}
