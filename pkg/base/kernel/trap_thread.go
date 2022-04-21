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
	"fmt"
	"runtime"
	"syscall"
	"unsafe"

	"github.com/bhojpur/kernel/pkg/base/drivers/pic"
	"github.com/bhojpur/kernel/pkg/base/kernel/trap"
	"github.com/bhojpur/kernel/pkg/base/log"
)

var (
	// irqset，IRQ_BASE+1<<bit
	irqset uintptr

	traptask threadptr
)

func runTrapThread() {
	runtime.LockOSThread()
	var trapset uintptr
	var err syscall.Errno
	const setsize = unsafe.Sizeof(irqset) * 8

	my := Mythread()
	traptask = (threadptr)(unsafe.Pointer(my))
	log.Infof("[trap] tid:%d", my.id)

	for {
		trapset, _, err = syscall.Syscall(SYS_WAIT_IRQ, 0, 0, 0)
		if err != 0 {
			throw("bad SYS_WAIT_IRQ return")
		}
		for i := uintptr(0); i < setsize; i++ {
			if trapset&(1<<i) == 0 {
				continue
			}
			trapno := uintptr(pic.IRQ_BASE + i)

			handler := trap.Handler(int(trapno))
			if handler == nil {
				fmt.Printf("trap handler for %d not found\n", trapno)
				pic.EOI(trapno)
				continue
			}
			handler()
		}
	}
}

//go:nosplit
func wakeIRQ(no uintptr) {
	irqset |= 1 << (no - pic.IRQ_BASE)
	wakeup(&irqset, 1)
	Yield()
}

//go:nosplit
func waitIRQ() uintptr {
	if irqset != 0 {
		ret := irqset
		irqset = 0
		return ret
	}
	sleepon(&irqset)
	ret := irqset
	irqset = 0
	return ret
}
