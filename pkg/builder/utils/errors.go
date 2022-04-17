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
)

type fatalErr struct {
	code int
	error
}

func (f fatalErr) ExitStatus() int {
	return f.code
}

type exitStatus interface {
	ExitStatus() int
}

// Fatal returns an error that will cause builder to print out the
// given args and exit with the given exit code.
func Fatal(code int, args ...interface{}) error {
	return fatalErr{
		code:  code,
		error: errors.New(fmt.Sprint(args...)),
	}
}

// Fatalf returns an error that will cause builder to print out the
// given message and exit with the given exit code.
func Fatalf(code int, format string, args ...interface{}) error {
	return fatalErr{
		code:  code,
		error: fmt.Errorf(format, args...),
	}
}

// ExitStatus queries the error for an exit status.  If the error is nil, it
// returns 0.  If the error does not implement ExitStatus() int, it returns 1.
// Otherwise it retiurns the value from ExitStatus().
func ExitStatus(err error) int {
	if err == nil {
		return 0
	}
	exit, ok := err.(exitStatus)
	if !ok {
		return 1
	}
	return exit.ExitStatus()
}
