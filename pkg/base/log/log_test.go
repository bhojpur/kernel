package log_test

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
	"testing"

	"github.com/bhojpur/kernel/pkg/base/log"
)

func ExampleSetLogLevel() {
	// Set the log level to Debug
	log.Level = log.LoglvlDebug
}

func TestLog_SetLevel(t *testing.T) {
	for _, test := range []struct {
		name        string
		l           log.LogLevel
		expectError bool
	}{
		// These two cases require users to do something pretty explcitly
		// wrong to hit, but they're worth catching
		{"log level is too low", log.LogLevel(-1), true},
		{"log level is too high", log.LogLevel(6), true},

		// Included log levels
		{"log.LoglvlDebug is valid", log.LoglvlDebug, false},
		{"log.LoglvlInfo is valid", log.LoglvlInfo, false},
		{"log.LoglvlWarn is valid", log.LoglvlWarn, false},
		{"log.LoglvlError is valid", log.LoglvlError, false},
		{"log.LoglvlNone is valid", log.LoglvlNone, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := log.SetLevel(test.l)

			if err == nil && test.expectError {
				t.Error("expected error, received none")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error: %#v", err)
			}

			if !test.expectError {
				if log.Level != test.l {
					t.Errorf("log.Level should be %#v, received %#v", log.Level, test.l)
				}
			}
		})
	}
}
