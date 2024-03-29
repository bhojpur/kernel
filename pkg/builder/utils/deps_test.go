package utils

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
	"sync/atomic"
	"testing"
	"time"
)

func TestDepsRunOnce(t *testing.T) {
	done := make(chan struct{})
	f := func() {
		done <- struct{}{}
	}
	go Deps(f, f)
	select {
	case <-done:
		// cool
	case <-time.After(time.Millisecond * 100):
		t.Fatal("func not run in a reasonable amount of time.")
	}
	select {
	case <-done:
		t.Fatal("func run twice!")
	case <-time.After(time.Millisecond * 100):
		// cool... this should be plenty of time for the goroutine to have run
	}
}

func TestDepsOfDeps(t *testing.T) {
	ch := make(chan string, 3)
	// this->f->g->h
	h := func() {
		ch <- "h"
	}
	g := func() {
		Deps(h)
		ch <- "g"
	}
	f := func() {
		Deps(g)
		ch <- "f"
	}
	Deps(f)

	res := <-ch + <-ch + <-ch

	if res != "hgf" {
		t.Fatal("expected h then g then f to run, but got " + res)
	}
}

func TestSerialDeps(t *testing.T) {
	ch := make(chan string, 3)
	// this->f->g->h
	h := func() {
		ch <- "h"
	}
	g := func() {
		ch <- "g"
	}
	f := func() {
		SerialDeps(g, h)
		ch <- "f"
	}
	Deps(f)

	res := <-ch + <-ch + <-ch

	if res != "ghf" {
		t.Fatal("expected g then h then f to run, but got " + res)
	}
}

func TestDepError(t *testing.T) {
	// TODO: this test is ugly and relies on implementation details. It should
	// be recreated as a full-stack test.

	f := func() error {
		return errors.New("ouch")
	}
	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("expected panic, but didn't get one")
		}
		actual := fmt.Sprint(err)
		if "ouch" != actual {
			t.Fatalf(`expected to get "ouch" but got "%s"`, actual)
		}
	}()
	Deps(f)
}

func TestDepFatal(t *testing.T) {
	f := func() error {
		return Fatal(99, "ouch!")
	}
	defer func() {
		v := recover()
		if v == nil {
			t.Fatal("expected panic, but didn't get one")
		}
		actual := fmt.Sprint(v)
		if "ouch!" != actual {
			t.Fatalf(`expected to get "ouch!" but got "%s"`, actual)
		}
		err, ok := v.(error)
		if !ok {
			t.Fatalf("expected recovered val to be error but was %T", v)
		}
		code := ExitStatus(err)
		if code != 99 {
			t.Fatalf("Expected exit status 99, but got %v", code)
		}
	}()
	Deps(f)
}

func TestDepTwoFatal(t *testing.T) {
	f := func() error {
		return Fatal(99, "ouch!")
	}
	g := func() error {
		return Fatal(11, "bang!")
	}
	defer func() {
		v := recover()
		if v == nil {
			t.Fatal("expected panic, but didn't get one")
		}
		actual := fmt.Sprint(v)
		// order is non-deterministic, so check for both orders
		if "ouch!\nbang!" != actual && "bang!\nouch!" != actual {
			t.Fatalf(`expected to get "ouch!" and "bang!" but got "%s"`, actual)
		}
		err, ok := v.(error)
		if !ok {
			t.Fatalf("expected recovered val to be error but was %T", v)
		}
		code := ExitStatus(err)
		// two different error codes returns, so we give up and just use error
		// code 1.
		if code != 1 {
			t.Fatalf("Expected exit status 1, but got %v", code)
		}
	}()
	Deps(f, g)
}

func TestDepWithUnhandledFunc(t *testing.T) {
	defer func() {
		err := recover()
		_, ok := err.(error)
		if !ok {
			t.Fatalf("Expected type error from panic")
		}
	}()
	var NotValid func(string) string = func(a string) string {
		return a
	}
	Deps(NotValid)
}

func TestDepsErrors(t *testing.T) {
	var hRan, gRan, fRan int64

	h := func() error {
		atomic.AddInt64(&hRan, 1)
		return errors.New("oops")
	}
	g := func() {
		Deps(h)
		atomic.AddInt64(&gRan, 1)
	}
	f := func() {
		Deps(g, h)
		atomic.AddInt64(&fRan, 1)
	}

	defer func() {
		err := recover()
		if err == nil {
			t.Fatal("expected f to panic")
		}
		if hRan != 1 {
			t.Fatalf("expected h to run once, but got %v", hRan)
		}
		if gRan > 0 {
			t.Fatalf("expected g to panic before incrementing gRan to run, but got %v", gRan)
		}
		if fRan > 0 {
			t.Fatalf("expected f to panic before incrementing fRan to run, but got %v", fRan)
		}
	}()
	f()
}
