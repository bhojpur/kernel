package sys

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

import "unsafe"

const PtrSize = 4 << (^uintptr(0) >> 63) // unsafe.Sizeof(uintptr(0)) but an ideal const

const PageSize = 4 << 10

//go:nosplit
func Outb(port uint16, data byte)

//go:nosplit
func Inb(port uint16) byte

//go:nosplit
func Outl(port uint16, data uint32)

//go:nosplit
func Inl(port uint16) uint32

//go:nosplit
func Cli()

//go:nosplit
func Sti()

//go:nosplit
func Hlt()

//go:nosplit
func Cr2() uintptr

//go:nosplit
func Flags() uintptr

//go:nosplit
func UnsafeBuffer(p uintptr, n int) []byte {
	return (*[1 << 30]byte)(unsafe.Pointer(p))[:n]
}

//go:nosplit
func Memclr(p uintptr, n int) {
	s := (*(*[1 << 30]byte)(unsafe.Pointer(p)))[:n]
	// the compiler will emit runtime.memclrNoHeapPointers
	for i := range s {
		s[i] = 0
	}
}

// funcPC returns the entry PC of the function f.
// It assumes that f is a func value. Otherwise the behavior is undefined.
// CAREFUL: In programs with plugins, funcPC can return different values
// for the same function (because there are actually multiple copies of
// the same function in the address space). To be safe, don't use the
// results of this function in any == expression. It is only safe to
// use the result as an address at which to start executing code.
//go:nosplit
func FuncPC(f interface{}) uintptr {
	return **(**uintptr)(unsafe.Pointer((uintptr(unsafe.Pointer(&f)) + PtrSize)))
}

//go:nosplit
func Fxsave(addr uintptr)

//go:nosplit
func SetAX(val uintptr)

//go:nosplit
func CS() uintptr
