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
	"fmt"
	"runtime/debug"
	"sort"
)

type AppEntry func(ctx *Context) error

var apps = map[string]AppEntry{}

func Register(name string, app AppEntry) {
	apps[name] = app
}

func Get(name string) AppEntry {
	return apps[name]
}

func AppNames() []string {
	var l []string
	for name := range apps {
		l = append(l, name)
	}
	sort.Strings(l)
	return l
}

func Run(name string, ctx *Context) error {
	entry := Get(name)
	if entry == nil {
		return fmt.Errorf("command not found: %s", name)
	}
	defer func() {
		err := recover()
		if err == nil {
			return
		}
		stack := debug.Stack()
		fmt.Fprintf(ctx.Stderr, "panic:%s\n", err)
		ctx.Stderr.Write(stack)
	}()
	return entry(ctx)
}
