package kernel

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
	"unsafe"

	"github.com/bhojpur/kernel/pkg/base/drivers/qemu"
	"github.com/bhojpur/kernel/pkg/base/drivers/uart"
	"github.com/bhojpur/kernel/pkg/base/kernel/sys"
	"github.com/bhojpur/kernel/pkg/base/log"
)

var (
	panicPcs [32]uintptr
)

//go:nosplit
func throw(msg string) {
	sys.Cli()
	tf := Mythread().tf
	throwtf(tf, msg)
}

//go:nosplit
func throwtf(tf *trapFrame, msg string) {
	sys.Cli()
	n := callers(tf, panicPcs[:])
	uart.WriteString(msg)
	uart.WriteByte('\n')

	log.PrintStr("0x")
	log.PrintHex(tf.IP)
	log.PrintStr("\n")
	for i := 0; i < n; i++ {
		log.PrintStr("0x")
		log.PrintHex(panicPcs[i])
		log.PrintStr("\n")
	}

	qemu.Exit(0xff)
	for {
	}
}

//go:nosplit
func callers(tf *trapFrame, pcs []uintptr) int {
	fp := tf.BP
	var i int
	for i = 0; i < len(pcs); i++ {
		pc := deref(fp + 8)
		pcs[i] = pc
		fp = deref(fp)
		if fp == 0 {
			break
		}
	}
	return i
}

//go:nosplit
func deref(addr uintptr) uintptr {
	return *(*uintptr)(unsafe.Pointer(addr))
}
