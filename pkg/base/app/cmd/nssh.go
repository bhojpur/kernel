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
	"net"

	"github.com/bhojpur/kernel/pkg/base/app"
)

func nsshmain(ctx *app.Context) error {
	l, err := net.Listen("tcp", "0.0.0.0:22")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go func() {
			fmt.Fprintf(ctx.Stdout, "conn from:%s\n", conn.RemoteAddr())
			shell := app.Get("sh")
			ctx := &app.Context{
				Stdin:  conn,
				Stdout: conn,
				Stderr: conn,
			}
			ctx.Init()
			shell(ctx)
			conn.Close()
			fmt.Fprintf(ctx.Stdout, "conn %s closed\n", conn.RemoteAddr())
		}()
	}
}

func init() {
	app.Register("nssh", nsshmain)
}
