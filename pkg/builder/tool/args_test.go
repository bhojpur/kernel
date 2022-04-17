package tool

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
	"bytes"
	"testing"
)

func TestArgs(t *testing.T) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/args",
		Stderr: stderr,
		Stdout: stdout,
		Args:   []string{"status", "say", "hi", "bob", "count", "5", "status", "wait", "5ms", "cough", "false"},
	}
	code := Invoke(inv)
	if code != 0 {
		t.Log(stderr.String())
		t.Fatalf("expected 1, but got %v", code)
	}
	actual := stdout.String()
	expected := `status
saying hi bob
01234
status
waiting 5ms
not coughing
`
	if actual != expected {
		t.Fatalf("output is not expected:\n%q", actual)
	}
}

func TestBadIntArg(t *testing.T) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/args",
		Stderr: stderr,
		Stdout: stdout,
		Args:   []string{"count", "abc123"},
	}
	code := Invoke(inv)
	if code != 2 {
		t.Log("stderr:", stderr)
		t.Log("stdout:", stdout)
		t.Fatalf("expected code 2, but got %v", code)
	}
	actual := stderr.String()
	expected := "can't convert argument \"abc123\" to int\n"

	if actual != expected {
		t.Fatalf("output is not expected:\n%q", actual)
	}
}

func TestBadBoolArg(t *testing.T) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/args",
		Stderr: stderr,
		Stdout: stdout,
		Args:   []string{"cough", "abc123"},
	}
	code := Invoke(inv)
	if code != 2 {
		t.Log("stderr:", stderr)
		t.Log("stdout:", stdout)
		t.Fatalf("expected code 2, but got %v", code)
	}
	actual := stderr.String()
	expected := "can't convert argument \"abc123\" to bool\n"

	if actual != expected {
		t.Fatalf("output is not expected:\n%q", actual)
	}
}

func TestBadDurationArg(t *testing.T) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/args",
		Stderr: stderr,
		Stdout: stdout,
		Args:   []string{"wait", "abc123"},
	}
	code := Invoke(inv)
	if code != 2 {
		t.Log("stderr:", stderr)
		t.Log("stdout:", stdout)
		t.Fatalf("expected code 2, but got %v", code)
	}
	actual := stderr.String()
	expected := "can't convert argument \"abc123\" to time.Duration\n"

	if actual != expected {
		t.Fatalf("output is not expected:\n%q", actual)
	}
}

func TestMissingArgs(t *testing.T) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/args",
		Stderr: stderr,
		Stdout: stdout,
		Args:   []string{"say", "hi"},
	}
	code := Invoke(inv)
	if code != 2 {
		t.Log("stderr:", stderr)
		t.Log("stdout:", stdout)
		t.Fatalf("expected code 2, but got %v", code)
	}
	actual := stderr.String()
	expected := "not enough arguments for target \"Say\", expected 2, got 1\n"

	if actual != expected {
		t.Fatalf("output is not expected:\n%q", actual)
	}
}

func TestDocs(t *testing.T) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/args",
		Stderr: stderr,
		Stdout: stdout,
		Help:   true,
		Args:   []string{"say"},
	}
	code := Invoke(inv)
	if code != 0 {
		t.Log("stderr:", stderr)
		t.Log("stdout:", stdout)
		t.Fatalf("expected code 0, but got %v", code)
	}
	actual := stdout.String()
	expected := `Say says something. It's pretty cool. I think you should try it.

Usage:

	builder say <msg> <name>

Aliases: speak

`
	if actual != expected {
		t.Fatalf("output is not expected:\n%q", actual)
	}
}

func TestMgF(t *testing.T) {
	stderr := &bytes.Buffer{}
	stdout := &bytes.Buffer{}
	inv := Invocation{
		Dir:    "./testdata/args",
		Stderr: stderr,
		Stdout: stdout,
		Args:   []string{"HasDep"},
	}
	code := Invoke(inv)
	if code != 0 {
		t.Log("stderr:", stderr)
		t.Log("stdout:", stdout)
		t.Fatalf("expected code 0, but got %v", code)
	}
	actual := stdout.String()
	expected := "saying hi Susan\n"
	if actual != expected {
		t.Fatalf("output is not expected: %q", actual)
	}
}
