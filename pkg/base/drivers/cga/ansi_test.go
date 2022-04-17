package cga

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
	"reflect"
	"testing"
)

func TestAnsi(t *testing.T) {
	var cases = []struct {
		str    string
		action byte
		params []string
		err    error
	}{
		{
			str:    "\x1b[12;24G",
			action: 'G',
			params: []string{"12", "24"},
		},
		{
			str:    "\x1b[12;;24G",
			action: 'G',
			params: []string{"12", "", "24"},
		},
		{
			str:    "\x1b[G",
			action: 'G',
			params: []string{},
		},
		{
			str: "X\x1b[G",
			err: errNormalChar,
		},
		{
			str: "\x1bX[G",
			err: errInvalidChar,
		},
		{
			str: "\x1b[\x00G",
			err: errInvalidChar,
		},
		{
			str: "\x1b[12\x00",
			err: errInvalidChar,
		},
	}

	p := ansiParser{}
	for _, test := range cases {
		p.Reset()
	runcase:
		for i := range test.str {
			err := p.step(test.str[i])
			if err == nil {
				continue
			}
			if err != errCSIDone {
				if err == test.err {
					break runcase
				} else {
					t.Fatalf("%q[%d] expect %v got %v", test.str, i, test.err, err)
				}
			}
			if test.action != p.Action() {
				t.Fatalf("%q[%d] expect %v got %v", test.str, i, test.action, p.Action())
			}
			if !reflect.DeepEqual(test.params, p.Params()) {
				t.Fatalf("%q[%d] expect %q got %q", test.str, i, test.params, p.Params())
			}
		}
	}
}
