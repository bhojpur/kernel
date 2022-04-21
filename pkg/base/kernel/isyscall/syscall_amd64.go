package isyscall

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

// must sync with kernel.trapFrame
type trapFrame struct {
	AX, BX, CX, DX    uintptr
	BP, SI, DI, R8    uintptr
	R9, R10, R11, R12 uintptr
	R13, R14, R15     uintptr

	Trapno, Err uintptr

	// pushed by hardware
	IP, CS, FLAGS, SP, SS uintptr
}

func NewRequest(tf uintptr) Request {
	return Request{
		tf: (*trapFrame)(unsafe.Pointer(tf)),
	}
}

//go:nosplit
func (t *trapFrame) NO() uintptr {
	return t.AX
}

//go:nosplit
func (t *trapFrame) Arg(n int) uintptr {
	switch n {
	case 0:
		return t.DI
	case 1:
		return t.SI
	case 2:
		return t.DX
	case 3:
		return t.R10
	case 4:
		return t.R8
	case 5:
		return t.R9
	default:
		return 0
	}
}

//go:nosplit
func (t *trapFrame) SetRet(v uintptr) {
	t.AX = v
}

//go:nosplit
func (t *trapFrame) Ret() uintptr {
	return t.AX
}
