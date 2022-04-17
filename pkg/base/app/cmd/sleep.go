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
	"flag"
	"fmt"
	"time"

	"github.com/bhojpur/kernel/pkg/base/app"
)

func sleepmain(ctx *app.Context) error {
	var (
		flagset = flag.NewFlagSet(ctx.Args[0], flag.ContinueOnError)
		istick  = flagset.Bool("t", false, "use ticker")
	)
	err := flagset.Parse(ctx.Args[1:])
	if err != nil {
		return err
	}
	if len(flagset.Args()) == 0 {
		return errors.New("usage: sleep $duration")
	}
	dura, err := time.ParseDuration(flagset.Arg(0))
	if err != nil {
		return err
	}

	if !*istick {
		time.Sleep(dura)
		return nil
	}

	begin := time.Now()
	ticker := time.NewTicker(dura)
	for t := range ticker.C {
		fmt.Fprintln(ctx.Stdout, t.Sub(begin).Milliseconds())
	}
	return nil
}

func init() {
	app.Register("sleep", sleepmain)
}
