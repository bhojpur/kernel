package target

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
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewestModTime(t *testing.T) {
	t.Parallel()
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("error creating temp dir: %s", err.Error())
	}
	defer os.RemoveAll(dir)
	for _, name := range []string{"a", "b", "c", "d"} {
		out := filepath.Join(dir, name)
		if err := ioutil.WriteFile(out, []byte("hi!"), 0644); err != nil {
			t.Fatalf("error writing file: %s", err.Error())
		}
	}
	time.Sleep(10 * time.Millisecond)
	outName := filepath.Join(dir, "c")
	outfh, err := os.OpenFile(outName, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		t.Fatalf("error opening file to append: %s", err.Error())
	}
	if _, err := outfh.WriteString("\nbye!\n"); err != nil {
		t.Fatalf("error appending to file: %s", err.Error())
	}
	if err := outfh.Close(); err != nil {
		t.Fatalf("error closing file: %s", err.Error())
	}

	afi, err := os.Stat(filepath.Join(dir, "a"))
	if err != nil {
		t.Fatalf("error stating unmodified file: %s", err.Error())
	}

	cfi, err := os.Stat(outName)
	if err != nil {
		t.Fatalf("error stating modified file: %s", err.Error())
	}
	if afi.ModTime().Equal(cfi.ModTime()) {
		t.Fatal("modified and unmodified file mtimes equal")
	}

	newest, err := NewestModTime(dir)
	if err != nil {
		t.Fatalf("error finding newest mod time: %s", err.Error())
	}
	if !newest.Equal(cfi.ModTime()) {
		t.Fatal("expected newest mod time to match c")
	}
}

func TestOldestModTime(t *testing.T) {
	t.Parallel()
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		t.Fatalf("error creating temp dir: %s", err.Error())
	}
	defer os.RemoveAll(dir)
	for _, name := range []string{"a", "b", "c", "d"} {
		out := filepath.Join(dir, name)
		if err := ioutil.WriteFile(out, []byte("hi!"), 0644); err != nil {
			t.Fatalf("error writing file: %s", err.Error())
		}
	}
	time.Sleep(10 * time.Millisecond)
	for _, name := range []string{"a", "b", "d"} {
		outName := filepath.Join(dir, name)
		outfh, err := os.OpenFile(outName, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			t.Fatalf("error opening file to append: %s", err.Error())
		}
		if _, err := outfh.WriteString("\nbye!\n"); err != nil {
			t.Fatalf("error appending to file: %s", err.Error())
		}
		if err := outfh.Close(); err != nil {
			t.Fatalf("error closing file: %s", err.Error())
		}
	}

	afi, err := os.Stat(filepath.Join(dir, "a"))
	if err != nil {
		t.Fatalf("error stating unmodified file: %s", err.Error())
	}

	outName := filepath.Join(dir, "c")
	cfi, err := os.Stat(outName)
	if err != nil {
		t.Fatalf("error stating modified file: %s", err.Error())
	}
	if afi.ModTime().Equal(cfi.ModTime()) {
		t.Fatal("modified and unmodified file mtimes equal")
	}

	newest, err := OldestModTime(dir)
	if err != nil {
		t.Fatalf("error finding oldest mod time: %s", err.Error())
	}
	if !newest.Equal(cfi.ModTime()) {
		t.Fatal("expected newest mod time to match c")
	}
}
