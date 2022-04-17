package libc

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

// Copyright 2020 The Libc Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

import (
	"fmt"
	"unsafe"

	"github.com/bhojpur/kernel/nk/libc/sys/types"
)

func Bool32(b bool) int32 {
	if b {
		return 1
	}

	return 0
}

func CString(s string) uintptr {
	n := len(s)
	p := Xmalloc(nil, types.Size_t(n)+1)
	if p == 0 {
		return 0
	}

	copy((*RawMem)(unsafe.Pointer(p))[:n:n], s)
	*(*byte)(unsafe.Pointer(p + uintptr(n))) = 0
	return p
}

func GoString(s uintptr) string {
	if s == 0 {
		return ""
	}

	var buf []byte
	for {
		b := *(*byte)(unsafe.Pointer(s))
		if b == 0 {
			return string(buf)
		}

		buf = append(buf, b)
		s++
	}
}

// GoBytes returns a byte slice from a C char* having length len bytes.
func GoBytes(s uintptr, len int) []byte {
	if len == 0 {
		return nil
	}

	return (*RawMem)(unsafe.Pointer(s))[:len:len]
}

func X__assert_fail(tls *TLS, assertion, file uintptr, line uint32, function uintptr) {
	panic(fmt.Sprintf("Assert failed:%s %s:%d.%s", GoString(assertion), GoString(file), line, GoString(function)))
}
