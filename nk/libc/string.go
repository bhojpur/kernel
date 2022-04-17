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
	"unsafe"

	"github.com/bhojpur/kernel/nk/libc/sys/types"
)

// RawMem represents the biggest byte array the runtime can handle
type RawMem [1<<30 - 1]byte

// RawMem64 represents the biggest uint64 array the runtime can handle.
type RawMem64 [unsafe.Sizeof(RawMem{}) / unsafe.Sizeof(uint64(0))]uint64

// size_t strlen(const char *s)
func Xstrlen(t *TLS, s uintptr) (r types.Size_t) {
	if s == 0 {
		return 0
	}

	for ; *(*int8)(unsafe.Pointer(s)) != 0; s++ {
		r++
	}
	return r
}

// void *memset(void *s, int c, size_t n)
func Xmemset(t *TLS, s uintptr, c int32, n types.Size_t) uintptr {
	if n != 0 {
		c := byte(c & 0xff)

		//this will make sure that on platforms where they are not equally alligned
		//we clear out the first few bytes until allignment
		bytesBeforeAllignment := s % unsafe.Alignof(uint64(0))
		if bytesBeforeAllignment > uintptr(n) {
			bytesBeforeAllignment = uintptr(n)
		}
		b := (*RawMem)(unsafe.Pointer(s))[:bytesBeforeAllignment:bytesBeforeAllignment]
		n -= types.Size_t(bytesBeforeAllignment)
		for i := range b {
			b[i] = c
		}
		if n >= 8 {
			i64 := uint64(c) + uint64(c)<<8 + uint64(c)<<16 + uint64(c)<<24 + uint64(c)<<32 + uint64(c)<<40 + uint64(c)<<48 + uint64(c)<<56
			b8 := (*RawMem64)(unsafe.Pointer(s + bytesBeforeAllignment))[: n/8 : n/8]
			for i := range b8 {
				b8[i] = i64
			}
		}
		if n%8 != 0 {
			b = (*RawMem)(unsafe.Pointer(s + bytesBeforeAllignment + uintptr(n-n%8)))[: n%8 : n%8]
			for i := range b {
				b[i] = c
			}
		}
	}
	return s
}

// void *memcpy(void *dest, const void *src, size_t n);
func Xmemcpy(t *TLS, dest, src uintptr, n types.Size_t) (r uintptr) {
	if n == 0 {
		return dest
	}

	s := (*RawMem)(unsafe.Pointer(src))[:n:n]
	d := (*RawMem)(unsafe.Pointer(dest))[:n:n]
	copy(d, s)
	return dest
}
