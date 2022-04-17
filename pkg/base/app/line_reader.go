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
	"bufio"
	"fmt"
	"io"

	"github.com/peterh/liner"
)

type LineReader interface {
	Prompt(string) (string, error)
	AppendHistory(string)
	SetAutoComplete(func(string) []string)
	Close() error
}

type simpleLineReader struct {
	r *bufio.Scanner
	w io.Writer
}

func newSimpleLineReader(r io.Reader, w io.Writer) LineReader {
	return &simpleLineReader{
		r: bufio.NewScanner(r),
		w: w,
	}
}

func (r *simpleLineReader) Prompt(prompt string) (string, error) {
	fmt.Fprintf(r.w, "%s", prompt)
	ok := r.r.Scan()
	if !ok {
		return "", io.EOF
	}
	return r.r.Text(), nil
}

func (r *simpleLineReader) AppendHistory(string) {
}

func (r *simpleLineReader) SetAutoComplete(f func(string) []string) {
}

func (r *simpleLineReader) Close() error {
	return nil
}

type lineEditor struct {
	*liner.State
}

func (l lineEditor) SetAutoComplete(f func(string) []string) {
	l.SetCompleter(f)
}

func newLineEditor() LineReader {
	r := liner.NewLiner()
	return lineEditor{
		State: r,
	}
}
