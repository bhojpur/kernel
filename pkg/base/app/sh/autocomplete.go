package sh

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
	"path"
	"strings"

	"github.com/bhojpur/kernel/pkg/base/app"
	"github.com/spf13/afero"
)

func autocompleteWrapper(ctx *app.Context) func(line string) []string {
	return func(line string) []string {
		return autocomplete(ctx, line)
	}
}

func autocomplete(ctx *app.Context, line string) []string {
	line = strings.TrimLeft(line, " ")

	list := strings.Split(line, " ")

	var (
		last   = ""
		hascmd bool
		l      []string
	)

	if len(list) != 0 && list[0] == "go" {
		list = list[1:]
	}

	switch len(list) {
	case 0:
	case 1:
		last = list[0]
	default:
		hascmd = true
		last = list[len(list)-1]
	}

	if !hascmd {
		l = app.AppNames()
	} else {
		l = completeFile(ctx, last)
	}

	var r []string
	for _, s := range l {
		if strings.HasPrefix(s, last) {
			r = append(r, line+strings.TrimPrefix(s, last))
		}
	}

	return r
}

func completeFile(fs afero.Fs, prefix string) []string {
	if prefix == "" {
		prefix = "."
	}

	joinPrefix := func(dir string, l []string) []string {
		for i := range l {
			l[i] = path.Join(dir, l[i])
		}
		return l
	}

	f, err := fs.Open(prefix)
	// user input a complete file name
	if err == nil {
		defer f.Close()

		stat, err := f.Stat()
		if err != nil {
			return nil
		}
		if !stat.IsDir() {
			return nil
		}
		names, err := f.Readdirnames(-1)
		if err != nil {
			return nil
		}
		return joinPrefix(prefix, names)
	}

	// complete dir entries
	dir := path.Dir(prefix)
	f, err = fs.Open(dir)
	if err != nil {
		return nil
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil
	}
	return joinPrefix(dir, names)
}
