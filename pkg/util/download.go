package util

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
	"io"

	"gopkg.in/cheggaaa/pb.v1"
)

// WriteCounter counts the number of bytes written to it.
type writeCounter struct {
	current int64 // Total # of bytes transferred
	total   int64 // Expected length
	bar     *pb.ProgressBar
}

func newWriteCounter(total int64) *writeCounter {
	return &writeCounter{
		total: total,
		bar:   pb.StartNew(int(total)),
	}
}

// Write implements the io.Writer interface.
//
// Always completes and never returns an error.
func (wc *writeCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.current += int64(n)
	wc.bar.Set(int(wc.current))
	if wc.current >= wc.total-1 {
		wc.bar.FinishPrint("download complete")
	}
	return n, nil
}

func ReaderWithProgress(r io.Reader, total int64) io.Reader {
	return io.TeeReader(r, newWriteCounter(total))
}
